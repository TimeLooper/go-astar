package finder

type Point struct {
	X int
	Y int
}

type PathFinder interface {
	SetWalkableChecker(f func(x, y int32) bool)
	Find(fromX, fromY, toX, toY int) []*Point
}
