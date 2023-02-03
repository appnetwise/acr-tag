package cmd

import (
	"os"
	"time"

	"github.com/blakyaks/acr-tag/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var latestCmd = &cobra.Command{
	Use:              "latest",
	Args:             cobra.NoArgs,
	TraverseChildren: true,
	Short:            "Retrieve the latest tag",
	Long:             `Retrieve all tags for a given image and print the latest one.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		environment, _ := cmd.Flags().GetString("environment")
		registry, _ := cmd.Flags().GetString("registry")
		repository, _ := cmd.Flags().GetString("repository")
		debug, _ := cmd.Flags().GetBool("debug")
		tag.LatestCmd(username, password, environment, registry, repository, debug)
	},
}

func init() {
	latestCmd.Flags().StringP("username", "u", "", "Username to authenticate to the registry")
	latestCmd.Flags().StringP("password", "p", "", "Password to authenticate to the registry")
	latestCmd.Flags().StringP("environment", "e", "", "Environment [dev, staging, prod]")
	latestCmd.Flags().StringP("registry", "r", "", "Azure Container Registry URL")
	latestCmd.Flags().StringP("repository", "i", "", "Repository Image Name")
	latestCmd.Flags().BoolP("debug", "", false, "Debug")
	rootCmd.AddCommand(latestCmd)

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
