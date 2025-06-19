package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"gonum.org/v1/gonum/spatial/r2"
)

type Link struct {
	Src int64
	Dst int64
}
type Node struct {
	IDVal int64
	Pos   r2.Vec
	Disp  r2.Vec
}

type Node2 struct {
	Id int64
	X  float64
	Y  float64
}

func (n Node) ID() int64 {
	return n.IDVal
}

func main() {

	http.HandleFunc("/graph", graph)

	http.ListenAndServe(":8080", nil)
}

const WIDTH = 1920.0
const HEIGHT = 1080.0
const SCALE = 500

func graph(w http.ResponseWriter, r *http.Request) {

	fmt.Println("request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	numLinks := 500
	nodes := make([]*Node, 1000)
	links := make([]*Link, 0)
	for i := range nodes {
		nodes[i] = &Node{IDVal: int64(i), Pos: r2.Vec{X: WIDTH/2 + rand.Float64()*WIDTH/4, Y: HEIGHT/2 + rand.Float64()*HEIGHT/4}}

	}
	for range numLinks {
		dst := rand.Intn(len(nodes))
		src := rand.Intn(len(nodes))
		if dst == src {
			continue
		}

		// Add edges
		links = append(links, &Link{Src: int64(src), Dst: int64(dst)})
	}

	area := WIDTH * HEIGHT
	k := math.Sqrt(area / float64(len(nodes)))

	start := time.Now()

	fmt.Println("Start sim")
	iterations := 100
	skipcount := 0
	for iter := 0; iter < iterations; iter++ {
		// Reset displacements
		for _, v := range nodes {
			v.Disp = r2.Vec{}
		}

		// Repulsive forces
		for i := range nodes {
			for j := i + 1; j < len(nodes); j++ {
				v := nodes[i]
				u := nodes[j]
				if v.ID() == u.ID() {
					continue
				}
				delta := r2.Sub(v.Pos, u.Pos)
				dist := math.Max(r2.Norm(delta), 0.01)
				if dist >= 250 {
					skipcount++
					continue
				}
				repulsiveForce := k * k / dist * 1000
				v.Disp = r2.Add(v.Disp, r2.Scale(repulsiveForce/dist, delta))
			}
		}

		// Attractive forces
		for _, e := range links {
			src := nodes[e.Src]
			dst := nodes[e.Dst]
			delta := r2.Sub(src.Pos, dst.Pos)
			dist := math.Max(r2.Norm(delta), 0.01)
			attractiveForce := (dist * dist) / k * 100000
			forceVec := r2.Scale(attractiveForce/dist, delta)
			src.Disp = r2.Sub(src.Disp, forceVec)
			dst.Disp = r2.Add(dst.Disp, forceVec)
		}

		// Update positions
		temp := float64(iterations-iter) / float64(iterations) * 10
		for _, v := range nodes {
			dispNorm := math.Max(r2.Norm(v.Disp), 0.01)
			v.Pos = r2.Add(v.Pos, r2.Scale(math.Min(dispNorm, temp)/dispNorm, v.Disp))
			v.Pos.X = math.Min(WIDTH, math.Max(0, v.Pos.X))
			v.Pos.Y = math.Min(HEIGHT, math.Max(0, v.Pos.Y))
		}
	}

	stop := time.Now()
	delta := stop.Sub(start)
	fmt.Println("duration", delta)

	fmt.Println("transform")
	nodes2 := make([]*Node2, len(nodes))
	for i, v := range nodes {
		vn := Node2{
			Id: v.IDVal,
			X:  v.Pos.X,
			Y:  v.Pos.Y,
		}
		nodes2[i] = &vn
	}
	exp := Export{Nodes: nodes2, Links: links}

	fmt.Println("done!")
	Encode(w, exp)
	fmt.Println(skipcount)
}

type Export struct {
	Nodes []*Node2 `json:"nodes"`
	Links []*Link  `json:"links"`
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

type Vertex struct {
	id  int
	pos r2.Vec
}

type Edge struct {
	src, dst int
}

type Graph struct {
	converged bool
	step      float64
	energy    float64
	energy0   float64
	vertices  []*Vertex
	edges     []*Edge
}

func NewGraph(vertices []*Vertex, edges []*Edge) Graph {
	return Graph{
		converged: false,
		step:      1.0,
		energy:    math.Inf(1),
		vertices:  vertices,
		edges:     edges,
	}
}

func (g *Graph) ForceDirectedGraph() {
	for !g.converged {
		g.energy0 = g.energy
		g.energy = 0
		for i := range g.vertices {
			fx := 0.0
			fy := 0.0
			for _, e := range g.edges {
				if e.dst == i || e.src == i {
					vx := g.vertices[e.src]
					ux := g.vertices[e.dst]
					delta := r2.Sub(ux.pos, vx.pos)
					fx += (f_ax(e.dst, e.src) / r2.Norm(delta)) * delta.X
					fy += (f_ay(e.dst, e.src) / r2.Norm(delta)) * delta.Y
				}
			}
			for j := i + 1; j < len(g.vertices); j++ {
				vx := g.vertices[i]
				ux := g.vertices[j]
				delta := r2.Sub(ux.pos, vx.pos)
				//f += f_r(i, j)
			}

			g.energy += math.Pow(math.Abs(fx), 2)
		}

	}
}

func f_ax(i int, j int) float64 {
	return 0.0
}
func f_ay(i int, j int) float64 {
	return 0.0
}

func f_rx(i int, j int) float64 {
	return 0.0
}
func f_ry(i int, j int) float64 {
	return 0.0
}
