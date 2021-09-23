package dataStruct

import (
	"time"
)

//用户注册记录:store的时候一定带上appid
//应该是根据一定规则打散到不同的broker,提高整体读写的能力
type Record struct {
	Id        int //分布式自增数据
	AppId     int
	AccountId string     //用户直属组织id
	Name      string     //用户名
	Mobile    string     //用户电话
	Extend    [1024]byte //额外信息json的Marshal得到
	Token     string     //根据得到秘钥：方向解析要可以得到过期时间、数据id
	CreateAt  time.Duration
}

//叶子节点,考虑无限层级,首期写3层就行[账户级-部门级-频道级]若有其他情形在扩展
type RoomNode struct {
	RoomId string
	Data   []*PageChain
}

//企业数
type MemoryTree struct {
	AccountId string
	APPId     string
	Count     int //部门上房间数
	Room      map[string]RoomNode
}

//查询数据是否存在在网关api那里加布隆过滤，避免无效数据攻击
//页表,每页存65535条数据
//问题：怎么快速让大量记录排序然后分隔成多个分组：加跳表
type PageChain struct {
	PageIndex   string //页号
	PageStartId uint64 //跳表的最小id值
	PageEndId   uint64 //跳表的最大id值
	Records     *SkipSortList
	IsDirty     bool //是否脏页
}

//当前值做1级
func New(appId int, AccountId string) *MemoryTree {
	node := new(MemoryTree)
	node.Count = 0
	node.AccountId = AccountId

	return node
}

func (node *MemoryTree) HasValue(RoomId string, value *RecordStatus) bool {
	if node.Count == 0 {
		return false
	}
	dept, ok := node.Room[RoomId]
	if !ok {
		return false
	}
	for _, v := range dept.Data {
		//数据跳表内的,判断是否新增更新过期
		if v.PageStartId <= value.Id && v.PageEndId >= value.Id {
			_, ok := v.Records.HasNode(value.Id)
			//当数据达到什么情况,对跳表执行分裂todo
			return ok
		}
	}
	return false
}

func (node *MemoryTree) AddValue(RoomId string, value *RecordStatus) (bool, bool) {
	if value.Id <= 0 || RoomId == "" {
		return false, false
	}
	room, ok := node.Room[RoomId]
	if !ok {
		room = new(RoomNode)
		room.RoomId = RoomId
		room.Data = make([]*PageChain, 1)
		skipNode := NewSkipSortList()
		skipNode.CreateOrUpdate(value.Id, value)
		room.Data[0].Records = skipNode
		room.Data[0].PageStartId = value.Id
		room.Data[0].PageEndId = value.Id
	}
	for _, v := range room.Data {
		//数据跳表内的,判断是否新增更新过期
		if v.PageStartId <= value.Id && v.PageEndId >= value.Id {
			//当数据达到什么情况,对跳表执行分裂todo
			return v.Records.CreateOrUpdate(value.Id, value)
		}
		//其他情形还有很多考虑点。todo
	}
	return true, true
}
