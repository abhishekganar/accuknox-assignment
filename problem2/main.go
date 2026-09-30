package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type bpfObjects struct {
	ProcessFilter *ebpf.Program `ebpf:"process_filter"`
}

func main() {
	// Load the compiled eBPF object
	spec, err := ebpf.LoadCollectionSpec("ebpf/process_filter.o")
	if err != nil {
		log.Fatalf("loading eBPF object: %v", err)
	}

	// Load the eBPF program
	var objs bpfObjects

	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("loading eBPF program: %v", err)
	}
	defer objs.ProcessFilter.Close()

	// Open the cgroup
	cgroup, err := os.Open("/sys/fs/cgroup")
	if err != nil {
		log.Fatalf("opening cgroup: %v", err)
	}
	defer cgroup.Close()

	// Attach eBPF program to cgroup/connect4
	cgLink, err := link.AttachCgroup(link.CgroupOptions{
		Path:    "/sys/fs/cgroup",
		Attach:  ebpf.AttachCGroupInet4Connect,
		Program: objs.ProcessFilter,
	})
	if err != nil {
		log.Fatalf("attaching eBPF program: %v", err)
	}
	defer cgLink.Close()

	fmt.Println("eBPF process filter attached.")
	fmt.Println("Allowing myprocess to use TCP port 4040 only.")
	fmt.Println("Press Ctrl+C to stop.")

	select {}
}
