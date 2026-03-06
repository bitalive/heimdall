# Heimdall Engine (Divine Core) 🛡️

**Heimdall** is a high-performance, nanosecond-precision network engine built for Golang. Designed as the base layer for **Bitalive** and **Conceit**, it provides an "Invisible I/O" transport that bypasses Go runtime overhead to achieve extreme packet processing rates.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg?style=for-the-badge)](LICENSE)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-orange.svg?style=for-the-badge)](LICENSE-BSD)

---

## ⚡ Performance
Heimdall is built for one purpose: **Speed**.

| Metric | Result | Environment |
| :--- | :--- | :--- |
| **Integrated Throughput** | **102,739,064 PPS** | Single Core (Intel i5-9400 @ 2.9GHz) |
| **Memory Allocation** | **0 allocs/op** | Zero-copy / Zero-GC pressure |
| **GetFrame Latency** | **< 0.9 ns** | Overlay Memory Mapping |
| **Syscall Efficiency** | **Batching (Up to 256)** | recvmmsg / sendmmsg |

---

## ✨ Features

- **Invisible I/O**: Custom x86_64 Assembly SYSCALL implementation that bypasses `runtime.entersyscall` and `runtime.exitsyscall`.
- **Zero-Allocation Hot Path**: Uses **Chronos Arena** for all internal buffers. No Go heap allocations during Poll/Send cycles.
- **Kernel-Level Resilience**: Built-in `EINTR` retry logic directly in the Assembly layer.
- **Multi-Backend Fallback**: High-performance Assembly on Linux/AMD64 with automatic fallback to standard `net.UDPConn` on macOS, Windows, and other architectures.
- **Micro-Telemetry**: Built-in atomic counters for real-time monitoring of PPS, Bytes, Syscalls, and Errors with zero performance impact.

---

## 🚀 Installation

```bash
go get github.com/bitalive/heimdall
```

---

## 🛠️ Usage

Heimdall is designed to be simple and professional.

```go
package main

import (
    "fmt"
    "github.com/bitalive/heimdall/io"
)

func main() {
    // 1. Configure the engine
    cfg := io.Config{
        BatchSize: 128,
        MsgSize:   64,
    }

    // 2. Initialize (Heimdall handles Arena and Sockets internally)
    engine, err := io.NewEngine(9999, cfg)
    if err != nil {
        panic(err)
    }

    fmt.Println("Heimdall Engine Active on Port 9999")

    for {
        // 3. Poll for a batch of packets
        n := engine.Poll()
        
        for i := 0; i < n; i++ {
            // 4. GetFrame (Zero-copy []byte view of the Arena)
            data := engine.GetFrame(i)
            
            // Handle your business logic here...
            _ = data
        }

        // 5. Send responses in batch
        engine.Send(n)
        
        // Optional: Check real-time stats
        // stats := engine.GetStats()
    }
}
```

---

## 📚 Architecture

Heimdall works as a **"Network Gearbox"**:
1. **Gathers**: Uses `recvmmsg` (Assembly) to harvest 128+ packets in one syscall.
2. **Stores**: Kernel writes directly into a pre-allocated **Chronos Arena**.
3. **Presents**: Provides zero-copy `[]byte` slices to your application.
4. **Dispatches**: Flushes responses using `sendmmsg` (Assembly).

---

## 🛡️ License

This project is dual-licensed:
- **GNU General Public License v3.0** (GPL-3.0)
- **3-Clause BSD License**

Choose the license that best fits your project requirements.

---

## 🤝 Contribution

Heimdall is an infrastructure component for the **Bitalive Ecosystem**. We welcome contributions that focus on:
- AVX-512 Protocol Parsing.
- Persistent Storage primitives.
- io_uring event loop integration.

---
Copyright © 2026 **Bitalive Team**. Built for the age of nanoseconds.
