package tag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

func latestTag(username string, password string, environment string, registry string,
	repository string, debug bool) {
	var regex string
	var currentVersion *version.Version

	di := DockerImage{
		Username:   username,
		Password:   password,
		Registry:   registry,
		Repository: repository,
	}

	switch environment {
	case "dev":
		regex = DEV_REGEX
	case "staging":
		regex = STAGING_REGEX
	case "prod":
		regex = PROD_REGEX
	}

	tags, err := getTags(&di, regex)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Sort tags
	sort.Sort(version.Collection(tags))

	// No versions yet
	if len(tags) == 0 {
		switch environment {
		case "dev":
			currentVersion, _ = version.NewVersion("v0.0.0-dev.0")
		case "staging":
			currentVersion, _ = version.NewVersion("v0.0.0-rc.0")
		case "prod":
			currentVersion, _ = version.NewVersion("v0.0.0")
		}
	} else {
		currentVersion = tags[len(tags)-1]
	}

	fmt.Println(currentVersion.Original())
}

func nextTag(username string, password string,
	tagType string, environment string, registry string,
	repository string, debug bool) {
	var regex string
	var currentVersion *version.Version

	di := DockerImage{
		Username:   username,
		Password:   password,
		Registry:   registry,
		Repository: repository,
	}

	switch environment {
	case "dev":
		regex = DEV_REGEX
	case "staging":
		regex = STAGING_REGEX
	case "prod":
		regex = PROD_REGEX
	}

	tags, err := getTags(&di, regex)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Sort tags
	sort.Sort(version.Collection(tags))

	// No versions yet
	if len(tags) == 0 {
		switch environment {
		case "dev":
			currentVersion, _ = version.NewVersion("v0.0.0-dev.0")
		case "staging":
			currentVersion, _ = version.NewVersion("v0.0.0-rc.0")
		case "prod":
			currentVersion, _ = version.NewVersion("v0.0.0")
		}
	} else {
		currentVersion = tags[len(tags)-1]
	}

	nextVersion, err := getNextVersion(currentVersion, tagType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Debugf("Tags: %v\n", tags)
	log.Debugf("Current version: %s\n", currentVersion.Original())
	fmt.Println(nextVersion.Original())
}

func getNextVersion(v *version.Version, tagType string) (*version.Version, error) {
	var vstr string
	switch tagType {
	case "major":
		rMajor := regexp.MustCompile(`^v([0-9]+)(.*)$`)
		new := v.Segments64()[0] + 1
		vstr = rMajor.ReplaceAllString(v.Original(), fmt.Sprintf("v%d${2}", new))
	case "minor":
		rMinor := regexp.MustCompile(`^(v[0-9]+\.)([0-9]+)(.*)$`)
		new := v.Segments64()[1] + 1
		vstr = rMinor.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d${3}", new))
	case "patch":
		rPatch := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.)([0-9]+)(.*)$`)
		new := v.Segments64()[2] + 1
		vstr = rPatch.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d${3}", new))
	case "rc":
		rRc := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+-rc\.)([0-9]+)$`)
		match := rRc.FindStringSubmatch(v.Original())
		n, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		vstr = rRc.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d", n+1))
	case "dev":
		rDev := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+-dev\.)([0-9]+)$`)
		match := rDev.FindStringSubmatch(v.Original())
		n, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		vstr = rDev.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d", n+1))
	}

	nextVersion, err := version.NewVersion(vstr)
	if err != nil {
		return nil, err
	}

	return nextVersion, nil
}

func getTags(di *DockerImage, regex string) ([]*version.Version, error) {

	r, err := getRepo(di)
	tags := []*version.Version{}
	if err != nil {
		return nil, err
	}

	for _, tag := range r.Tags {
		matched, _ := regexp.MatchString(regex, tag.Name)
		if matched {
			v, _ := version.NewVersion(tag.Name)
			tags = append(tags, v)
		}
	}

	return tags, nil
}

func getRepo(di *DockerImage) (*Repository, error) {

	repo := Repository{}
	url := fmt.Sprintf("%s/acr/v1/%s/_tags", di.Registry, di.Repository)
	r, err := doGet(url, nil, di.Username, di.Password)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func doGet(url string, querystring map[string]string, username string, password string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	q := req.URL.Query()
	// Set custom querystring pairs
	for k, v := range querystring {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	log.Debug("GET: ", req.URL)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return r, nil
}
