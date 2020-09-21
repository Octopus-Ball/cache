package cache

import (
	"sync"
)

// ************************************************************

// nodeData 节点数据
type nodeData struct {
	key  string      // 该节点在map里的key
	data interface{} // 数据
}

// listNode 链表节点
type listNode struct {
	deadline int64     // 过期时间点(纳秒时间戳)
	data     nodeData  // 节点数据
	pre      *listNode // 节点前驱
	next     *listNode // 节点后继
}

// list 双链表(带头带尾)
type list struct {
	cap  int       // 链表容量
	len  int       // 链表长度
	head *listNode // 头节点
	end  *listNode // 尾节点

	m *sync.Mutex // 互斥锁
}

// ************************************************************

// newList 新建双链表
func newList(cap int) *list {
	list := &list{
		cap:  cap,
		len:  0,
		head: newListNode("", nil),
		end:  newListNode("", nil),
	}
	list.head.next = list.end
	list.end.pre = list.head
	list.m = new(sync.Mutex)

	return list
}

// ************************************************************

// addNode 添加新节点
func (l *list) addNode(key string, data interface{}) (newNode *listNode) {
	newNode = newListNode(key, data)
	l.insertList(newNode)

	return
}

// upToHead 将某链表节点提到头部
func (l *list) upToHead(node *listNode) {
	l.takeNode(node)
	l.insertList(node)
}

// insertList 头插链表
func (l *list) insertList(node *listNode) {
	l.m.Lock()
	defer l.m.Unlock()
	node.next = l.head.next
	l.head.next.pre = node
	l.head.next = node
	node.pre = l.head
	l.len++
}

// takeNode 取下链表某个节点
func (l *list) takeNode(node *listNode) {
	l.m.Lock()
	defer l.m.Unlock()
	node.pre.next = node.next
	node.next.pre = node.pre
	l.len--
}

// newListNode 新建单链表节点
func newListNode(key string, data interface{}) *listNode {
	nodeData := nodeData{
		key:  key,
		data: data,
	}
	listNode := &listNode{
		data: nodeData,
		pre:  nil,
		next: nil,
	}
	return listNode
}

// ************************************************************
