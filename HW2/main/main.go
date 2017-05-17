package main

import (
	"CSE237B/HW2/config"
	"CSE237B/HW2/graph"
	"log"
)

const (
	CONFIG_PATH = "./config.json"
)

func main() {
	log.Printf("Start DataFlow")

	// load the configration file
	topology, err := config.ParseGraphConfig(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the graph processing engine from topology
	g := graph.ConstructGraph(topology)

	// Initialize all nodes, create goroutines
	g.InitAllNode()

	// Trigger the source node to start processing
	g.Start()

	// Wait for the drain node finishing processing
	g.WaitEnd()
}
