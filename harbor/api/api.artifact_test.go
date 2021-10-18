package harborapi

import (
	"encoding/json"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func (t *APITestSuite) TestGet() {

	// Normal use case
	response := &Artifact{
		ID: 10,
	}
	responder, err := httpmock.NewJsonResponder(200, response)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest?with_scan_overview=true", responder)
	artifact, err := t.client.Artifact().Get("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), response.ID, artifact.ID)

	// Not found
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest?with_scan_overview=true",
		httpmock.NewStringResponder(404, ""))
	artifact, err = t.client.Artifact().Get("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Nil(t.T(), artifact)

	// Not authorized
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest?with_scan_overview=true",
		httpmock.NewStringResponder(403, ""))
	artifact, err = t.client.Artifact().Get("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// error use cases
	_, err = t.client.Artifact().Get("", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().Get("projectTest", "", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().Get("projectTest", "repositoryTest", "")
	assert.Error(t.T(), err)

}

func (t *APITestSuite) TestGetVulnerabilities() {
	// Normal use case
	response := VulnerabilityReportResponse{
		"test": &VulnerabilityReport{
			Severity: "High",
		},
	}
	responder, err := httpmock.NewJsonResponder(200, response)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/additions/vulnerabilities", responder)
	reportResponse, err := t.client.Artifact().GetVulnerabilities("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), response["test"].Severity, reportResponse["test"].Severity)

	// Not found
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/additions/vulnerabilities",
		httpmock.NewStringResponder(404, ""))
	reportResponse, err = t.client.Artifact().GetVulnerabilities("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Nil(t.T(), reportResponse)

	// Not authorized
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/additions/vulnerabilities",
		httpmock.NewStringResponder(403, ""))
	_, err = t.client.Artifact().GetVulnerabilities("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// error use cases
	_, err = t.client.Artifact().GetVulnerabilities("", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().GetVulnerabilities("projectTest", "", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().GetVulnerabilities("projectTest", "repositoryTest", "")
	assert.Error(t.T(), err)
}

func (t *APITestSuite) TestDelete() {
	// Normal use case
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest",
		httpmock.NewStringResponder(200, ""))
	err := t.client.Artifact().Delete("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)

	// Not found
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest",
		httpmock.NewStringResponder(404, ""))
	err = t.client.Artifact().Delete("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)

	// Not authorized
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest",
		httpmock.NewStringResponder(403, ""))
	err = t.client.Artifact().Delete("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// error use cases
	err = t.client.Artifact().Delete("", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().Delete("projectTest", "", "artifactTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().Delete("projectTest", "repositoryTest", "")
	assert.Error(t.T(), err)
}

func (t *APITestSuite) TestAddTag() {

	// Normal use case
	tag := &Tag{
		Name: "tagTest",
	}
	tagTmp := new(Tag)
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags",
		func(r *http.Request) (*http.Response, error) {

			if err := json.NewDecoder(r.Body).Decode(tagTmp); err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return httpmock.NewStringResponse(200, ""), nil
		})
	err := t.client.Artifact().AddTag("projectTest", "repositoryTest", "artifactTest", tag.Name)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), tag.Name, tagTmp.Name)

	// Artifact not exist
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags",
		httpmock.NewStringResponder(404, ""))
	err = t.client.Artifact().AddTag("projectTest", "repositoryTest", "artifactTest", "tagTest")
	assert.Error(t.T(), err)

	// Not authorized
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags",
		httpmock.NewStringResponder(403, ""))
	err = t.client.Artifact().AddTag("projectTest", "repositoryTest", "artifactTest", "tagTest")
	assert.Error(t.T(), err)

	// error use cases
	err = t.client.Artifact().AddTag("", "repositoryTest", "artifactTest", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().AddTag("projectTest", "", "artifactTest", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().AddTag("projectTest", "repositoryTest", "", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().AddTag("projectTest", "repositoryTest", "artifactTest", "")
	assert.Error(t.T(), err)
}

func (t *APITestSuite) TestDeleteTag() {
	// Normal use case
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags/tagTest",
		httpmock.NewStringResponder(200, ""))
	err := t.client.Artifact().DeleteTag("projectTest", "repositoryTest", "artifactTest", "tagTest")
	assert.NoError(t.T(), err)

	// Not found
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags/tagTest",
		httpmock.NewStringResponder(404, ""))
	err = t.client.Artifact().DeleteTag("projectTest", "repositoryTest", "artifactTest", "tagTest")
	assert.NoError(t.T(), err)

	// Not authorized
	httpmock.RegisterResponder("DELETE", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags/tagTest",
		httpmock.NewStringResponder(403, ""))
	err = t.client.Artifact().DeleteTag("projectTest", "repositoryTest", "artifactTest", "tagTest")
	assert.Error(t.T(), err)

	// error use cases
	err = t.client.Artifact().DeleteTag("", "repositoryTest", "artifactTest", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().DeleteTag("projectTest", "", "artifactTest", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().DeleteTag("projectTest", "repositoryTest", "", "tagTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().DeleteTag("projectTest", "repositoryTest", "artifactTest", "")
	assert.Error(t.T(), err)
}

func (t *APITestSuite) TestGetTags() {

	// Normal use case
	response := []Tag{
		{
			Name: "tagTest",
		},
	}
	responder, err := httpmock.NewJsonResponder(200, response)
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags", responder)
	tags, err := t.client.Artifact().GetTags("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), response[0].Name, tags[0].Name)

	// Not found
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags",
		httpmock.NewStringResponder(404, ""))
	tags, err = t.client.Artifact().GetTags("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)
	assert.Empty(t.T(), tags)

	// Not authorized
	httpmock.RegisterResponder("GET", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/tags",
		httpmock.NewStringResponder(403, ""))
	tags, err = t.client.Artifact().GetTags("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// error use cases
	_, err = t.client.Artifact().GetTags("", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().GetTags("projectTest", "", "artifactTest")
	assert.Error(t.T(), err)
	_, err = t.client.Artifact().GetTags("projectTest", "repositoryTest", "")
	assert.Error(t.T(), err)

}
