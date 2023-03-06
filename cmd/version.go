package cmd

import (
	"os"
	"time"

	"github.com/carbonplace/acr-tag/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var versionCmd = &cobra.Command{
	Use:              "version",
	Args:             cobra.NoArgs,
	TraverseChildren: true,
	Short:            "Shows the version information",
	Long:             `Shows the version information for this utility`,
	Run: func(cmd *cobra.Command, args []string) {
		tag.VersionCmd()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = false
	formatter.ForceColors = true
	formatter.TimestampFormat = time.RFC1123

	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
