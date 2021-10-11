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

	return deleteArtifact(c.String("project"), c.String("repository"), c.String("artifact"), client)

}

func deleteArtifact(project, repository, artifact string, client *harbor.Client) error {

	return client.API.Artifact().Delete(project, repository, artifact)

}
