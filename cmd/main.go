package cmd

import "github.com/spf13/cobra"

var (
	Create bool
	Push bool
	NumOfWorkers int
	Verbose bool
)

const (
	ALL = iota
	DRIVE
	SHARED
)
const (
	BOTH = iota
	FOLDER
	FILE
)

var (
	Show bool
	Type int
	FileType int
)

var (
	Search = false
	SearchQuery string
)

var rootCmd = &cobra.Command{
	Use:   "gdriveSync",
	Short: "Main command for the application",
	Long:  `This command serves as the entry point for the gdriveSync application.`,
}

var listcmd = &cobra.Command{
	Use:   "list",
	Short: "List Google Drive files",
	Long:  `This command lists all files in your Google Drive.`,
	Run: func(cmd *cobra.Command, args []string) {
		Show = true
		FileType = FOLDER
	},
}

var searchcmd = &cobra.Command{
	Use:   "search",
	Short: "Search Google Drive files",
	Long:  `This command allows you to search for files in your Google Drive.`,
	Run: func(cmd *cobra.Command, args []string) {
		Search = true
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Push, "push", "p", false, "Push a new changes")
	rootCmd.PersistentFlags().BoolVarP(&Create, "create", "c", false, "Create a new GDrive init file")
	rootCmd.PersistentFlags().IntVarP(&NumOfWorkers, "workers", "w", 10, "Number of workers to use for uploading files")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose logging")
	listcmd.PersistentFlags().IntVarP(&Type, "type", "t", DRIVE, "Type of files to show (0: All, 1: Drive, 2: Shared)")
	listcmd.PersistentFlags().IntVarP(&FileType, "file-type", "x", FOLDER, "Type of files to show (0: Both, 1: Folder, 2: File)")
	searchcmd.PersistentFlags().IntVarP(&Type, "type", "t", DRIVE, "Type of files to search (0: All, 1: Drive, 2: Shared)")
	searchcmd.PersistentFlags().IntVarP(&FileType, "file-type", "x", BOTH, "Type of files to search (0: Both, 1: Folder, 2: File)")
	searchcmd.PersistentFlags().StringVarP(&SearchQuery, "query", "q", "", "Search query")
	rootCmd.AddCommand(listcmd)
	rootCmd.AddCommand(searchcmd)
}

func Execute() error {
	return rootCmd.Execute()
}
