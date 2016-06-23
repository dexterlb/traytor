package mesh

import (
	"fmt"

	"github.com/DexterLB/traytor/maths"
)

// KDtree represents a node in a KD tree
type KDtree struct {
	Axis      int
	Median    float64
	Triangles []int
	Children  [2]*KDtree
}

// NewLeaf returns a new KDtree with triangles - the list of the given triangles and sets the axis to leaf
func NewLeaf(triangles []int) *KDtree {
	return &KDtree{
		Axis:      maths.Leaf,
		Triangles: triangles,
	}
}

// NewNode returns a new KDtree node with the givena axis and median
func NewNode(median float64, axis int) *KDtree {
	var children [2]*KDtree
	return &KDtree{
		Axis:     axis,
		Median:   median,
		Children: children,
	}
}

// String returns the string representation of the KDtree in the form of
// axis{median}(child1, child2)
func (t *KDtree) String() string {
	if t.Axis == maths.Leaf {
		return fmt.Sprintf("%v", t.Triangles)
	}
	return fmt.Sprintf("%d{%.3g}(%s, %s)", t.Axis, t.Median, t.Children[0], t.Children[1])
}
