// Package nocgroup implements the KUASAR_NO_CGROUP temporary operating mode.
//
// KUASAR_NO_CGROUP=1 makes the node execution path skip cgroup v2 setup:
// run-sandbox/run-builder do not prepare a delegated VMM cgroup, sandbox-ctl is
// started without --cgroup-path (the runtime's documented no-cgroup mode), and
// the reconciler omits cgroup-only resource policy (overhead, watermark_high)
// while pinning allocatable to capacity.
//
// This is a development/portability escape hatch for hosts that only provide
// cgroup v1. It disables per-sandbox CPU/memory enforcement, memory.high budget
// enforcement, memory accounting and OOM deflation; node admission is not
// enforced. Production deployments must leave it unset.
package nocgroup

import "os"

// EnvVar is the environment variable that enables the mode.
const EnvVar = "KUASAR_NO_CGROUP"

// Enabled reports whether the no-cgroup mode is active.
func Enabled() bool {
	return os.Getenv(EnvVar) == "1"
}
