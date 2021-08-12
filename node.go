package easystargo

type node struct {
	x                      int32
	y                      int32
	costSoFar              int32
	simpleDistanceToTarget int32
	parent                 *node
	index                  int
}

type nodeHeap []*node

func (nh nodeHeap) Len() int {
	return len(nh)
}

func (nh nodeHeap) Less(i, j int) bool {
	a, b := nh[i], nh[j]
	return a.costSoFar+a.simpleDistanceToTarget < b.costSoFar+b.simpleDistanceToTarget
}

func (nh nodeHeap) Swap(i, j int) {
	nh[j], nh[i] = nh[i], nh[j]
	nh[i].index = i
	nh[j].index = j
}

func (nh *nodeHeap) Push(x interface{}) {
	tmp := *nh
	n := x.(*node)
	n.index = len(tmp)
	*nh = append(tmp, n)
}

func (nh *nodeHeap) Pop() interface{} {
	tmp := *nh
	length := len(tmp)
	n := tmp[length-1]
	n.index = -1
	*nh = tmp[0 : length-1]
	return n
}
