package cmd

import "github.com/spf13/cobra"

var (
	Create bool
)

var rootCmd = &cobra.Command{
	Use:   "gdriveSync",
	Short: "Main command for the application",
	Long:  `This command serves as the entry point for the gdriveSync application.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Create, "create", "c", false, "Create a new GDrive init file")
}

func Execute() error {
	return rootCmd.Execute()
}
