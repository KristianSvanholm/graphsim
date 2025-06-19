package forces

import (
	"fmt"
	"math"
)

var initialRadius float64 = 15
var initialAngle float64 = math.Pi * (3 - math.Sqrt(5))

const DISTANCE float64 = 30
const STRENGTH float64 = -30

type Link struct {
	Src      int
	Dst      int
	strength float64
	bias     float64
}

type Node struct {
	Id       int
	X        float64
	Y        float64
	vx       float64
	vy       float64
	Strength float64
}

type Export struct {
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`
}

type Sim struct {
	alpha         float64
	alphaMin      float64
	alphaDecay    float64
	alphaTarget   float64
	velocityDecay float64
	iterations    int
	nodes         []*Node
	links         []*Link
}

func Sim_init(nodes []*Node, links []*Link) Sim {
	am := 0.001
	return Sim{
		alpha:         1,
		alphaMin:      am,
		alphaDecay:    0.02276277904,
		alphaTarget:   0,
		velocityDecay: 0.6,
		nodes:         nodes,
		links:         links,
	}
}

func (s *Sim) Export() Export {
	return Export{
		Nodes: s.nodes,
		Links: s.links,
	}
}

func (s *Sim) Step() bool {
	s.tick()
	return s.alpha < s.alphaMin
}

func (s *Sim) tick() {

	s.alpha += (s.alphaTarget - s.alpha) * s.alphaDecay
	s.linkforces()
	//s.manyBody()
	s.center(1920/2, 1080/2)

	// Apply velocities to node position
	for _, n := range s.nodes {
		n.vx *= s.velocityDecay // Decay velocity
		n.vy *= s.velocityDecay
		n.X += n.vx // Apply to position
		n.Y += n.vy
		//fmt.Println(n.X, n.Y)
	}
	//fmt.Println("========")
}

// Set initial node positions
func (s *Sim) Init() {
	s.nodeInit()
	s.linkInit()
}

func (s *Sim) nodeInit() {
	for i, n := range s.nodes {
		n.Id = i
		radius := initialRadius * math.Sqrt(0.5+float64(i))
		angle := float64(i) * initialAngle
		n.X = radius * math.Cos(angle)
		n.Y = radius * math.Sin(angle)
	}
}

func (s *Sim) linkInit() {

	count := make([]int, len(s.nodes))
	for _, l := range s.links {
		count[l.Src] += 1
		count[l.Dst] += 1
	}

	for _, l := range s.links {
		l.bias = float64(count[l.Src]) / float64(count[l.Src]+count[l.Dst])
		l.strength = 1 / math.Min(float64(count[l.Src]), float64(count[l.Dst]))

		fmt.Println(count[l.Src], count[l.Dst], l.strength, l.bias)
	}
}

func (s *Sim) linkforces() {
	for _, l := range s.links {
		src := s.nodes[l.Src]
		dst := s.nodes[l.Dst]

		x := dst.X + dst.vx - src.X - src.vx
		y := dst.Y + dst.vy - src.Y - src.vy

		d := math.Sqrt(x*x + y*y)
		d = (d - DISTANCE) / d * s.alpha * l.strength

		x *= d
		y *= d

		dst.vx -= x * l.bias
		dst.vy -= y * l.bias
	}
}

func (s *Sim) manyBodyInit() {
	for _, n := range s.nodes {
		n.Strength = STRENGTH
	}
}

func (s *Sim) manyBody() {

}

func (s *Sim) center(width int, height int) {
	sx := 0.0
	sy := 0.0

	for _, n := range s.nodes {
		sx += n.X
		sy += n.Y
	}

	const STRENGTH = 1

	flen := float64(len(s.nodes))
	fw := float64(width)
	fh := float64(height)

	sx = (sx/flen - fw) * STRENGTH
	sy = (sy/flen - fh) * STRENGTH

	for _, n := range s.nodes {
		n.X -= sx
		n.Y -= sy
	}
}
