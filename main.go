package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"forces/internal"
)

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	http.HandleFunc("/graph", graph)
	http.ListenAndServe(":8080", nil)
}

func graph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := r.URL.Query()
	itt, err := strconv.Atoi(params.Get("itt"))
	if err != nil {
		itt = 100
	}

	var data Export
	bytes, err := os.ReadFile("dump.json")
	if err != nil {
		fmt.Println("err: ", err)
		http.Error(w, "nope", http.StatusInternalServerError)
		return
	}

	json.Unmarshal(bytes, &data)

	start := time.Now()
	nodes, links, quads := internal.Simulate(data.Nodes, data.Links, itt)
	delta := time.Now().Sub(start)
	fmt.Println("itterations:", itt, "duration:", delta, "mspi:", float64(delta)/float64(itt))

	exp := Export{Nodes: nodes, Links: links, Quads: quads}

	Encode(w, exp)
}

type Export struct {
	Nodes []*internal.Node `json:"nodes"`
	Links []*internal.Link `json:"links"`
	Quads []*internal.Quad `json:"quads"`
}

// Encode any struct given to it
func Encode(w http.ResponseWriter, data any) error {
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}
