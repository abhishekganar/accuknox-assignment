Problem 1 – Drop TCP Traffic on a Specific Port
Objective

Drop TCP packets destined for a specific port using eBPF/XDP.

The default blocked port is `4040`.

### Implementation

The solution consists of two parts:

- eBPF program (`problem1/ebpf/xdp_filter.c`)
  - Runs in the Linux kernel using XDP.
  - Inspects incoming Ethernet, IP, and TCP headers.
  - Checks the destination TCP port.
  - Returns `XDP_DROP` when the packet matches the configured port.
  - Returns `XDP_PASS` for other traffic.

- Go userspace program (`problem1/main.go`)
  - Loads the compiled eBPF program.
  - Configures the blocked port through a BPF map.
  - Attaches the XDP program to the network interface.
  - Keeps the program running until it receives `Ctrl+C`.

How to Run

From the `problem1` directory:

```bash
go run main.go
