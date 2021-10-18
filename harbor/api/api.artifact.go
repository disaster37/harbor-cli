package harborapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	basePathArtifact       = "/projects/%s/repositories/%s/artifacts/%s"
	basePathArtifactSearch = "/projects/%s/repositories/%s/artifacts"
)

type ArtifactAPIImpl struct {
	client *resty.Client
}

func NewArtifactAPI(client *resty.Client) ArtifactAPI {
	return &ArtifactAPIImpl{
		client: client,
	}
}

type Artifact struct {
	ID           int64        `json:"id"`
	Digest       string       `json:"digest"`
	Size         int64        `json:"size,omitempty"`
	PushTime     time.Time    `json:"push_time,omitempty"`
	PullTime     time.Time    `json:"pull_time,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	RepositoryID int64        `json:"repository_id,omitempty"`
	ProjectID    int64        `json:"project_id,omitempty"`
	Type         string       `json:"type,omitempty"`
	ScanOverview ScanOverview `json:"scan_overview,omitempty"`
}

type Tag struct {
	ID           int64     `json:"id,omitempty"`
	ArtifactID   int64     `json:"artifact_id,omitempty"`
	RepositoryID int64     `json:"repository_id,omitempty"`
	Name         string    `json:"name"`
	Immutable    bool      `json:"immutable,omitempty"`
	Signed       bool      `json:"signed,omitempty"`
	PullTime     time.Time `json:"pull_time,omitempty"`
	PushTime     time.Time `json:"push_time,omitempty"`
}

func (api *ArtifactAPIImpl) Get(project, repositoryName, artifactName string) (*Artifact, error) {
	if project == "" {
		return nil, errors.New("You must need provide project")
	}
	if repositoryName == "" {
		return nil, errors.New("You must need provide repository name")
	}
	if artifactName == "" {
		return nil, errors.New("You must need provide artifact name")
	}

	path := fmt.Sprintf(basePathArtifact, project, repositoryName, artifactName)

	resp, err := api.client.R().
		SetQueryParam("with_scan_overview", "true").
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, errors.Errorf("Error when get artifact %s/%s/%s: %s", project, repositoryName, artifactName, resp.Body())
	}

	artifact := new(Artifact)
	if err := json.Unmarshal(resp.Body(), artifact); err != nil {
		return nil, err
	}

	log.Debugf("Artifact: %+v", artifact)

	return artifact, nil
}

func (api *ArtifactAPIImpl) GetFromTag(project, repositoryName, tagName string) (*Artifact, error) {
	if project == "" {
		return nil, errors.New("You must need provide project")
	}
	if repositoryName == "" {
		return nil, errors.New("You must need provide repository name")
	}
	if tagName == "" {
		return nil, errors.New("You must need provide tag name")
	}

	path := fmt.Sprintf(basePathArtifactSearch, project, repositoryName)

	resp, err := api.client.R().
		SetQueryParam("q", fmt.Sprintf("tags=%s", tagName)).
		Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, errors.Errorf("Error when get artifact %s/%s:%s: %s", project, repositoryName, tagName, resp.Body())
	}

	artifacts := make([]*Artifact, 0)

	if err := json.Unmarshal(resp.Body(), &artifacts); err != nil {
		return nil, err
	}

	if len(artifacts) == 0 {
		return nil, nil
	}

	log.Debugf("Artifact: %+v", artifacts[0])

	return artifacts[0], nil
}

func (api *ArtifactAPIImpl) GetVulnerabilities(project, repositoryName, artifactName string) (VulnerabilityReportResponse, error) {
	if project == "" {
		return nil, errors.New("You must need provide project")
	}
	if repositoryName == "" {
		return nil, errors.New("You must need provide repository name")
	}
	if artifactName == "" {
		return nil, errors.New("You must need provide artifact name")
	}

	path := fmt.Sprintf(basePathArtifact+"/additions/vulnerabilities", project, repositoryName, artifactName)

	resp, err := api.client.R().Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, errors.Errorf("Error when get vulnerabilities for %s/%s/%s: %s", project, repositoryName, artifactName, resp.Body())
	}

	report := new(VulnerabilityReportResponse)
	if err := json.Unmarshal(resp.Body(), report); err != nil {
		return nil, err
	}

	log.Debugf("Vulnerability report: %+v", report)

	return *report, nil
}

func (api *ArtifactAPIImpl) Delete(project, repositoryName, artifactName string) error {
	if project == "" {
		return errors.New("You must need provide project")
	}
	if repositoryName == "" {
		return errors.New("You must need provide repository name")
	}
	if artifactName == "" {
		return errors.New("You must need provide artifact name")
	}

	path := fmt.Sprintf(basePathArtifact, project, repositoryName, artifactName)

	resp, err := api.client.R().
		Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil
		}
		return errors.Errorf("Error when delete artifact %s/%s/%s: %s", project, repositoryName, artifactName, resp.Body())
	}

	log.Debugf("Artifact successfully deleted %s/%s/%s", project, repositoryName, artifactName)

	return nil
}

func (api *ArtifactAPIImpl) AddTag(project, repository, artifact, tagName string) error {
	if project == "" {
		return errors.New("You must need provide project")
	}
	if repository == "" {
		return errors.New("You must need provide repository")
	}
	if artifact == "" {
		return errors.New("You must need provide artifact")
	}
	if tagName == "" {
		return errors.New("You must need provide tag")
	}

	tag := &Tag{
		Name: tagName,
	}

	path := fmt.Sprintf(basePathArtifact+"/tags", project, repository, artifact)

	resp, err := api.client.R().
		SetBody(tag).
		Post(path)

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 300 {
		// Ignore if tag already exist
		if resp.StatusCode() != 409 {
			return errors.Errorf("Error when add tag %s on artifact %s/%s/%s: %s", tagName, project, repository, artifact, resp.Body())
		}
	}

	log.Debugf("Tag %s successfully added to %s/%s/%s", tagName, project, repository, artifact)

	return nil
}

func (api *ArtifactAPIImpl) DeleteTag(project, repository, artifact, tagName string) error {
	if project == "" {
		return errors.New("You must need provide project")
	}
	if repository == "" {
		return errors.New("You must need provide repository")
	}
	if artifact == "" {
		return errors.New("You must need provide artifact")
	}
	if tagName == "" {
		return errors.New("You must need provide tag")
	}

	path := fmt.Sprintf(basePathArtifact+"/tags/%s", project, repository, artifact, tagName)

	resp, err := api.client.R().
		Delete(path)

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 300 {
		// Ignore if tag not exist
		if resp.StatusCode() != 404 {
			return errors.Errorf("Error when delete tag %s on artifact %s/%s/%s: %s", tagName, project, repository, artifact, resp.Body())
		}
	}

	log.Debugf("Tag %s successfully deleted to %s/%s/%s", tagName, project, repository, artifact)

	return nil
}

func (api *ArtifactAPIImpl) GetTags(project, repository, artifact string) (listTags []Tag, err error) {
	if project == "" {
		return nil, errors.New("You must need provide project")
	}
	if repository == "" {
		return nil, errors.New("You must need provide repository")
	}
	if artifact == "" {
		return nil, errors.New("You must need provide artifact")
	}

	path := fmt.Sprintf(basePathArtifact+"/tags", project, repository, artifact)

	resp, err := api.client.R().
		Get(path)

	if err != nil {
		return nil, err
	}

	listTags = make([]Tag, 0)

	// Return error only if not 404
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return listTags, nil
		}
		return nil, errors.Errorf("Error when list tags on artifact %s/%s/%s: %s", project, repository, artifact, resp.Body())
	}

	if err := json.Unmarshal(resp.Body(), &listTags); err != nil {
		return nil, err
	}

	log.Debugf("Successfully get list tags to %s/%s/%s", project, repository, artifact)

	return listTags, nil
}
