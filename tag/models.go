package tag

import "time"

type DockerImage struct {
	Username   string
	Password   string
	Registry   string
	Repository string
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
