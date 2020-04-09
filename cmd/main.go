package main

import (
	"encoding/json"
	"flag"
	"github.com/struckoff/kvrouter"
	"github.com/struckoff/kvrouter/config"
	"github.com/struckoff/kvrouter/rpcapi"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf config.Config
	errCh := make(chan error)
	// If config implies use of consul, consul agent name  will be  used as name.
	// Otherwise, hostname will be used instead.
	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	configFile, err := os.Open(*cfgPath)
	if err != nil {
		return err
	}
	defer configFile.Close()
	if err := json.NewDecoder(configFile).Decode(&conf); err != nil {
		return err
	}

	h, err := kvrouter.NewHost(&conf)
	if err != nil {
		return err
	}

	//Run API servers
	go func(errCh chan error) {
		if err := h.RunHTTPServer(conf.Address); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	go func(errCh chan error, h *kvrouter.Host, conf *config.Config) {
		if err := RunRPCServer(h, conf); err != nil {
			errCh <- err
			return
		}
	}(errCh, h, &conf)

	return <-errCh
}

func RunRPCServer(h rpcapi.RPCBalancerServer, conf *config.Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	rpcapi.RegisterRPCBalancerServer(s, h)

	if err := s.Serve(inbound); err != nil {
		return err
	}
	return nil
}
