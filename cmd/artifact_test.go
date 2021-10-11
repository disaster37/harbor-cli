package cmd

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func (t CmdTestSuite) TestDeleteArtifact() {

	// Normale use case
	t.mockClient.EXPECT().Artifact().Return(t.mockArtifact)
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil)

	err := deleteArtifact("projectTest", "repositoryTest", "artifactTest", t.client)
	assert.NoError(t.T(), err)

	// When error
	t.mockClient.EXPECT().Artifact().Return(t.mockArtifact)
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(errors.New("fake error"))
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", t.client)
	assert.Error(t.T(), err)

}
