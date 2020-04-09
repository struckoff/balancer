package kvrouter

import (
	"context"
	"github.com/pkg/errors"
	"github.com/struckoff/kvrouter/balancer_adapter"
	"github.com/struckoff/kvrouter/config"
	"github.com/struckoff/kvrouter/router"
	"github.com/struckoff/kvrouter/rpcapi"
	"github.com/struckoff/kvrouter/ttl"
	"log"
	"net/http"
)

type Host struct {
	kvr    *router.Router
	checks *ttl.ChecksMap
}

func deadHandler(nodeID string) func() {
	return func() {
		log.Printf("node(%s) seems to be dead", nodeID)
	}
}

func (h *Host) removeHandler(nodeID string) func() {
	return func() {
		if err := h.kvr.RemoveNode(nodeID); err != nil {
			log.Printf("Error removing node(%s): %s", nodeID, err.Error())
		}
		h.checks.Delete(nodeID)
		log.Printf("node(%s) removed", nodeID)
	}
}

func (h *Host) RPCRegister(ctx context.Context, in *rpcapi.NodeMeta) (*rpcapi.Empty, error) {
	en, err := router.NewExternalNode(in)
	if err != nil {
		return nil, err
	}

	onDead := deadHandler(en.ID())
	onRemove := h.removeHandler(en.ID())
	check, err := ttl.NewTTLCheck(in.Check, onDead, onRemove)
	if err != nil {
		return nil, err
	}
	h.checks.Store(en.ID(), check)
	if err := h.kvr.AddNode(en); err != nil {
		return nil, err
	}
	log.Printf("node(%s) registered", en.ID())
	return &rpcapi.Empty{}, nil
}

func (h *Host) RPCHeartbeat(ctx context.Context, in *rpcapi.Ping) (*rpcapi.Empty, error) {
	if ok := h.checks.Update(in.NodeID); !ok {
		return nil, errors.Errorf("unable to found check for node(%s)", in.NodeID)
	}
	return &rpcapi.Empty{}, nil
}

func (h *Host) RunHTTPServer(addr string) error {
	r := h.kvr.HTTPHandler()
	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}

func NewHost(conf *config.Config) (*Host, error) {
	bal, err := balancer_adapter.NewSFCBalancer(conf.Balancer)
	if err != nil {
		return nil, err
	}
	kvr, err := router.NewRouter(bal)
	if err != nil {
		return nil, err
	}
	h := &Host{
		kvr:    kvr,
		checks: ttl.NewChecksMap(),
	}
	return h, nil
}
