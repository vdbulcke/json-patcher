package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"github.com/vdbulcke/json-patcher/src/patcher"
	"github.com/vdbulcke/json-patcher/src/tui"
	"go.uber.org/zap"
)

func init() {
	// bind to root command
	rootCmd.AddCommand(interactive)
	// add flags to sub command
	interactive.Flags().StringVarP(&patchFile, "patch-file", "p", "", "file containing a list of patches")
	interactive.Flags().StringVarP(&skipTags, "skip-tags", "s", "", "comma separated list of tags to skip")
	interactive.Flags().BoolVarP(&allowUnescapedHTML, "allow-unescaped-html", "", false, "allow unescaped HTML in JSON output")

	// required flags
	//nolint
	interactive.MarkFlagRequired("patch-file")

}

var interactive = &cobra.Command{
	Use:   "interactive",
	Short: "interactive  list of patches",
	// Long: "",
	Run: interactiveCmd,
}

// interactiveCmd
func interactiveCmd(cmd *cobra.Command, args []string) {

	logger := logger.GetZapLogger(Debug)

	c, err := config.ParseConfig(patchFile)
	if err != nil {
		logger.Error("Error parsing config", zap.String("filename", patchFile), zap.Error(err))
		os.Exit(1)
	}

	err = tui.StartUI(c, patcher.NewOptions(skipTags, allowUnescapedHTML, Debug))
	if err != nil {
		logger.Error("Error applying patches", zap.Error(err))
		os.Exit(1)
	}

}
