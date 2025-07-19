package cmd

import "github.com/spf13/cobra"

var (
	Create bool
	Push bool
	NumOfWorkers int
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "gdriveSync",
	Short: "Main command for the application",
	Long:  `This command serves as the entry point for the gdriveSync application.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Push, "push", "p", false, "Push a new changes")
	rootCmd.PersistentFlags().BoolVarP(&Create, "create", "c", false, "Create a new GDrive init file")
	rootCmd.PersistentFlags().IntVarP(&NumOfWorkers, "workers", "w", 10, "Number of workers to use for uploading files")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose logging")
}

func Execute() error {
	return rootCmd.Execute()
}
