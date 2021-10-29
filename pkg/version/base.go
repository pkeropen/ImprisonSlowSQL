package version

var (
	gitTag       = ""
	gitCommit    = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	version      = "1.0.0"     			  // state of git tree, either "clean" or "dirty"
	buildDate    = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)
