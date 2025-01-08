package main

import (
	"encoding/json"
	"fmt"
	"forces/forces"
	"math/rand"
	"net/http"
)


func main(){


    http.HandleFunc("/graph", graph)

    http.ListenAndServe(":8080", nil)
   }

func graph(w http.ResponseWriter, r *http.Request) {


    w.Header().Set("Access-Control-Allow-Origin", "*")
    numLinks := 50
    nodes := make([]*forces.Node, 250)
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
    
    Encode(w, sim.Export()) 
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
