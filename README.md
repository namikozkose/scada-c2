# SCADA-C2 🎛️🌐

A centralized Command & Control (C2) / SOC Server for Distributed ICS Honeypots, written in Go.

I developed this central server to manage multiple edge defense agents (like my ModbusTrap project) deployed across an Operational Technology (OT) network. When an edge node detects industrial sabotage, it securely transmits the threat telemetry to this C2 server via a JSON-based REST API.

## ⚙️ Core Architecture

*   **High-Concurrency REST API:** Built using Go's native `net/http` package to handle simultaneous incoming alerts from hundreds of edge sensors without blocking or dropping packets.
*   **Thread-Safe Processing:** Utilizes `sync.Mutex` to prevent race conditions in memory when multiple agents report attacks at the exact same millisecond.
*   **Standardized Threat Telemetry:** Parses incoming JSON payloads containing critical incident data: Agent ID, Attacker IP, specific Modbus Function Codes (e.g., `0x05` for sabotage), and raw hexadecimal dumps for forensic analysis.

## 🚀 Deployment

```bash
# Clone the repository
git clone [https://github.com/namikozkose/scada-c2.git](https://github.com/namikozkose/scada-c2.git)
cd scada-c2

# Build the SOC server
go build -o scada-c2 main.go

# Start the server (Listens on port 8080)
./scada-c2