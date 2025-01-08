package main

import (
	"fmt"
	"forces/forces"
	"math/rand"
)


func main(){

    numLinks := 4
    nodes := make([]*forces.Node, 10)
    links := make([]*forces.Link, 0)
    for i := range len(nodes) {
        nodes[i] = &forces.Node{Id: i} 
    }
    for range numLinks {
        dst := rand.Intn(len(nodes))
        src := rand.Intn(len(nodes))
        if dst == src {
            fmt.Println("crash")
            continue
        }
        links =append(links, &forces.Link{Src:src, Dst:dst})
        fmt.Println(dst,src)
    }

    sim := forces.Sim_init(nodes, links)
    sim.Init()

    for !sim.Step() {
        //fmt.Println(sim)
    }

}
