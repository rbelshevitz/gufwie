package version

// These variables are intended to be overridden at build time via -ldflags "-X".
var (
	Version = "0.2"
	Commit  = "dev"
	Date    = "unknown"
)
