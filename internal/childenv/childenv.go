// Package childenv adjusts the environment handed to short-lived Go helper
// child processes (sandbox-ctl per task, node-ctl checkpoint-cleanup per
// pause cleanup).
//
// Go initializes one scheduler P per logical CPU during schedinit, before
// any user code runs. On many-core hosts (e.g. 384 cores) every short-lived
// helper pays 15-20 ms of runtime initialization that no amount of in-process
// tuning can recover, because P structures are already allocated by the time
// main executes. Capping the scheduler width at exec time removes that start
// tax; see kuasar-sandbox issue #6 for the measurements.
//
// Long-running services (conductor, runner, proxy) keep the default
// GOMAXPROCS: their concurrency needs are unrelated, and a global env var
// would couple them to this start-up optimization.
package childenv

import "strings"

// GOMAXPROCS caps a short-lived Go helper's scheduler width unless the parent
// environment already sets GOMAXPROCS. Explicit settings (parent env or unit
// spec carried through taskEnv) always win, so per-host tuning stays possible.
func GOMAXPROCS(env []string, value string) []string {
	for _, kv := range env {
		if strings.HasPrefix(kv, "GOMAXPROCS=") {
			return env
		}
	}
	return append(env, "GOMAXPROCS="+value)
}
