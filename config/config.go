package config

import (
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	"strings"
)

type Config struct {
	Address    string
	RPCAddress string
	Balancer   *BalancerConfig
}

// If config implies use of consul, this options will be taken from consul KV.
// Otherwise it will be taken from config file.
type BalancerConfig struct {
	Dimensions uint64    //Amount of space filling curve dimensions
	Size       uint64    //Size of space filling curve
	Curve      CurveType //Space filling curve type
}

type CurveType struct {
	curve.CurveType
}

func (ct *CurveType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "morton":
		ct.CurveType = curve.Morton
		return nil
	case "hilbert":
		ct.CurveType = curve.Hilbert
		return nil
	default:
		return errors.New("unknown curve type")
	}
}
