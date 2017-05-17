package graph

import (
	"log"
)

// Node is a concurrent processing unit
type Node struct {
	Name     string
	IsSource bool           // is it a source node?
	IsDrain  bool           // is it a drain node?
	Inputs   map[string]int // the inputs information map[node_name](data quantity)
	Outputs  map[string]int // the outputs information map[node_name](data quantity)
}

func (n *Node) Run() {
	log.Printf("Node (%s): Initiated\n", n.Name)

	/*
		To be completed:
		1. For source node, block on source channel and wait for triggering;
		2. For darin node, once it finishes, trigger the drain channel;
		3. All nodes need to block on input channels, and to start processing once it collects enough inputs,
		  	then sends messages to output channels;
		4. Print information such as:
			log.Printf("Node (%s): Receive <%d> from (%s)\n", n.Name, msg.Quantity, inputNodeName)
			log.Printf("Node (%s): ----- Start processing data ----- \n", n.Name)
			log.Printf("Node (%s): Send <%d> to (%s)\n", n.Name, quantity, outputNodeName)
	*/
}
