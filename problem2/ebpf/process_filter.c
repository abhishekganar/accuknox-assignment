#include <linux/bpf.h>
#include <linux/in.h>
#include "bpf_helpers.h"

SEC("cgroup/connect4")
int process_filter(struct bpf_sock_addr *ctx)
{
    char comm[16];

    bpf_get_current_comm(&comm, sizeof(comm));

    if (comm[0] != 'm' ||
        comm[1] != 'y' ||
        comm[2] != 'p' ||
        comm[3] != 'r' ||
        comm[4] != 'o' ||
        comm[5] != 'c' ||
        comm[6] != 'e' ||
        comm[7] != 's' ||
        comm[8] != 's' ||
        comm[9] != '\0')
    {
        return 1;
    }

    if (ctx->user_port != __constant_htons(4040))
        return 0;

    return 1;
}

char LICENSE[] SEC("license") = "GPL";
