package cmd

import (
	"github.com/disaster37/harbor-cli/harbor"
	"github.com/urfave/cli/v2"
)

func DeleteArtifact(c *cli.Context) error {
	client, err := getClientWrapper(c)
	if err != nil {
		return err
	}

	if err := deleteArtifact(c.String("project"), c.String("repository"), c.String("artifact"), client); err != nil {
		return err
	}

	log.Infof("Successfully remove artifact %s/%s/%s", c.String("project"), c.String("repository"), c.String("artifact"))

}

func deleteArtifact(project, repository, artifact string, client *harbor.Client) error {

	return client.API.Artifact().Delete(project, repository, artifact)

}
