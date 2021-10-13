package finder

import (
	"container/heap"
)

const (
	directValue  = 10
	obliqueValue = 14
)

// Point grid descriptor
type Point struct {
	X int32
	Y int32
}

// AStarFinder A* algorighm
type AStarFinder struct {
	walkableChecker func(x, y int32) bool
}

// NewAStarFinder create A* finder
func NewAStarFinder() *AStarFinder {
	return &AStarFinder{}
}

// SetWalkableChecker grid walkable callback
func (asf *AStarFinder) SetWalkableChecker(f func(x, y int32) bool) {
	asf.walkableChecker = f
}

// Find calculate road from start grid to end grid
func (asf *AStarFinder) Find(fromX, fromY, toX, toY int32) []*Point {
	nm := make(nodeMap)
	openList := &nodeHeap{}
	heap.Init(openList)
	fromNode := nm.get(fromX, fromY, toX, toY, 0, nil)
	// 设置为open
	fromNode.setBitState(true, 0)
	heap.Push(openList, fromNode)
	var result []*Point
	for {
		if openList.Len() <= 0 {
			return nil
		}
		current := heap.Pop(openList).(*node)
		// 设置为非open状态
		current.setBitState(false, 0)
		// 设置为close状态
		current.setBitState(true, 1)
		if current.x == toX && current.y == toY {
			curr := current
			for curr != nil {
				result = append(result, &Point{X: int32(curr.x), Y: int32(curr.y)})
				curr = curr.parent
			}
			length := len(result)
			// 反转即为路径结果
			for i := 0; i < length>>1; i++ {
				result[i], result[length-i-1] = result[length-i-1], result[i]
			}
			return result
		}
		asf.checkAdjacentNode(nm, openList, current, current.x, current.y-1, toX, toY, directValue)
		asf.checkAdjacentNode(nm, openList, current, current.x+1, current.y, toX, toY, directValue)
		asf.checkAdjacentNode(nm, openList, current, current.x, current.y+1, toX, toY, directValue)
		asf.checkAdjacentNode(nm, openList, current, current.x-1, current.y, toX, toY, directValue)

		if asf.walkableChecker(current.x, current.y-1) && asf.walkableChecker(current.x+1, current.y) {
			asf.checkAdjacentNode(nm, openList, current, current.x+1, current.y-1, toX, toY, obliqueValue)
		}
		if asf.walkableChecker(current.x+1, current.y) && asf.walkableChecker(current.x, current.y+1) {
			asf.checkAdjacentNode(nm, openList, current, current.x+1, current.y+1, toX, toY, obliqueValue)
		}
		if asf.walkableChecker(current.x, current.y+1) && asf.walkableChecker(current.x-1, current.y) {
			asf.checkAdjacentNode(nm, openList, current, current.x-1, current.y+1, toX, toY, obliqueValue)
		}
		if asf.walkableChecker(current.x-1, current.y) && asf.walkableChecker(current.x, current.y-1) {
			asf.checkAdjacentNode(nm, openList, current, current.x-1, current.y-1, toX, toY, obliqueValue)
		}
	}
}

func (asf *AStarFinder) checkAdjacentNode(nm nodeMap, openList *nodeHeap, searchNode *node, x, y, toX, toY, cost int32) {
	if asf.walkableChecker(x, y) {
		node := nm.get(x, y, toX, toY, cost, searchNode)
		if !node.getBitState(0) && !node.getBitState(1) {
			node.setBitState(true, 0)
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

func (n *node) setBitState(state bool, index int) {
	if state {
		n.flag |= byte(0x1) << index
	} else {
		n.flag &= (^(byte(0x1) << index))
	}
}

func (n *node) getBitState(index int) bool {
	return (n.flag&(0x1<<index))>>index != 0
}

type nodeHeap []*node

func (nh nodeHeap) Len() int {
	return len(nh)
}

func (nh nodeHeap) Less(i, j int) bool {
	a, b := nh[i], nh[j]
	return a.hValue+a.gValue < b.hValue+b.gValue
}

func (nh nodeHeap) Swap(i, j int) {
	nh[i], nh[j] = nh[j], nh[i]
	nh[i].index = int32(i)
	nh[j].index = int32(j)
}

func (nh *nodeHeap) Pop() interface{} {
	old := *nh
	n := len(old)
	no := old[n-1]
	no.index = -1
	*nh = old[0 : n-1]
	return no
}

func (nh *nodeHeap) Push(x interface{}) {
	old := *nh
	item := x.(*node)
	item.index = int32(len(old))
	*nh = append(old, item)
}

type nodeMap map[int64]*node

func (nm nodeMap) get(x, y, toX, toY, cost int32, parent *node) *node {
	n, ok := nm[int64(x)<<32|int64(y)]
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
		nm[int64(x)<<32|int64(y)] = n
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
