package java

import (
	"fmt"
	"time"

	"github.com/grafana/alloy/internal/component/discovery"
	"github.com/grafana/alloy/internal/component/pyroscope"
)

type Arguments struct {
	Targets   []discovery.Target     `alloy:"targets,attr"`
	ForwardTo []pyroscope.Appendable `alloy:"forward_to,attr"`

	TmpDir          string          `alloy:"tmp_dir,attr,optional"`
	ProfilingConfig ProfilingConfig `alloy:"profiling_config,block,optional"`

	// undocumented
	Dist string `alloy:"dist,attr,optional"`
}

type ProfilingConfig struct {
	Interval time.Duration `alloy:"interval,attr,optional"`

	// CPU Profiling Options
	CPU        bool   `alloy:"cpu,attr,optional"`         // Enable CPU profiling
	Event      string `alloy:"event,attr,optional"`       // -e: Profiling event (cpu, itimer, wall, etc.)
	SampleRate int    `alloy:"sample_rate,attr,optional"` // Sample rate for CPU profiling (samples per second)
	Wall       string `alloy:"wall,attr,optional"`        // --wall: Wall clock profiling interval
	AllUser    bool   `alloy:"all_user,attr,optional"`    // --all-user: Include only user-mode events
	PerThread  bool   `alloy:"per_thread,attr,optional"`  // -t: Profile threads separately
	Filter     string `alloy:"filter,attr,optional"`      // --filter: Profile only threads with specified ids
	Sched      bool   `alloy:"sched,attr,optional"`       // --sched: Group threads by scheduling policy
	TTSP       bool   `alloy:"ttsp,attr,optional"`        // --ttsp: Time-to-safepoint profiling
	Begin      string `alloy:"begin,attr,optional"`       // --begin: Auto-start profiling when function is executed
	End        string `alloy:"end,attr,optional"`         // --end: Auto-stop profiling when function is executed
	NoStop     bool   `alloy:"nostop,attr,optional"`      // --nostop: Don't stop profiling outside --begin/--end window
	Proc       string `alloy:"proc,attr,optional"`        // --proc: Collect system process statistics interval
	TargetCPU  int    `alloy:"target_cpu,attr,optional"`  // --target-cpu: Sample only threads on specified CPU
	RecordCPU  bool   `alloy:"record_cpu,attr,optional"`  // --record-cpu: Capture which CPU a sample was taken on

	// Memory Profiling Options
	Alloc     string `alloy:"alloc,attr,optional"`      // --alloc: Allocation profiling interval
	Live      bool   `alloy:"live,attr,optional"`       // --live: Retain only live objects (not collected)
	NativeMem string `alloy:"native_mem,attr,optional"` // --nativemem: Native memory allocation profiling interval
	NoFree    bool   `alloy:"no_free,attr,optional"`    // --nofree: Don't record free calls in native memory profiling

	// Lock Profiling Options
	Lock       string `alloy:"lock,attr,optional"`        // --lock: Lock profiling threshold
	NativeLock string `alloy:"native_lock,attr,optional"` // --nativelock: Native lock (pthread) profiling threshold

	// Miscellaneous Options
	All         bool     `alloy:"all,attr,optional"`         // --all: Enable cpu, wall, alloc, live, nativemem, lock simultaneously
	LogLevel    string   `alloy:"log_level,attr,optional"`   // -L: Log level (debug, info, warn, error, none)
	Quiet       bool     `alloy:"quiet,attr,optional"`       // Suppress profiler output
	Include     []string `alloy:"include,attr,optional"`     // -I: Include stack trace patterns
	Exclude     []string `alloy:"exclude,attr,optional"`     // -X: Exclude stack trace patterns
	JStackDepth int      `alloy:"jstackdepth,attr,optional"` // -j: Maximum Java stack depth
	CStack      string   `alloy:"cstack,attr,optional"`      // --cstack: C stack walking mode (fp, dwarf, lbr, vm, vmx, no)
	Features    []string `alloy:"features,attr,optional"`    // -F: Stack walking features (stats, vtable, comptask, pcaddr)
	Trace       []string `alloy:"trace,attr,optional"`       // --trace: Java methods to trace with optional latency threshold
	JFRSync     string   `alloy:"jfrsync,attr,optional"`     // --jfrsync: Sync Java Flight Recording config
	Signal      string   `alloy:"signal,attr,optional"`      // --signal: Alternative signal for profiling
	Clock       string   `alloy:"clock,attr,optional"`       // --clock: Clock source for JFR timestamps (tsc, monotonic)
}

func (rc *Arguments) UnmarshalAlloy(f func(any) error) error {
	*rc = DefaultArguments()
	type config Arguments
	return f((*config)(rc))
}

func (arg *Arguments) Validate() error {
	switch arg.ProfilingConfig.Event {
	case "itimer", "cpu", "wall":
		return nil
	default:
		return fmt.Errorf("invalid event: '%s'. Event must be one of 'itimer', 'cpu' or 'wall'", arg.ProfilingConfig.Event)
	}
}

func DefaultArguments() Arguments {
	return Arguments{
		TmpDir: "/tmp",
		ProfilingConfig: ProfilingConfig{
			Interval: 60 * time.Second,

			// CPU Profiling Options
			CPU:        true,
			Event:      "itimer",
			SampleRate: 100,
			Wall:       "",
			AllUser:    false,
			PerThread:  false,
			Filter:     "",
			Sched:      false,
			TTSP:       false,
			Begin:      "",
			End:        "",
			NoStop:     false,
			Proc:       "",
			TargetCPU:  -1,
			RecordCPU:  false,

			// Memory Profiling Options
			Alloc:     "512k",
			Live:      false,
			NativeMem: "",
			NoFree:    false,

			// Lock Profiling Options
			Lock:       "10ms",
			NativeLock: "",

			// Miscellaneous Options
			All:         false,
			LogLevel:    "INFO",
			Quiet:       false,
			Include:     []string{},
			Exclude:     []string{},
			JStackDepth: 2048,
			CStack:      "",
			Features:    []string{},
			Trace:       []string{},
			JFRSync:     "",
			Signal:      "",
			Clock:       "tsc",
		},
	}
}
