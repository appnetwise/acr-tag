package tag

import "errors"

type Environment string

const (
	ENV_DEMO    Environment = "demo"
	ENV_DEV     Environment = "dev"
	ENV_QA      Environment = "qa"
	ENV_UAT     Environment = "uat"
	ENV_STAGING Environment = "staging"
	ENV_PROD    Environment = "prod"
)

func (e Environment) Regex() (string, error) {
	switch e {
	case ENV_DEMO:
		return `^v[0-9]+\.[0-9]+.[0-9]+-demo\.[0-9]+$`, nil	
	case ENV_DEV:
		return `^v[0-9]+\.[0-9]+.[0-9]+-dev\.[0-9]+$`, nil
	case ENV_QA:
		return `^v[0-9]+\.[0-9]+.[0-9]+-qa\.[0-9]+$`, nil
	case ENV_UAT:
		return `^v[0-9]+\.[0-9]+.[0-9]+-uat\.[0-9]+$`, nil
	case ENV_STAGING:
		return `^v[0-9]+\.[0-9]+.[0-9]+-rc\.[0-9]+$`, nil
	case ENV_PROD:
		return `^v[0-9]+\.[0-9]+.[0-9]+$`, nil
	}
	return "", errors.New("an invalid or unsupported environment type was provided")
}

func (e Environment) DefaultVersion() (string, error) {
	switch e {
	case ENV_DEMO:
		return "v0.0.0-demo.0", nil	
	case ENV_DEV:
		return "v0.0.0-dev.0", nil
	case ENV_QA:
		return "v0.0.0-qa.0", nil
	case ENV_UAT:
		return "v0.0.0-uat.0", nil
	case ENV_STAGING:
		return "v0.0.0-rc.0", nil
	case ENV_PROD:
		return "v0.0.0", nil
	}
	return "", errors.New("an invalid or unsupported environment type was provided")
}

func (e Environment) IsValid() error {
	switch e {
	case ENV_DEMO, ENV_DEV, ENV_QA, ENV_UAT, ENV_STAGING, ENV_PROD:
		return nil
	}
	return errors.New("an invalid or unsupported environment type was provided")
}

func (e Environment) ValidateTag(t TagType) error {
	if e == ENV_DEMO && t == TAG_RC {
		return errors.New("release candidate tags cannot be used with demo environments")
	}
	if e == ENV_DEV && t == TAG_RC {
		return errors.New("release candidate tags cannot be used with development environments")
	}
	if e == ENV_QA && t == TAG_RC {
		return errors.New("release candidate tags cannot be used with QA environments")
	}
	if e == ENV_UAT && t == TAG_RC {
		return errors.New("release candidate tags cannot be used with QA environments")
	}
	if e == ENV_STAGING && t == TAG_DEV {
		return errors.New("development tags cannot be used with staging environments")
	}
	if e == ENV_PROD && (t == TAG_DEV || t == TAG_RC || t == TAG_QA || t == TAG_DEMO) {
		return errors.New("production tags cannot be used with staging or development environments")
	}
	return nil
}
