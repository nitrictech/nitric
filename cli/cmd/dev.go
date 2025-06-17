package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/simulation"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func (p *PrefixWriter) Write(content []byte) (int, error) {
	value := strings.TrimSuffix(string(content), "\n")

	split := strings.Split(value, "\n")
	value = strings.Join(split, "\n"+p.prefix) + "\n"

	_, err := fmt.Fprintf(p.writer, "%s%s", p.prefix, value)
	if err != nil {
		return 0, err
	}

	return len(content), nil
}

func NewPrefixWriter(prefix string, writer io.Writer) *PrefixWriter {
	return &PrefixWriter{
		prefix: prefix,
		writer: writer,
	}
}

var dev = &cobra.Command{
	Use:   "dev",
	Short: "Run the Nitric application in development mode",
	Long:  `Run the Nitric application in development mode, allowing local testing of resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Extract common loading logic into a separate function
		// (see build.go for common loading logic)

		// 1. Load the App Spec
		// Read the nitric.yaml file
		fs := afero.NewOsFs()

		appSpec, err := schema.LoadFromFile(fs, "nitric.yaml")
		cobra.CheckErr(err)

		simserver := simulation.NewSimulationServer(fs, appSpec)
		err = simserver.Start(os.Stdout)
		cobra.CheckErr(err)
	},
}

func init() {
	// Add the dev command to the root command
	rootCmd.AddCommand(dev)

	// Add flags for the dev command if needed
	// e.g., dev.Flags().StringVarP(&flagName, "flag", "f", "defaultValue", "Description of the flag")
}

//
