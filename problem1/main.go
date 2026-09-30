package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type bpfObjects struct {
	XdpFilter   *ebpf.Program `ebpf:"xdp_filter"`
	BlockedPort *ebpf.Map     `ebpf:"blocked_port"`
}

func main() {
	// Load the compiled eBPF object
	spec, err := ebpf.LoadCollectionSpec("ebpf/xdp_filter.o")
	if err != nil {
		log.Fatalf("loading eBPF object: %v", err)
	}

	// Load the eBPF program and map
	var objs bpfObjects

	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("loading eBPF program: %v", err)
	}
	defer objs.XdpFilter.Close()
	defer objs.BlockedPort.Close()

	// Port we want to block
	port := uint16(4040)
	key := uint32(0)

	// Put the port into the BPF map
	if err := objs.BlockedPort.Update(key, port, ebpf.UpdateAny); err != nil {
		log.Fatalf("updating blocked port: %v", err)
	}

	// Find network interface
	iface, err := net.InterfaceByName("enp0s3")
	if err != nil {
		log.Fatalf("finding interface: %v", err)
	}

	// Attach eBPF program to the network interface
	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpFilter,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("attaching XDP: %v", err)
	}
	defer xdpLink.Close()

	fmt.Printf("XDP filter attached to %s\n", iface.Name)
	fmt.Printf("Blocking TCP port %d\n", port)
	fmt.Println("Press Ctrl+C to stop.")

	// Keep program running
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
