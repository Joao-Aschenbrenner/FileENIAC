// SPDX-License-Identifier: MIT
package transports

import "sync"

type TransportConstructor func(cfg TransportConfig) (Transport, error)

var (
	registryMu sync.RWMutex
	registry   = make(map[string]TransportConstructor)
)

func Register(protocol string, constructor TransportConstructor) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[protocol] = constructor
}

func lookup(protocol string) (TransportConstructor, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	c, ok := registry[protocol]
	return c, ok
}

func Registered() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	protocols := make([]string, 0, len(registry))
	for p := range registry {
		protocols = append(protocols, p)
	}
	return protocols
}
