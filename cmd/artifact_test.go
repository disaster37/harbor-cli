package cmd

import (
	"github.com/disaster37/harbor-cli/harbor/api"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func (t CmdTestSuite) TestDeleteArtifact() {

	t.mockClient.EXPECT().Artifact().AnyTimes().Return(t.mockArtifact)

	// Normale use case when no tag provided

	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil)

	err := deleteArtifact("projectTest", "repositoryTest", "artifactTest", "", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when tag provided and no tag on current artifact
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return([]harborapi.Tag{}, nil)
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil)
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when tag provided and one tag on current artifact that match the tag
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return([]harborapi.Tag{
			{
				Name: "tagTest",
			},
		}, nil)
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil)
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when tag provided and one tag on current artifact that not match the tag
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return([]harborapi.Tag{
			{
				Name: "fake",
			},
		}, nil)
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when tag provided and multi tag on current artifact
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return([]harborapi.Tag{
			{
				Name: "fake",
			},
			{
				Name: "tagTest",
			},
		}, nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("tagTest")).
		Return(nil)
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest", t.client)
	assert.NoError(t.T(), err)

	// When error
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(errors.New("fake error"))
	err = deleteArtifact("projectTest", "repositoryTest", "artifactTest", "", t.client)
	assert.Error(t.T(), err)

}

func (t CmdTestSuite) TestPromoteArtifact() {

	t.mockClient.EXPECT().Artifact().AnyTimes().Return(t.mockArtifact)

	// Normale use case
	t.mockArtifact.
		EXPECT().
		GetFromTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("testTag2")).
		Return(nil, nil)
	t.mockArtifact.
		EXPECT().
		AddTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag2")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag1")).
		Return(nil)
	err := promoteArtifact("projectTest", "repositoryTest", "artifactTest", "testTag1", "testTag2", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when target tag already in use and it alone
	t.mockArtifact.
		EXPECT().
		GetFromTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("testTag2")).
		Return(&harborapi.Artifact{
			Digest: "test",
		}, nil)
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("test")).
		Return([]harborapi.Tag{
			{
				Name: "testTag2",
			},
		}, nil)
	t.mockArtifact.
		EXPECT().
		Delete(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("test")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		AddTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag2")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag1")).
		Return(nil)
	err = promoteArtifact("projectTest", "repositoryTest", "artifactTest", "testTag1", "testTag2", t.client)
	assert.NoError(t.T(), err)

	// Normale use case when target tag already in use and it not alone
	t.mockArtifact.
		EXPECT().
		GetFromTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("testTag2")).
		Return(&harborapi.Artifact{
			Digest: "test",
		}, nil)
	t.mockArtifact.
		EXPECT().
		GetTags(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("test")).
		Return([]harborapi.Tag{
			{
				Name: "testTag2",
			},
			{
				Name: "fake",
			},
		}, nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("test"), gomock.Eq("testTag2")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		AddTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag2")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("testTag1")).
		Return(nil)
	err = promoteArtifact("projectTest", "repositoryTest", "artifactTest", "testTag1", "testTag2", t.client)
	assert.NoError(t.T(), err)

	// When error on delete tag
	t.mockArtifact.
		EXPECT().
		GetFromTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("tagTest2")).
		Return(nil, nil)
	t.mockArtifact.
		EXPECT().
		AddTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("tagTest2")).
		Return(nil)
	t.mockArtifact.
		EXPECT().
		DeleteTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("tagTest1")).
		Return(errors.New("fake error"))
	err = promoteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest1", "tagTest2", t.client)
	assert.Error(t.T(), err)

	// When error on add tag
	t.mockArtifact.
		EXPECT().
		GetFromTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("tagTest2")).
		Return(nil, nil)
	t.mockArtifact.
		EXPECT().
		AddTag(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest"), gomock.Eq("tagTest2")).
		Return(errors.New("fake error"))
	err = promoteArtifact("projectTest", "repositoryTest", "artifactTest", "tagTest1", "tagTest2", t.client)
	assert.Error(t.T(), err)

}
