package version

import (
	"fmt"
	"runtime"
)

// Info contains versioning information.
type Info struct {
	Version   string `json:"version"`
	GitTag    string `json:"gitTag"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info *Info) String() string {
	return info.GitTag
}

// Get 返回详细的版本信息
func Get() Info {
	return Info{
		Version:   version,
		GitTag:    gitTag,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
