package version

// Current holds the current version of the application.
// It is set at build time using -ldflags.
var Current = "dev"

// Commit holds the git commit hash of the build.
// It is set at build time using -ldflags.
var Commit = "unknown"

// FullVersion returns the full version string including the commit hash.
func FullVersion() string {
	return Current + " (" + Commit + ")"
}
func init() {
	Current = "unknown"
}
