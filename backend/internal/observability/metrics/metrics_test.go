// SPDX-License-Identifier: MIT

package metrics

import (
	"testing"
)

func TestNoOp_Timer(t *testing.T) {
	m := noOpMetrics{}
	end := m.Timer("test")
	end()
}

func TestNoOp_Counter(t *testing.T) {
	m := noOpMetrics{}
	m.Counter("test", 1)
}

func TestNoOp_Gauge(t *testing.T) {
	m := noOpMetrics{}
	m.Gauge("test", 1.0)
}

func TestGlobal_DefaultNoOp(t *testing.T) {
	m := Get()
	if m == nil {
		t.Fatal("Get() must not return nil")
	}
	m.Timer("test")()
	m.Counter("test", 1)
	m.Gauge("test", 1.0)
}

func TestSet_ReplaceGlobal(t *testing.T) {
	Set(noOpMetrics{})
	if Get() == nil {
		t.Fatal("Get() must not return nil after Set")
	}
}

func TestSet_NilResetsNoOp(t *testing.T) {
	Set(nil)
	m := Get()
	if m == nil {
		t.Fatal("Get() must not return nil after Set(nil)")
	}
	m.Timer("test")()
}
