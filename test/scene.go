package test

type Scene struct {
	width    int32
	height   int32
	sceneMap [][]int8
}

func NewScene(sc [][]int8) *Scene {
	s := &Scene{
		sceneMap: sc,
	}
	s.width = int32(len(sc))
	s.height = int32(len(sc[0]))
	return s
}

func (s *Scene) Walkable(x, y int32) bool {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return false
	}
	return s.sceneMap[x][y] > 0
}
