# Problem 2 - Process-Specific TCP Port Filtering with eBPF

## Problem

Write an eBPF program to allow TCP traffic only on a specific port (default 4040) for a given process name (`myprocess`).

All TCP traffic from that process to other ports should be dropped.

## Approach

This implementation uses a Linux cgroup `connect4` eBPF hook.

The eBPF program checks:

1. The name of the process making the connection.
2. The destination TCP port.

The rule is:

- `myprocess` connecting to TCP port `4040` -> allowed
- `myprocess` connecting to any other TCP port -> blocked
- Other processes -> allowed

## Files

### `ebpf/process_filter.c`

The eBPF program runs in the Linux kernel.

It uses `bpf_get_current_comm()` to obtain the current process name and checks `ctx->user_port` to determine the destination port.

For a cgroup connect hook:

- `return 1` allows the connection.
- `return 0` blocks the connection.

### `main.go`

The Go program runs in userspace.

It:

1. Loads the compiled eBPF object.
2. Loads the eBPF program into the kernel.
3. Attaches it to the cgroup using the `connect4` hook.
4. Keeps the program running until it is stopped.

## Build

Compile the eBPF program:

```bash
cd ebpf

clang -O2 -g -target bpf \
-I/root/go/pkg/mod/github.com/cilium/ebpf@v0.22.0/examples/headers \
-c process_filter.c -o process_filter.o

