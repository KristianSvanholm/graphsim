package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/spatial/r2"
)

type Link struct {
	Src int64 `json="src"`
	Dst int64 `json="dst"`
}
type Node struct {
	Id   int64 `json="id"`
	Pos  r2.Vec
	Disp r2.Vec
}

type Node2 struct {
	Id int64
	X  float64
	Y  float64
}

func (n Node) ID() int64 {
	return n.Id
}

func main() {

	http.HandleFunc("/graph", graph)

	http.ListenAndServe(":8080", nil)
}

const HEIGHT = 1080.0
const SCALE = 500

type Quadtree struct {
	root Quad
}

type Quad struct {
	Pos    r2.Vec
	Size   float64
	quads  [4]*Quad
	nodes  []*Node
	weight float64
}

func (q *Quad) addQuads() {
	size := q.Size * 0.5
	northwest := Quad{Pos: q.Pos, Size: size}
	northeast := Quad{Pos: r2.Vec{q.Pos.X + size, q.Pos.Y}, Size: size}
	southwest := Quad{Pos: r2.Vec{q.Pos.X, q.Pos.Y + size}, Size: size}
	southeast := Quad{Pos: r2.Vec{q.Pos.X + size, q.Pos.Y + size}, Size: size}
	q.quads[0] = &northwest
	q.quads[1] = &northeast
	q.quads[2] = &southwest
	q.quads[3] = &southeast
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

	return r2.Vec{minX, minY}, size
}

func generateQuadtree(nodes []*Node) Quadtree {
	pos, size := findBounds(nodes)
	root := Quad{Pos: pos, Size: size, nodes: nodes}
	qt := Quadtree{root: root}
	run(&qt.root)
	return qt
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

func fetchQuads(nodes []*Node) []*Quad {
	qt := generateQuadtree(nodes)
	run(&qt.root)
	return visit(&qt.root, 0)
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

func superNodes(q *Quad, node r2.Vec, k_area float64) r2.Vec {
	var disp r2.Vec

	if len(q.nodes) == 0 {
		return disp
	}

	center := r2.Vec{X: q.Pos.X + q.Size/2, Y: q.Pos.Y + q.Size/2}
	delta := r2.Sub(node, center)
	distance := math.Max(r2.Norm(delta), 0.1)
	if len(q.nodes) <= 1 || q.Size/distance <= 1.2 { // TODO:: Magic number. Put in const
		repulsion := k_area / distance * 1000 * float64(len(q.nodes)) // Todo:: Magic number ´1000´. Put in const for repulsion force
		disp = r2.Add(disp, r2.Scale(repulsion/distance, delta))
	} else {
		for _, qs := range q.quads {
			disp = r2.Add(disp, superNodes(qs, node, k_area))
		}
	}
	return disp
}

func graph(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var struc Export
	data, err := os.ReadFile("dump.json")
	if err != nil {
		fmt.Println("err: ", err)
		http.Error(w, "nope", http.StatusInternalServerError)
		return
	}

	rng := rand.New(rand.NewSource(1))
	json.Unmarshal(data, &struc)
	for _, v := range struc.Nodes {
		v.Pos.X = HEIGHT/3 + rng.Float64()
		v.Pos.Y = HEIGHT/3 + rng.Float64()
	}

	nodes := struc.Nodes
	links := struc.Links
	area := HEIGHT * HEIGHT
	k := math.Sqrt(area / float64(len(nodes)))
	k_area := k * k

	start := time.Now()

	fmt.Println("Start sim ", len(nodes))
	iterations := 100

	for iter := 0; iter < iterations; iter++ {

		// Reset displacements
		for _, v := range nodes {
			v.Disp = r2.Vec{}
		}

		inner_start := time.Now()

		qt := generateQuadtree(nodes)

		quads_stop := time.Now()

		// Repulsive forces
		for _, v := range nodes {
			v.Disp = r2.Add(v.Disp, superNodes(&qt.root, v.Pos, k_area))
		}
		/*for i := range nodes {
			for j := i + 1; j < len(nodes); j++ {
				v := nodes[i]
				u := nodes[j]
				delta := r2.Sub(v.Pos, u.Pos)
				dist := math.Max(r2.Norm(delta), 0.01)
				if dist >= 250 {
					continue
				}
				repulsiveForce := k * k / dist * 1000
				v.Disp = r2.Add(v.Disp, r2.Scale(repulsiveForce/dist, delta))
			}
		}*/

		node_stop := time.Now()

		// Attractive forces
		for _, e := range links {
			src := nodes[e.Src]
			dst := nodes[e.Dst]
			delta := r2.Sub(src.Pos, dst.Pos)
			dist := math.Max(r2.Norm(delta), 0.01)
			attractiveForce := (dist * dist) / k * 12500
			forceVec := r2.Scale(attractiveForce/dist, delta)
			src.Disp = r2.Sub(src.Disp, forceVec)
			dst.Disp = r2.Add(dst.Disp, forceVec)
		}

		link_stop := time.Now()

		// Update positions
		temp := float64(iterations-iter) / float64(iterations) * 10
		for _, v := range nodes {
			dispNorm := math.Max(r2.Norm(v.Disp), 0.01)
			v.Pos = r2.Add(v.Pos, r2.Scale(math.Min(dispNorm, temp)/dispNorm, v.Disp))
		}

		update_stop := time.Now()
		fmt.Println(fmt.Sprintf("Quads: %d | Repulse: %d | Attract: %d | Update: %d", quads_stop.Sub(inner_start).Microseconds(), node_stop.Sub(quads_stop).Microseconds(), link_stop.Sub(node_stop).Microseconds(), update_stop.Sub(link_stop).Microseconds()))
	}

	quads := fetchQuads(nodes)

	stop := time.Now()
	delta := stop.Sub(start)
	fmt.Println("duration", delta)

	exp := Export{Nodes: nodes, Links: links, Quads: quads}

	Encode(w, exp)
}

type Export struct {
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`
	Quads []*Quad `json:"quads"`
}

// Encode any struct given to it
func Encode(w http.ResponseWriter, data interface{}) error {

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func graph2(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\n\n\n")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var struc graphExport
	data, err := os.ReadFile("dump.json")
	if err != nil {
		fmt.Println("err: ", err)
		http.Error(w, "nope", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(data, &struc)
	for _, v := range struc.Vertices {
		v.Pos.X = HEIGHT/2 + rand.Float64()*HEIGHT/4
		v.Pos.Y = HEIGHT/2 + rand.Float64()*HEIGHT/4
	}
	fmt.Println("pos", struc.Vertices[0].Pos, struc.Vertices[0].Id)
	g := NewGraph(struc.Vertices, struc.Edges, HEIGHT, HEIGHT)
	g.ForceDirectedGraph()

	Encode(w, g.export())

}

type Vertex struct {
	Id  int
	Pos r2.Vec
}

type Edge struct {
	Src, Dst int
}

type graphExport struct {
	Vertices []*Vertex `json:"nodes"`
	Edges    []*Edge   `json:"links"`
}

type Graph struct {
	converged bool
	progress  int     // Adaptive step length counter
	t         float64 // Step length modifier
	step      float64 // Step length
	energy    float64 // Total energy in graph
	energy0   float64 // Previous energy in graph
	vertices  []*Vertex
	edges     []*Edge
	x         []float64
	x0        []float64
	k         float64 // "Optimal distance"
	c         float64 // Relative strength of attractive and repulsive forces
	tol       float64
}

func (g *Graph) export() graphExport {
	return graphExport{
		Vertices: g.vertices,
		Edges:    g.edges,
	}
}

func NewGraph(vertices []*Vertex, edges []*Edge, w float64, h float64) Graph {
	return Graph{
		converged: false,
		step:      1,
		t:         0.9,
		k:         math.Sqrt((w*h)/float64(len(vertices))) / 10,
		x:         flattenPositions(vertices),
		c:         0.2,
		tol:       1,
		energy:    math.Inf(1),
		vertices:  vertices,
		edges:     edges,
	}
}

func flattenPositions(vs []*Vertex) []float64 {
	x := make([]float64, 0)
	for _, v := range vs {
		x = append(x, v.Pos.X, v.Pos.Y)
	}
	return x
}

func (g *Graph) checkConverged() {
	g.converged = floats.Distance(g.x, g.x0, 2) < g.k*g.tol
}

func (g *Graph) ForceDirectedGraph() {

	empty := r2.Vec{0.0, 0.0}
	fmt.Println("run")
	start := time.Now()
	for !g.converged {

		g.x0 = g.x
		g.energy0 = g.energy
		g.energy = 0

		for i := range g.vertices {
			var f r2.Vec
			for _, e := range g.edges {
				if e.Src == i {

					// NOTE:: Not knowing which is which here might be a bug.
					vx := g.vertices[e.Src]
					ux := g.vertices[e.Dst]
					delta := r2.Sub(ux.Pos, vx.Pos)
					f = r2.Add(f, r2.Scale(g.fa(e.Src, e.Dst)/r2.Norm(delta), delta))
				}
			}
			for j := i + 1; j < len(g.vertices); j++ {
				vx := g.vertices[i]
				ux := g.vertices[j]
				delta := r2.Sub(ux.Pos, vx.Pos)
				f = r2.Add(f, r2.Scale(g.fr(i, j), delta))
			}
			v := g.vertices[i]
			if f != empty {
				v.Pos = r2.Add(v.Pos, r2.Scale(g.step, r2.Unit(f)))
			}
			g.energy += r2.Norm2(f)

		}
		g.updateSteplen()
		g.x = flattenPositions(g.vertices)
		g.checkConverged()
	}

	stop := time.Now()
	delta := stop.Sub(start)
	fmt.Println(delta)
}

func (g *Graph) updateSteplen() {
	if g.energy < g.energy0 {
		//fmt.Println("helloless")
		g.progress += 1
		if g.progress >= 5 {
			g.progress = 0
			g.step /= g.t
		}
	} else {
		//fmt.Println("more")
		g.progress = 0
		g.step *= g.t
	}
}

func (g *Graph) fa(i int, j int) float64 {
	v := g.vertices[i].Pos
	u := g.vertices[j].Pos
	return r2.Norm2(r2.Sub(v, u)) / g.k
}

func (g *Graph) fr(i int, j int) float64 {
	v := g.vertices[i].Pos
	u := g.vertices[j].Pos
	return -g.c * math.Pow(g.k, 2) / r2.Norm(r2.Sub(v, u))
}
