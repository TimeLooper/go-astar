package finder

import (
	"container/heap"
)

const (
	directValue  = 10
	obliqueValue = 14
)

type Point struct {
	X int32
	Y int32
}

type AStarFinder struct {
	walkableChecker func(x, y int32) bool
}

func NewAStarFinder() *AStarFinder {
	return &AStarFinder{}
}

func (this *AStarFinder) SetWalkableChecker(f func(x, y int32) bool) {
	this.walkableChecker = f
}

func (this *AStarFinder) Find(fromX, fromY, toX, toY int32) []*Point {
	nm := make(nodeMap)
	openList := &nodeHeap{}
	heap.Init(openList)
	fromNode := nm.get(fromX, fromY, toX, toY, 0, nil)
	// 设置为open
	fromNode.SetBitState(true, 0)
	heap.Push(openList, fromNode)
	var result []*Point
	for {
		if openList.Len() <= 0 {
			return nil
		}
		current := heap.Pop(openList).(*node)
		// 设置为非open状态
		current.SetBitState(false, 0)
		// 设置为close状态
		current.SetBitState(true, 1)
		if current.x == toX && current.y == toY {
			curr := current
			for curr != nil {
				result = append(result, &Point{X: int(curr.x), Y: int(curr.y)})
				curr = curr.parent
			}
			length := len(result)
			// 反转即为路径结果
			for i := 0; i < length>>1; i++ {
				result[i], result[length-i-1] = result[length-i-1], result[i]
			}
			return result
		}
		this.checkAdjacentNode(nm, openList, current, current.x, current.y-1, toX, toY, directValue)
		this.checkAdjacentNode(nm, openList, current, current.x+1, current.y, toX, toY, directValue)
		this.checkAdjacentNode(nm, openList, current, current.x, current.y+1, toX, toY, directValue)
		this.checkAdjacentNode(nm, openList, current, current.x-1, current.y, toX, toY, directValue)

		if this.walkableChecker(current.x, current.y-1) && this.walkableChecker(current.x+1, current.y) {
			this.checkAdjacentNode(nm, openList, current, current.x+1, current.y-1, toX, toY, obliqueValue)
		}
		if this.walkableChecker(current.x+1, current.y) && this.walkableChecker(current.x, current.y+1) {
			this.checkAdjacentNode(nm, openList, current, current.x+1, current.y+1, toX, toY, obliqueValue)
		}
		if this.walkableChecker(current.x, current.y+1) && this.walkableChecker(current.x-1, current.y) {
			this.checkAdjacentNode(nm, openList, current, current.x-1, current.y+1, toX, toY, obliqueValue)
		}
		if this.walkableChecker(current.x-1, current.y) && this.walkableChecker(current.x, current.y-1) {
			this.checkAdjacentNode(nm, openList, current, current.x-1, current.y-1, toX, toY, obliqueValue)
		}
	}
}

func (this *AStarFinder) checkAdjacentNode(nm nodeMap, openList *nodeHeap, searchNode *node, x, y, toX, toY, cost int32) {
	if this.walkableChecker(x, y) {
		node := nm.get(x, y, toX, toY, cost, searchNode)
		if !node.GetBitState(0) && !node.GetBitState(1) {
			node.SetBitState(true, 0)
			heap.Push(openList, node)
		} else if searchNode.gValue+cost < node.gValue {
			node.gValue = searchNode.gValue + cost
			node.parent = searchNode
			heap.Fix(openList, int(node.index))
		}
	}
}

type node struct {
	x      int32
	y      int32
	hValue int32
	gValue int32
	index  int32
	flag   byte
	parent *node
}

func (this *node) SetBitState(state bool, index int) {
	if state {
		this.flag |= byte(0x1) << index
	} else {
		this.flag &= (^(byte(0x1) << index))
	}
}

func (this *node) GetBitState(index int) bool {
	return (this.flag&(0x1<<index))>>index != 0
}

type nodeHeap []*node

func (this nodeHeap) Len() int {
	return len(this)
}

func (this nodeHeap) Less(i, j int) bool {
	a, b := this[i], this[j]
	return a.hValue+a.gValue < b.hValue+b.gValue
}

func (this nodeHeap) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
	this[i].index = int32(i)
	this[j].index = int32(j)
}

func (this *nodeHeap) Pop() interface{} {
	old := *this
	n := len(old)
	no := old[n-1]
	no.index = -1
	*this = old[0 : n-1]
	return no
}

func (this *nodeHeap) Push(x interface{}) {
	old := *this
	item := x.(*node)
	item.index = int32(len(old))
	*this = append(old, item)
}

type nodeMap map[int64]*node

func (this nodeMap) get(x, y, toX, toY, cost int32, parent *node) *node {
	n, ok := this[int64(x)<<32|int64(y)]
	if !ok {
		n = &node{
			x: int32(x),
			y: int32(y),
		}
		n.gValue = cost
		if parent != nil {
			n.gValue = parent.gValue + cost
		}
		n.parent = parent
		n.hValue = manhattanDistance(x, y, toX, toY)
		this[int64(x)<<32|int64(y)] = n
	}
	return n
}

func manhattanDistance(fromX, fromY, toX, toY int32) int32 {
	x := fromX - toX
	y := fromY - toY
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return (x + y) * directValue
}
