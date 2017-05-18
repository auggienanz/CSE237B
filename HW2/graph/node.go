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
	
	// Create a buffer to store data to be "processed"
	dataBuffer := make(map[string]int)
	for ;; {
		// For source node, block on source channel and wait for triggering;
		if n.IsSource {
			<- sourceChannel
		} else {
			// Listen on each input until enough data is received to process
			for inputNodeName, quantityNeeded := range n.Inputs {
				for dataBuffer[inputNodeName] < quantityNeeded {
					msg := <-channels[inputNodeName + "-" + n.Name]
					log.Printf("Node (%s): Receive <%d> from (%s)\n", n.Name, msg.Quantity, inputNodeName)
					dataBuffer[inputNodeName] += msg.Quantity
				}
			}
		}

		log.Printf("Node (%s): ----- Start processing data -----\n", n.Name)
		// Got all the data we need, so remove that data from the buffer
		for inputNodeName, quantityProcessed := range n.Inputs {
			dataBuffer[inputNodeName] -= quantityProcessed
		}
		
		// For drain node, once it finishes, trigger the drain channel
		if n.IsDrain {
			drainChannel <- Message{}
			break
		} else {
			// Send the appropriate data to the output nodes
			for outputNodeName, quantity := range n.Outputs {
				log.Printf("Node (%s): Send <%d> to (%s)\n", n.Name, quantity, outputNodeName)
				channels[n.Name + "-" + outputNodeName] <- Message{quantity}
			}
		}
	}

}
