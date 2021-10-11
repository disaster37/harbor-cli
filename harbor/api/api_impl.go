package harborapi

import "github.com/go-resty/resty/v2"

type APIImpl struct {
	artifact ArtifactAPI
}

func New(client *resty.Client) API {
	return &APIImpl{
		artifact: NewArtifactAPI(client),
	}
}

func (api *APIImpl) Artifact() ArtifactAPI {
	return api.artifact
}
