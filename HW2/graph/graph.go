package graph

import (
	"CSE237B/HW2/config"
	"log"
)

const (
	// Set a buffer size for all channels between nodes
	CHANNEL_SIZE = 10
)

var (
	// channels is a map that stores all channels between the nodes
	// The name of a channel is "input_node_name-output_node_name"
	channels = map[string](chan Message){}

	// Send a message to the source channel to trigger the processing
	sourceChannel = make(chan Message)

	// Block on the drain channel to wait for the drain node
	drainChannel = make(chan Message)
)

// Graph is a map that contains all nodes
type Graph struct {
	Nodes map[string]*Node
}

func (g *Graph) InitAllNode() {
	for _, node := range g.Nodes {
		log.Printf("Create goroutine for node (%s)\n", node.Name)
		// Create a goroutine
		go node.Run()
	}
}

func (g *Graph) Start() {
	sourceChannel <- Message{}
}

func (g *Graph) WaitEnd() {
	<-drainChannel
	log.Println("Graph processing ends")
}

type Message struct {
	Quantity int
}

// ConstructGraph constructs a new graph processing engine
// based on a graph configuration
func ConstructGraph(graphConfig *config.GraphConfig) *Graph {
	// Create a new empty graph
	graph := &Graph{
		Nodes: make(map[string]*Node),
	}

	// Parse the topology to construct the graph
	for nodeName, nodeConfig := range graphConfig.Topology {
		// Check whether this node has been created
		node, ok := graph.Nodes[nodeName]
		if !ok {
			// If not, create a new one
			node = &Node{
				Name:     nodeName,
				IsSource: false,
				IsDrain:  false,
				Inputs:   make(map[string]int),
				Outputs:  make(map[string]int),
			}
			graph.Nodes[nodeName] = node
		}

		/*
			To be completed:
			1. Initialize the newly created nodes;
			2. Create channels that are necessary for data communication.
		*/
	}

	return graph
}
