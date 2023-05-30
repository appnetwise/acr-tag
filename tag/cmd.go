package tag

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/morikuni/aec"
)

var (
	app_version = "dev"
	commit      = "n/a"
	date        = "n/a"
)

const asciiYak = `
____  _____  _______\_ |__    ____    ____  ______  |  |  _____     ____   ____   
_/ ___\ \__  \ \_  __ \| __ \  /  _ \  /    \ \____ \ |  |  \__  \  _/ ___\_/ __ \  
\  \___  / __ \_|  | \/| \_\ \(  <_> )|   |  \|  |_> >|  |__ / __ \_\  \___\  ___/  
 \___  >(____  /|__|   |___  / \____/ |___|  /|   __/ |____/(____  / \___  >\___  > 
     \/      \/            \/              \/ |__|               \/      \/     \/  

`

func VersionCmd() {
	yak := aec.RedF.Apply(asciiYak)
	fmt.Print(yak)
	fmt.Println(aec.LightYellowF.Apply(fmt.Sprintf("Version %s (commit %s, %s)\n", app_version, commit, date)))
}

func LatestCmd(username string, password string, environment string,
	registry string, repository string, debug bool) {

	e := Environment(environment)
	i := ImageTarget{
		Username:   username,
		Password:   password,
		Registry:   registry,
		Repository: repository,
	}

	// Validate Image Target
	if err := i.Validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Validate environment
	if err := e.IsValid(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	latestTag(e, i, debug)
}

func NextCmd(username string, password string, tagType string, environment Environment,
	registry string, repository string, debug bool, version string) {

	e := Environment(environment)
	t := TagType(tagType)
	i := ImageTarget{
		Username:   username,
		Password:   password,
		Registry:   registry,
		Repository: repository,
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Check for static version check
	if version == "" {

		// Validate environment
		if err := e.IsValid(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Check valid tagType is passed to environment
		if err := e.ValidateTag(t); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Validate Image Target
		if err := i.Validate(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Query registry for next tag
		nextTag(e, i, t, debug)

	} else {

		// Check valid tagType is passed to command
		if err := t.IsValid(); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Perform local calculation based on version
		nextTagFromVersionString(version, t)

	}

}
