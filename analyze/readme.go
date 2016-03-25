package analyze

import "github.com/spf13/cobra"

var (
	ReadmeCommand = &cobra.Command{
		Use:   "readme",
		Short: "Generates README markdown.",
		RunE:  ReadmeCommandFunc,
	}

	readmeDir string
)

func init() {
	ReadmeCommand.PersistentFlags().StringVarP(&readmeDir, "readme-dir", "d", "", "Directory path to generate README.")
}

func ReadmeCommandFunc(cmd *cobra.Command, args []string) error {
	return nil
}
