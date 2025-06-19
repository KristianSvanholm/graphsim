package forces

import "math"

type Quadtree struct {
	root    *Quad
	content []*Node
	nodes   []Quadtree
	total   int
}

type Quad struct {
	nodes []*Quad
	x     float64
	y     float64
	w     float64
	h     float64
}

func (q *Quadtree) addAll(nodes []*Node) {
	xz := make([]float64, len(nodes))
	yz := make([]float64, len(nodes))
	x0 := math.MaxFloat64
	y0 := x0
	x1 := -x0
	y1 := x1

	// Compute the points and their extent
	for i, n := range nodes {
		xz[i] = n.X
		yz[i] = n.Y

		// Shrink the bounds
		if n.X < x0 {
			x0 = n.X
		}
		if n.X > x1 {
			x1 = n.X
		}
		if n.Y < y0 {
			y0 = n.Y
		}
		if n.Y > y1 {
			y0 = n.Y
		}
	}

	// If there were no valid points, abort
	if x0 > x1 || y0 > y1 {
		return
	}

	// Expand the tree to cover new points
	q.cover(x0, y0).cover(x1, y1)

	// Add new points
	for i, n := range nodes {
		q.add(xz[i], yz[i], n)
	}

}

func (q *Quadtree) add(x float64, y float64, n *Node) {

}

func (q *Quadtree) cover(x float64, y float64) *Quadtree {
	_x := q.root.x
	_y := q.root.y
	_w := q.root.w
	_h := q.root.h

	b := q.root
	if b == nil {
		_x = math.Floor(x)
		_y = math.Floor(y)
	} else {
		var z float64
		if _w-_x == 0 {
			z = 1
		} else {
			z = _w - _x
		}

		node := q.root

		for b.x > x || x >= b.w || b.y > y || y >= b.h {
			i := btoi(y < b.y)<<1 | btoi(x < b.x)
			parent := make([]*Quad, 4)
			parent[i] = node
			z *= 2
			switch i {
			case 0:
				{
					_w = _x + z
					_h = _y + z
					break
				}
			case 1:
				{
					_x = _w + z
					_y = _h + z
					break
				}
			case 2:
				{
					_w = _x + z
					_y = _w + z
					break
				}
			case 3:
				{
					_x = _w + z
					_h = _y + z
					break
				}
			}

		}

		if q.root != nil && len(q.root.nodes) != 0 {
			q.root = node
		}

	}

	q.root.x = _x
	q.root.y = _y
	q.root.w = _w
	q.root.h = _h

	return q
}

func (q *Quadtree) visitAfter() {
	quads := make([]Quad, 0)
	//next := make([]Quad, 0)
	if q.root != nil {
		quads = append(quads, Quad{})
	}

}

func btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
