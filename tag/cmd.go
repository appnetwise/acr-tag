package tag

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	app_version = "dev"
	commit      = "none"
	date        = "unknown"
	builtBy     = "unknown"
)

func validateEnvType(environment string, tagType string) error {
	valid := true
	if environment == "dev" && tagType == "rc" {
		valid = false
	}

	if environment == "staging" && tagType == "dev" {
		valid = false
	}

	if environment == "prod" && (tagType == "dev" || tagType == "rc") {
		valid = false
	}
	if !valid {
		return errors.New("invalid environment and tagType pair")
	}

	return nil
}

func VersionCmd() {
	fmt.Printf("Version %s (commit %s)", app_version, commit)
	fmt.Printf("Built at %s by %s", date, builtBy)
}

func LatestCmd(username string, password string, environment string,
	registry string, repository string, debug bool) {

	// Check arguments
	if len(environment) == 0 || len(registry) == 0 ||
		len(repository) == 0 || len(username) == 0 || len(password) == 0 {
		log.Error("Arguments missing")
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	latestTag(username, password, environment, registry, repository, debug)
}

func NextCmd(username string, password string, tagType string, environment string,
	registry string, repository string, debug bool) {

	// Check arguments
	if len(environment) == 0 || len(registry) == 0 ||
		len(repository) == 0 || len(username) == 0 || len(password) == 0 {
		log.Error("arguments missing")
		os.Exit(1)
	}

	// Check valid pair environment, tagType
	err := validateEnvType(environment, tagType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	nextTag(username, password, tagType, environment, registry, repository, debug)
}
