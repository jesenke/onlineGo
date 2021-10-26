package dataStruct

import "sync/atomic"

const (
	RootType    int = 1
	MidType     int = 2
	HeadType    int = 0
	TailType    int = 3
	UnknownType int = -1
)

//链表最大长度
const LinkLimit = 1 << 16

const LeftKeyLimit uint64 = 0
const RightKeyLimit uint64 = 1 << 63

const HitType = 0       //命中
const HitLeftType = -1  //命中左侧
const HitRightType = 1  //命中右侧
const HitRangeType = -2 //命中右侧

type Node struct {
	Value *atomic.Value // 存储值
	Key   uint64
	Prev  *Node // 同层前节点
	Next  *Node // 同层后节点
	Down  *Node // 下层同节点
	//Above *Node // 上层同节点
}

func NewList() *Node {
	List := &Node{}
	return List
}

// 返回链表头节点
func (l *Node) Head() (head *ListNode) {
	head = l.head

	return
}

// 返回链表尾节点
func (l *Node) Tail() (tail *ListNode) {
	tail = l.tail

	return
}

func (node *Node) NodeType() int {
	if node == nil {
		return UnknownType
	}
	if node.Prev != nil && node.Next != nil {
		return MidType
	}
	//中间节点
	if node.Prev == nil && node.Next == nil {
		return RootType
	}
	//头部节点
	if node.Prev == nil {
		return HeadType
	}
	// 尾节点
	if node.Next == nil {
		return TailType
	}
	return UnknownType
}

func (node *Node) storeDown(nodeData *Node) *Node {
	node.Down = nodeData
	return node
}

func (node *Node) storeKV(key uint64, val interface{}) *Node {
	node.Value = new(atomic.Value)
	node.Value.Store(val)
	node.Key = key
	return node
}

// Next returns the next ring element. r must not be empty.
func (node *Node) NextNode() *Node {
	return node.Next
}

// Prev returns the previous ring element. r must not be empty.
func (node *Node) PrevNode() *Node {
	return node.Prev
}

//链表初始化
func newNode() *Node {
	node := new(Node)
	return node
}

//找到链表尾部
func (node *Node) tailNode() *Node {
	nodeData := node
	for nodeData.NodeType() > UnknownType {
		if nodeData.NodeType() == TailType || nodeData.NodeType() == RootType {
			break
		}
		nodeData = node.Next
	}
	return nodeData
}

func checkKey(key uint64) bool {
	if key >= LeftKeyLimit || key <= RightKeyLimit {
		return true
	}
	return false
}

func (node *Node) Compare(key uint64) int {
	if node.Key == key {
		return HitType
	} else if node.Key < key {
		return HitRightType
	} else {
		return HitLeftType
	}
}

//前置条件是非空双链表
//遍历查找,是否存在,返回一个维护好的联表结果,
//[在当前节点的哪个方向]
//(-1:前,0:直接命中找到,1:后, -1),
//找不到,确认大概范围在有效节点的后面
func (node *Node) Ergodic(key uint64) (*Node, int) {
	//数据存储日志写上add,update,expire,del
	var found = HitRangeType
	if !checkKey(key) {
		return nil, HitRangeType
	}
	child := node

	if child.NodeType() == RootType {
		return child, child.Compare(key)
	}
	if child.NodeType() == TailType {
		child = child.PrevNode()
	}
	if child.NodeType() == HeadType {
		child = child.NextNode()
	}
	//中序遍历
	for child.Key > LeftKeyLimit && child.Key < RightKeyLimit {
		//需要到下一个节点
		println(key)
		found = child.Compare(key)
		if found == HitType {
			break
		}
		if found > HitType {
			if child.NextNode() == nil || child.NextNode().Key < key {
				break
			}
			child = child.NextNode()
		} else {
			if child.PrevNode() == nil || child.PrevNode().Key > key {
				break
			}
			child = child.PrevNode()
		}
	}
	return child, found
}

//链表头部
func (node *Node) headNode() *Node {
	nodeData := node
	for nodeData.NodeType() > UnknownType {
		if nodeData.PrevNode() == nil {
			break
		}
		nodeData = nodeData.PrevNode()
	}
	return nodeData
}

//更新节点
//【写日志要保持顺序性】
//更新是判断下值是否相等
func (node *Node) Update(nodeValue interface{}) bool {
	//数据存储日志写上add,update,expire,del
	return true
}

//返回当前节点
func (node *Node) leftLink(prevNode *Node) *Node {
	if node.NodeType() == HeadType || node.NodeType() == RootType {
		node.Prev = prevNode
		prevNode.Next = node
		prevNode.Prev = nil
		return prevNode
	}
	//前节点打断
	prev := node.PrevNode()
	prevNode.Prev = prev
	prev.Next = prevNode
	node.Prev = prevNode
	prevNode.Next = node
	return prevNode
}

//返回当前节点
func (node *Node) rightLink(nextNode *Node) *Node {
	if node.NodeType() == TailType || node.NodeType() == RootType {
		node.Next = nextNode
		nextNode.Prev = node
		return nextNode
	}
	//后节点打断
	next := node.NextNode()
	next.Prev = nextNode
	nextNode.Next = next
	//联到前面去
	node.Next = nextNode
	nextNode.Prev = node
	return nextNode
}

//一个节点联入[无论哪里插入，都是前置、后置节点放入新加节点]
func (node *Node) LinkIn(nodeValue Node, direct bool) (*Node, bool) {
	//数据存储日志写上add,update,expire,del
	//当前节点插入,后一个节点作为插入节点的后一个节点
	//必须前后节点都为空,否则插入错误
	if nodeValue.NodeType() != RootType {
		return node, false
	}
	if direct {
		return node.rightLink(&nodeValue), true
	} else {
		return node.leftLink(&nodeValue), true
	}
}

//节点联入[无论哪里插入，都是前置、后置节点放入新加节点]
func (node *Node) Println(no int64) {
	//数据存储日志写上add,update,expire,del
	//当前节点插入,后一个节点作为插入节点的后一个节点
	nodeData := node
	for nodeData.NodeType() > UnknownType {
		if nodeData.NodeType() == HeadType || nodeData.NodeType() == RootType {
			break
		}
		nodeData = nodeData.PrevNode()
	}
	println(no, "start out")
	for nodeData.NodeType() > UnknownType {
		//前进指针
		if nodeData.NodeType() == TailType || nodeData.NodeType() == RootType {
			println(no, "dataKey", nodeData.Key, "value", nodeData.Value.Load())
			break
		}
		println(no, "dataKey", nodeData.Key, "value", nodeData.Value.Load())

		nodeData = nodeData.NextNode()
	}
}

//对节点进行过期回收
func (node *Node) GcNode() {
	//todo
}
