package runtime

import _ "embed"

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-gcp
var runtime []byte

func NitricGcpRuntime() []byte {
	return runtime
}
