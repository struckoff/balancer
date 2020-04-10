package router

type Balancer interface {
	AddNode(n Node) error               // Add node to the balancer
	RemoveNode(id string) error         // Remove node from the balancer
	SetNodes(ns []Node) error           // Remove all nodes from the balancer and set a new ones
	LocateKey(key string) (Node, error) // Return the node for the given key
	Nodes() ([]Node, error)             // Return list of nodes in the balancer
	GetNode(id string) (Node, error)    // Return the node with the given id
}
