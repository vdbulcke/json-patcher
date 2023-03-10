package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"github.com/vdbulcke/json-patcher/src/patcher"
	"go.uber.org/zap"
)

var patchFile string
var skipTags string
var allowUnescapedHTML bool

func init() {
	// bind to root command
	rootCmd.AddCommand(apply)
	// add flags to sub command
	apply.Flags().StringVarP(&patchFile, "patch-file", "p", "", "file containing a list of patches")
	apply.Flags().StringVarP(&skipTags, "skip-tags", "s", "", "comma separated list of tags to skip")
	apply.Flags().BoolVarP(&allowUnescapedHTML, "allow-unescaped-html", "", false, "allow unescaped HTML in JSON output")

	// required flags
	//nolint
	apply.MarkFlagRequired("patch-file")

}

var apply = &cobra.Command{
	Use:   "apply",
	Short: "apply  list of patches",
	// Long: "",
	Run: applyCmd,
}

// applyCmd
func applyCmd(cmd *cobra.Command, args []string) {

	logger := logger.GetZapLogger(Debug)

	c, err := config.ParseConfig(patchFile)
	if err != nil {
		logger.Error("Error parsing config", zap.String("filename", patchFile), zap.Error(err))
		os.Exit(1)
	}

	err = patcher.Apply(c, patcher.NewOptions(skipTags, allowUnescapedHTML, Debug))
	if err != nil {
		logger.Error("Error applying patches", zap.Error(err))
		os.Exit(1)
	}

}
