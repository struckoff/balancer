package router

import (
	"github.com/pkg/errors"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal Balancer
}

func NewRouter(bal Balancer) (*Router, error) {
	h := &Router{
		bal: bal,
	}
	return h, nil
}

// AddNode adds node to kvrouter
func (h *Router) AddNode(n Node) error {
	return h.bal.AddNode(n)
}

// RemoveNode removes node from kvrouter
func (h *Router) RemoveNode(id string) error {
	return h.bal.RemoveNode(id)
}

// Returns node from kvrouter by given key.
func (h *Router) GetNode(key string) (Node, error) {
	//di := DataItem(key)
	nb, err := h.bal.LocateKey(key)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (h *Router) SetNodes(ns []Node) error {
	return h.bal.SetNodes(ns)
}

//func (h *Router) relocate(locations map[Node][]string) {
//	for n, keys := range locations {
//		go func(n Node, keys []string) {
//			if err := h.n.Move(keys, n); err != nil {
//				log.Println(err.Error())
//				return
//			}
//		}(n, keys)
//	}
//}
