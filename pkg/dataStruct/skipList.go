package dataStruct

import (
	"sync"
	"time"
)

//当前设计只做4层跳表，每层跨度按照2的4次方来算,跨越16个数
//后续想扩容可以改
const (
	MaxLevel   = 4
	DataLimit  = 65536
	SplitLimit = 1 << 15 //达到这个量直接分裂跳表todo
)

const StepLength = 1 << 4

type RecordStatus struct {
	sync.Mutex //加并发锁，避免读取不一致
	Id         uint64
	ExpireAt   time.Duration
	Status     uint16 //最多1024个状态(0表示已过期需要删除):这里状态变更心跳要求客户端api前加随机值发送，避免并发问题导致socket大量创建服务器崩溃
}

/*
跳表
3 *                          *
2 *        *        *        *
1 *  *  *  *  *  *  *  *  *  *
0 ****************************
*/

type DriverNode struct {
	node   *Node
	direct bool //前后驱
}

type SkipSortList struct {
	Level      int     //跳表最大层级
	Length     uint64  //跳表当前数据量
	Status     int     //0默认，对数据：1、写入磁盘 2、GC回收过期 3、数据扩容[看情况在说]
	LevelNodes []*Node //层级列表
}

// 初始化跳表
func NewSkipSortList() *SkipSortList {
	list := new(SkipSortList)
	list.Level = -1                           // 设置层级别
	list.LevelNodes = make([]*Node, MaxLevel) // 初始化头节点数组
	list.LevelNodes[0] = nil
	return list
}

// 删除节点
// 如果节点是上层节点,那么down值也必须做相应前移或者后移改动[避免掉索引层失效]
func (list *SkipSortList) Delete(key uint64) bool {
	node, ok := list.HasNode(key)
	if !ok {
		return false
	}
	//如果节点存在才删除结束
	for node.Key != 0 {
		prevNode := node.Prev
		nextNode := node.Next
		prevNode.Next = nextNode
		if nextNode.Value != nil {
			//前后直接相连
			nextNode.Prev = prevNode
		}
		node = node.Down
	}
	return true
}

func (list SkipSortList) GetVal(key uint64) (interface{}, bool) {
	node, ok := list.HasNode(key)
	if !ok {
		return nil, false
	}
	return node.Value.Load(), true
}

func (list SkipSortList) GetList(start, offset int64) []interface{} {
	//todo
	listData := make([]interface{}, 0)
	return listData
}

func (list SkipSortList) Count(start, offset int64) uint64 {
	return list.Length
}

func (list SkipSortList) HasNode(key uint64) (*Node, bool) {
	// 只有层级在大于等于 0 的时候在进行循环判断，如果层级小于 0 说明是没有任何数据
	if list.Level < 0 && !checkKey(key) {
		return nil, false
	}
	level := list.Level
	node := list.LevelNodes[level]
	for i := level; i >= 0; i-- {
		println("loop:level", i, node.Value.Load())
		current, pos := node.Ergodic(key)
		if pos == 0 {
			return current, true
		} else if pos < 0 {
			//在有效节点前，肯定无效
			return nil, false
		}
		//下层为空结束
		if current.Down == nil {
			return nil, false
		}
		node = current.Down
	}
	return nil, false
}

// 添加数据到跳表中
// 数据做好修改忍到跳表,并注意使用管道写操作日志到中继日志:这里只有 add,update,expire 其他del只有过期策略中有
// 跳表跳跃度维护
func (list *SkipSortList) randIncrLevel(level int) bool {
	//0层,当数据是16的倍数时添加1条
	//1层,当数据是256的倍数时添加1条
	//2层,当数据是4096的倍数时添加1条
	stepLength := 1 << ((level + 1) * 4)
	res := int(list.Length)%stepLength == 0
	return res
}

//每层排序插入
//3层跨2的4次方16个数
//2高层2的8次方256个数
//1层2的12次方4096个数
//0层2的16次方16次方65535个数
//当链表数据条数大于最大步长时，出现链表复制，当前数据上层数据复制为上层
//返回存在1，前后值是否相等更新
func (list *SkipSortList) CreateOrUpdate(key uint64, value interface{}) (bool, bool) {
	// 插入最底层
	if !checkKey(key) {
		return false, false
	}
	newNode := newNode().storeKV(key, value)
	//空表判断,
	if list.Level == -1 {
		list.LevelNodes[0] = newNode
		list.Level = 0
		list.Length = 1
		return false, true
	}
	if list.Length > DataLimit {
		return false, false
	}
	//不是空表,就对每一层做插入判断,开始构筑跳表
	nodes := list.LevelNodes[list.Level]
	level := list.Level
	driver := make(map[int]*DriverNode)
	for level >= 0 {
		//每层最多遍历16次,16次还找不到,节点得给上层初始化
		//找到前驱节点
		node, pos := nodes.Ergodic(newNode.Key)
		if pos == HitType {
			//找到前驱了直接更新完成
			node.Update(value)
			return true, false
		}
		//重置当前驱动
		driverNode := new(DriverNode)
		driver[level] = driverNode
		driver[level].node = node
		driver[level].direct = pos > 0
		//到底了
		if level == 0 {
			break
		}
		level--
	}
	//找不到数据就插入，找的到就更新
	ok := list.unShiftSkipTable(*newNode, driver)
	return false, ok
}

//通过驱动去处理数据
func (list *SkipSortList) unShiftSkipTable(node Node, driver map[int]*DriverNode) bool {
	//对数据条数+1
	list.Length = list.Length + 1
	currentLevel := 0
	//0层没有跳跃点
	var addNode = node
	//当层数合法才进行
	for currentLevel < MaxLevel {
		//当前驱动指针一定有，最开始是0一定添加新节点
		currentNode, ok := driver[currentLevel]
		//初始化前值后值
		listData, ok := currentNode.node.LinkIn(addNode, currentNode.direct)
		if !ok {
			return false
		}
		list.LevelNodes[currentLevel] = listData
		//判断是不是向上一层添加数据:数据条目达到上一层节点间数据量16条就向上层维护数据
		if !list.randIncrLevel(currentLevel) || currentLevel > list.Level+1 {
			//不添加上层节点直接退出
			break
		}
		currentLevel++
		node := newNode().storeKV(node.Key, node.Value.Load()).storeDown(listData)
		addNode = *node
		//要跳下一层时，发现没有下一层，就新建一层并结束
		//没有下一层,直接新建一层
		if list.Level < currentLevel {
			tailNode := driver[currentLevel-1].node.tailNode() //用最后驱动去找最大值赋值给上层
			appendNode := newNode().storeKV(tailNode.Key, tailNode.Value.Load()).storeDown(tailNode)
			//存尾部数据并跳到头部并结束
			list.LevelNodes[currentLevel] = appendNode
			list.Level = currentLevel
			break
		}
	}
	return true
}
