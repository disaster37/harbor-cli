package harborapi

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func (t *APITestSuite) TestScan() {

	// Normale use case
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/scan",
		httpmock.NewStringResponder(200, ""))
	err := t.client.Artifact().Scan("projectTest", "repositoryTest", "artifactTest")
	assert.NoError(t.T(), err)

	// Not found
	httpmock.Reset()
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/scan",
		httpmock.NewStringResponder(404, ""))
	err = t.client.Artifact().Scan("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// Not autorized
	httpmock.Reset()
	httpmock.RegisterResponder("POST", "http://localhost/projects/projectTest/repositories/repositoryTest/artifacts/artifactTest/scan",
		httpmock.NewStringResponder(403, ""))
	err = t.client.Artifact().Scan("projectTest", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)

	// error use cases
	err = t.client.Artifact().Scan("", "repositoryTest", "artifactTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().Scan("projectTest", "", "artifactTest")
	assert.Error(t.T(), err)
	err = t.client.Artifact().Scan("projectTest", "repositoryTest", "")
	assert.Error(t.T(), err)
}
