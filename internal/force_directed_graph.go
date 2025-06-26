package internal

import (
	"gonum.org/v1/gonum/spatial/r2"
	"math"
	"math/rand"
)

const ITERATIONS = 125
const SCALE = 750.0
const ATTRACT = 12500
const REPULSE = 750
const SUPERNODE_THRESHOLD = 1.2

type Link struct {
	Src int64
	Dst int64
}
type Node struct {
	Id   int64
	Pos  r2.Vec
	Disp r2.Vec
}

func Simulate(nodes []*Node, links []*Link, itt int) ([]*Node, []*Link, []*Quad) {
	rng := rand.New(rand.NewSource(1))

	for _, v := range nodes {
		v.Pos.X = SCALE + rng.Float64()
		v.Pos.Y = SCALE/1.5 + rng.Float64()
	}

	area := SCALE * SCALE
	k := math.Sqrt(area / float64(len(nodes)))
	k_area := k * k

	for iter := range itt {

		// Reset displacements
		for _, v := range nodes {
			v.Disp = r2.Vec{}
		}

		root := generateQuadtree(nodes)

		// Repulsive forces
		for _, v := range nodes {
			v.Disp = r2.Add(v.Disp, superNodes(root, v.Pos, k_area))
		}

		// Attractive forces
		for _, e := range links {
			src := nodes[e.Src]
			dst := nodes[e.Dst]
			delta := r2.Sub(src.Pos, dst.Pos)
			dist := math.Max(r2.Norm(delta), 0.01)
			attractiveForce := (dist * dist) / k * ATTRACT
			forceVec := r2.Scale(attractiveForce/dist, delta)
			src.Disp = r2.Sub(src.Disp, forceVec)
			dst.Disp = r2.Add(dst.Disp, forceVec)
		}

		// Update positions
		temp := float64(itt-iter) / float64(itt) * 10
		for _, v := range nodes {
			dispNorm := math.Max(r2.Norm(v.Disp), 0.01)
			v.Pos = r2.Add(v.Pos, r2.Scale(math.Min(dispNorm, temp)/dispNorm, v.Disp))
		}
	}

	quads := fetchQuads(nodes)

	return nodes, links, quads

}

func superNodes(q *Quad, node r2.Vec, k_area float64) r2.Vec {
	var disp r2.Vec

	if len(q.nodes) == 0 {
		return disp
	}

	center := r2.Vec{X: q.Pos.X + q.Size/2, Y: q.Pos.Y + q.Size/2}
	delta := r2.Sub(node, center)
	distance := math.Max(r2.Norm(delta), 0.1)
	if len(q.nodes) <= 1 || q.Size/distance <= SUPERNODE_THRESHOLD {
		repulsion := k_area / distance * REPULSE * float64(len(q.nodes))
		disp = r2.Add(disp, r2.Scale(repulsion/distance, delta))
	} else {
		for _, qs := range q.quads {
			disp = r2.Add(disp, superNodes(qs, node, k_area))
		}
	}
	return disp
}
