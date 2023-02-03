package tag

import (
	"errors"
	"time"
)

type ImageTarget struct {
	Username   string
	Password   string
	Registry   string
	Repository string
}

func (i ImageTarget) Validate() error {
	if len(i.Repository) == 0 {
		return errors.New("a target repository has not been specified")
	}
	if len(i.Registry) == 0 {
		return errors.New("the target registry has not been specified")
	}
	return nil
}

type Repository struct {
	Registry  string `json:"registry"`
	ImageName string `json:"imageName"`
	Tags      []Tag  `json:"tags"`
}

type Tag struct {
	Name           string    `json:"name"`
	Digest         string    `json:"digest"`
	CreatedTime    time.Time `json:"createdTime"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}
