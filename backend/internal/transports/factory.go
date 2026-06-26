package transports

import "fmt"

func New(cfg TransportConfig) (Transport, error) {
	constructor, ok := lookup(cfg.Protocol)
	if !ok {
		return nil, fmt.Errorf("transport %q: protocol not registered", cfg.Protocol)
	}
	return constructor(cfg)
}
