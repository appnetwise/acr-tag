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

func latestTag(e Environment, i ImageTarget, debug bool) {

	var currentVersion *version.Version

	if rx, err := e.Regex(); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		tags, err := getTags(&i, rx)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Sort tags
		sort.Sort(version.Collection(tags))

		// No versions yet
		if len(tags) == 0 {
			currentVersion, _ = version.NewVersion(func() string {
				v, _ := e.DefaultVersion()
				return v
			}())
		} else {
			currentVersion = tags[len(tags)-1]
		}

	}

	fmt.Println(currentVersion.Original())
}

func nextTag(e Environment, i ImageTarget, t TagType, debug bool) {

	var currentVersion *version.Version

	tags, err := getTags(&i, func() string {
		r, _ := e.Regex()
		return r
	}())

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Sort tags
	sort.Sort(version.Collection(tags))
	log.Debugf("Tags: %v\n", tags)

	// No versions yet
	if len(tags) == 0 {
		currentVersion, _ = version.NewVersion(func() string {
			v, _ := e.DefaultVersion()
			return v
		}())
	} else {
		currentVersion = tags[len(tags)-1]
	}

	log.Debugf("Current version: %s\n", currentVersion.Original())

	nextVersion, err := getNextVersion(currentVersion, t)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	fmt.Println(nextVersion.Original())
}

func getNextVersion(v *version.Version, tagType TagType) (*version.Version, error) {
	var vstr string
	switch tagType {
	case TAG_MAJOR:
		rMajor := regexp.MustCompile(`^v([0-9]+)(.*)$`)
		new := v.Segments64()[0] + 1
		vstr = rMajor.ReplaceAllString(v.Original(), fmt.Sprintf("v%d${2}", new))
	case TAG_MINOR:
		rMinor := regexp.MustCompile(`^(v[0-9]+\.)([0-9]+)(.*)$`)
		new := v.Segments64()[1] + 1
		vstr = rMinor.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d${3}", new))
	case TAG_PATCH:
		rPatch := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.)([0-9]+)(.*)$`)
		new := v.Segments64()[2] + 1
		vstr = rPatch.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d${3}", new))
	case TAG_RC:
		rRc := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+-rc\.)([0-9]+)$`)
		match := rRc.FindStringSubmatch(v.Original())
		n, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		vstr = rRc.ReplaceAllString(v.Original(), fmt.Sprintf("${1}%d", n+1))
	case TAG_DEV:
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

func getTags(i *ImageTarget, regex string) ([]*version.Version, error) {

	log.Debugf("Searching for tags using regex: %s\n", regex)

	r, err := getRepo(i)
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

func getRepo(i *ImageTarget) (*Repository, error) {

	repo := Repository{}
	url := fmt.Sprintf("%s/acr/v1/%s/_tags", i.Registry, i.Repository)
	r, err := doGet(url, nil, i.Username, i.Password)
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
