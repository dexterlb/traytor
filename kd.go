package traytor

type KDtree struct {
	Axis      int
	Median    int
	Triangles []*Triangle
	Children  [2]*KDtree
}

func NewLeaf(triangles []*Triangle) *KDtree {
	return &KDtree{
		Axis:      Leaf,
		Triangles: triangles,
	}
}

func NewNode(median int, axis int) *KDtree {
	var children [2]*KDtree
	return &KDtree{
		Axis:     axis,
		Median:   median,
		Children: children,
	}
}
