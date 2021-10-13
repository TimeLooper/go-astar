package test

type scene struct {
	width    int32
	height   int32
	sceneMap [][]int8
}

func newScene(sc [][]int8) *scene {
	s := &scene{
		sceneMap: sc,
	}
	s.width = int32(len(sc))
	s.height = int32(len(sc[0]))
	return s
}

func (s *scene) Walkable(x, y int32) bool {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return false
	}
	return s.sceneMap[x][y] > 0
}
