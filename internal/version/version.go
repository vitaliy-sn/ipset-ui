package version

// Values are set at build time via -ldflags "-X".
var (
	BuildTime = "unknown"
	Commit    = "unknown"
)
