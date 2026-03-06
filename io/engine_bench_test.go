package io

import (
	"testing"
)

func BenchmarkEngine_Poll(b *testing.B) {
	cfg := Config{BatchSize: 128, MsgSize: 64}
	// Note: NewEngine creates a real socket, which is fine for measuring overhead.
	// We'll ignore the error if it fails (e.g., in environments without network perms)
	engine, err := NewEngine(0, cfg) // Port 0 for random port
	if err != nil {
		b.Skip("Skipping benchmark: could not create engine")
		return
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		engine.Poll()
	}
}

func BenchmarkEngine_Send(b *testing.B) {
	cfg := Config{BatchSize: 128, MsgSize: 64}
	engine, err := NewEngine(0, cfg)
	if err != nil {
		b.Skip("Skipping benchmark: could not create engine")
		return
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		engine.Send(128)
	}
}

func BenchmarkEngine_GetFrame(b *testing.B) {
	cfg := Config{BatchSize: 128, MsgSize: 64}
	engine, err := NewEngine(0, cfg)
	if err != nil {
		b.Skip("Skipping benchmark: could not create engine")
		return
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = engine.GetFrame(i % 128)
	}
}
