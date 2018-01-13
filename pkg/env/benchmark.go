package env

import "os"

var (
	// ScanID is used to provide the benchmark services with the current scan
	ScanID = Setting(scanID{})

	// Checks is used to provide the benchmark services with the checks that need to be run as part of the benchmark
	Checks = Setting(checks{})

	// BenchmarkName is used to provide the benchmark service with the benchmark name
	BenchmarkName = Setting(benchmarkName{})

	// BenchmarkCompletion is used to provide the benchmark service with whether or not the benchmark container should exit
	BenchmarkCompletion = Setting(benchmarkCompletion{})
)

type scanID struct{}

func (s scanID) EnvVar() string {
	return "ROX_APOLLO_SCAN_ID"
}

func (s scanID) Setting() string {
	return os.Getenv(s.EnvVar())
}

type checks struct{}

func (c checks) EnvVar() string {
	return "ROX_APOLLO_CHECKS"
}

func (c checks) Setting() string {
	return os.Getenv(c.EnvVar())
}

type benchmarkName struct{}

func (c benchmarkName) EnvVar() string {
	return "ROX_APOLLO_BENCHMARK_NAME"
}

func (c benchmarkName) Setting() string {
	return os.Getenv(c.EnvVar())
}

type benchmarkCompletion struct{}

func (c benchmarkCompletion) EnvVar() string {
	return "ROX_APOLLO_BENCHMARK_COMPLETION"
}

func (c benchmarkCompletion) Setting() string {
	return os.Getenv(c.EnvVar())
}
