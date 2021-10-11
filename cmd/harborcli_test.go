package cmd

import (
	"testing"
	"time"

	"github.com/disaster37/harbor-cli/harbor"
	"github.com/disaster37/harbor-cli/harbor/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdTestSuite struct {
	suite.Suite
	client       *harbor.Client
	mockClient   *mocks.MockAPI
	mockArtifact *mocks.MockArtifactAPI
	mockCtrl     *gomock.Controller
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

func (t *CmdTestSuite) BeforeTest(suiteName, testName string) {

	t.mockCtrl = gomock.NewController(t.T())
	t.mockClient = mocks.NewMockAPI(t.mockCtrl)
	t.mockArtifact = mocks.NewMockArtifactAPI(t.mockCtrl)

	t.client = &harbor.Client{
		API: t.mockClient,
	}

}

func (t *CmdTestSuite) AfterTest(suiteName, testName string) {
	defer t.mockCtrl.Finish()
}

func (t *CmdTestSuite) TestGetClient() {

	client, err := getClient("http://localhost", "user", "password", false, nil, 1*time.Second)
	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), client)
}
