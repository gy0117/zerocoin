package queue

import (
	"github.com/shopspring/decimal"
	"math/rand"
)

const (
	skipListMaxLevel = 32
	skipListAddLevel = 0.25 // 跳表加一层索引的概率
)

type NodeValue interface {
	GetId() string
	GetQuantity() decimal.Decimal
	GetUid() int64
	SetQuantity(decimal.Decimal)
}

type skipListLevel struct {
	forward *SkipListNode
	span    int64
}

type SkipListNode struct {
	backward *SkipListNode   // 向后指针
	level    []skipListLevel // 索引
	value    NodeValue       // 存储的值
	score    decimal.Decimal // 排序字段
}

func (slNode *SkipListNode) GetScore() decimal.Decimal {
	return slNode.score
}

func (slNode *SkipListNode) GetValue() NodeValue {
	return slNode.value
}

// Next 第i层的下一个元素
func (slNode *SkipListNode) Next(i int) *SkipListNode {
	return slNode.level[i].forward
}

// SetNext 设置第i层的下一个元素
func (slNode *SkipListNode) SetNext(i int, next *SkipListNode) {
	slNode.level[i].forward = next
}

// Span 第i层的span值
func (slNode *SkipListNode) Span(i int) int64 {
	return slNode.level[i].span
}

// SetSpan 设置第i层的span
func (slNode *SkipListNode) SetSpan(i int, span int64) {
	slNode.level[i].span = span
}

func NewSkipListNode(level int, score decimal.Decimal, value NodeValue) *SkipListNode {
	return &SkipListNode{
		backward: nil,
		level:    make([]skipListLevel, level),
		value:    value,
		score:    score,
	}
}

// https://blog.51cto.com/u_16099346/6379487，redis中的写法
type SkipList struct {
	head, tail *SkipListNode
	size       int64 // 节点个数
	level      int   // 当前最高level
	maxLevel   int
}

func NewSkipList() (*SkipList, error) {
	skipList := &SkipList{
		head:     nil,
		tail:     nil,
		size:     0,
		level:    1,
		maxLevel: skipListMaxLevel,
	}
	skipList.head = NewSkipListNode(skipList.maxLevel, decimal.NewFromInt(0), nil)
	return skipList, nil
}

func (list *SkipList) Insert(score decimal.Decimal, value NodeValue) *SkipListNode {
	// rank数组记录目标节点每一层的排名
	rank := make([]int64, list.maxLevel)
	// update数组将用于记录新节点在每一层索引的目标插入位置
	update := make([]*SkipListNode, list.maxLevel)

	p := list.head
	for i := list.level - 1; i >= 0; i-- {
		if i == list.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		// 下个节点存在，并且下个节点的score小于等于score时(score相同，按时间排序)
		for p.Next(i) != nil && p.Next(i).score.LessThanOrEqual(score) {
			rank[i] += p.level[i].span
			p = p.Next(i)
		}
		// 每一层最后一个小于给定score值的节点
		update[i] = p
	}

	// 新节点将要在哪些索引出现
	level := list.randLevel()

	// 如果level大于当前跳表的最高索引，为高出来的索引层赋初始值
	if level > list.level {
		for i := list.level; i < level; i++ {
			rank[i] = 0
			update[i] = list.head
			update[i].SetSpan(i, list.size)
		}
		list.level = level
	}
	newNode := NewSkipListNode(level, score, value)

	for i := 0; i < level; i++ {
		newNode.SetNext(i, update[i].Next(i))
		update[i].SetNext(i, newNode)

		newNode.SetSpan(i, update[i].Span(i)-(rank[0]-rank[i]))
		update[i].SetSpan(i, rank[0]-rank[i]+1)
	}

	// 处理新增节点的span
	for i := level; i < list.level; i++ {
		update[i].level[i].span++
	}
	// 处理新增节点的后退指针
	if update[0] == list.head {
		newNode.backward = nil
	} else {
		newNode.backward = update[0]
	}
	// 判断新插入的节点是不是最后一个节点
	if newNode.Next(0) != nil {
		newNode.Next(0).backward = newNode
	} else {
		list.tail = newNode
	}
	list.size++
	return newNode
}

// 随机索引层数
func (list *SkipList) randLevel() int {
	level := 1
	for rand.Int31n(100) < int32(100*skipListAddLevel) && level < list.maxLevel {
		level++
	}
	return level
}

func (list *SkipList) Find(score decimal.Decimal, id string) (*SkipListNode, []*SkipListNode) {
	update := make([]*SkipListNode, list.maxLevel)

	p := list.head
	for i := list.level - 1; i >= 0; i-- {
		for p.Next(i) != nil && p.Next(i).score.LessThan(score) {
			p = p.Next(i)
		}
		update[i] = p
	}

	for p.Next(0) != nil && p.Next(0).score.LessThanOrEqual(score) {
		p = p.Next(0)
		if p.score.Equal(score) && p.value.GetId() == id {
			break
		}
		update[0] = p
	}

	if p.score.Equal(score) && p.value.GetId() == id {
		return p, update
	}
	return nil, nil
}

func (list *SkipList) Update(score decimal.Decimal, id string, value NodeValue) *SkipListNode {
	node, _ := list.Find(score, id)
	node.value = value
	return node
}

func (list *SkipList) Delete(score decimal.Decimal, id string) {
	node, update := list.Find(score, id)
	if node == nil || node == list.head {
		return
	}

	for i := 0; i < list.level; i++ {
		if update[i].Next(i) == node {
			update[i].SetSpan(i, update[i].Span(i)+node.Span(i)-1)
			update[i].SetNext(i, node.Next(i))
		} else {
			update[i].level[i].span--
		}
	}

	if node.Next(0) == nil {
		list.tail = update[0]
	} else {
		node.Next(0).backward = update[0]
	}

	for list.level > 1 && list.head.Next(list.level-1) == nil {
		list.level--
	}
	list.size--
}

func (list *SkipList) First() *SkipListNode {
	if list.size == 0 {
		return nil
	}
	return list.head.Next(0)
}
