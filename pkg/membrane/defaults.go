package membrane

import (
	"log"
	"os"
	"strings"

	"github.com/nitrictech/nitric/pkg/utils"
)

func DefaultMembraneOptions() *MembraneOptions {
	options := &MembraneOptions{}

	if len(os.Args) > 1 && len(os.Args[1:]) > 0 {
		options.ChildCommand = os.Args[1:]
	} else {
		options.ChildCommand = strings.Fields(utils.GetEnv("INVOKE", ""))
		if len(options.ChildCommand) > 0 {
			log.Default().Println("Warning: use of INVOKE environment variable is deprecated and may be removed in a future version")
		}
	}

	return options
}
