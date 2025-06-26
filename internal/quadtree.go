package internal

import (
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type Quad struct {
	Pos    r2.Vec
	Size   float64
	quads  [4]*Quad
	nodes  []*Node
	weight float64
}

func generateQuadtree(nodes []*Node) *Quad {
	pos, size := findBounds(nodes)
	root := &Quad{Pos: pos, Size: size, nodes: nodes}
	run(root)
	return root
}

func run(q *Quad) {
	if len(q.nodes) <= 1 {
		return
	}
	q.addQuads()
	for _, qs := range q.quads {
		qs.insertNodes(q.nodes)
		run(qs)
	}
}

func (q *Quad) addQuads() {
	size := q.Size * 0.5
	northwest := Quad{Pos: q.Pos, Size: size}
	northeast := Quad{Pos: r2.Vec{X: q.Pos.X + size, Y: q.Pos.Y}, Size: size}
	southwest := Quad{Pos: r2.Vec{X: q.Pos.X, Y: q.Pos.Y + size}, Size: size}
	southeast := Quad{Pos: r2.Vec{X: q.Pos.X + size, Y: q.Pos.Y + size}, Size: size}
	q.quads[0] = &northwest
	q.quads[1] = &northeast
	q.quads[2] = &southwest
	q.quads[3] = &southeast
}

func (q *Quad) insertNodes(nodes []*Node) {
	for _, v := range nodes {
		if q.testBounds(v) {
			q.nodes = append(q.nodes, v)
		}
	}
}

func (q *Quad) testBounds(v *Node) bool {
	innerX := v.Pos.X >= q.Pos.X && v.Pos.X < q.Pos.X+q.Size
	innerY := v.Pos.Y >= q.Pos.Y && v.Pos.Y < q.Pos.Y+q.Size
	return innerX && innerY
}

func fetchQuads(nodes []*Node) []*Quad {
	root := generateQuadtree(nodes)
	return visit(root, 0)
}

func visit(q *Quad, depth int) []*Quad {
	//fmt.Println(depth, len(q.quads))
	if len(q.nodes) <= 1 {
		return []*Quad{q}
	}
	list := make([]*Quad, 0)
	for _, qs := range q.quads {
		list = append(list, visit(qs, depth+1)...)
	}
	return list
}

func findBounds(nodes []*Node) (r2.Vec, float64) {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64

	for _, v := range nodes {
		minX = math.Min(minX, v.Pos.X)
		minY = math.Min(minY, v.Pos.Y)
		maxX = math.Max(maxX, v.Pos.X)
		maxY = math.Max(maxY, v.Pos.Y)
	}

	var size float64
	if maxX > maxY {
		size = maxX - minX
	} else {
		size = maxY - minY
	}

	return r2.Vec{X: minX, Y: minY}, size
}
