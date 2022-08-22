package cmd

import (
	"github.com/disaster37/harbor-cli/harbor"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func DeleteArtifact(c *cli.Context) error {
	client, err := getClientWrapper(c)
	if err != nil {
		return err
	}

	if err := deleteArtifact(c.String("project"), c.String("repository"), c.String("artifact"), c.String("tag"), client); err != nil {
		return err
	}

	if c.String("tag") == "" {
		log.Infof("Successfully remove artifact %s/%s/%s", c.String("project"), c.String("repository"), c.String("artifact"))
	} else {
		log.Infof("Successfully remove tag %s/%s/%s:%s", c.String("project"), c.String("repository"), c.String("artifact"), c.String("tag"))
	}

	return nil
}

func PromoteArtifact(c *cli.Context) error {
	client, err := getClientWrapper(c)
	if err != nil {
		return err
	}

	if err := promoteArtifact(c.String("project"), c.String("repository"), c.String("artifact"), c.String("source-tag"), c.StringSlice("target-tags"), client); err != nil {
		return err
	}

	for _, targetTag := range c.StringSlice("target-tags") {
		log.Infof("Successfully promote artifact %s/%s/%s:%s to %s/%s/%s:%s", c.String("project"), c.String("repository"), c.String("artifact"), c.String("source-tag"), c.String("project"), c.String("repository"), c.String("artifact"), targetTag)
	}
	return nil
}

func deleteArtifact(project, repository, artifact, tag string, client *harbor.Client) error {
	if tag == "" {
		return client.API.Artifact().Delete(project, repository, artifact)
	}

	tags, err := client.API.Artifact().GetTags(project, repository, artifact)
	if err != nil {
		return err
	}
	switch len(tags) {
	case 0:
		return client.API.Artifact().Delete(project, repository, artifact)
	case 1:
		if tags[0].Name == tag {
			return client.API.Artifact().Delete(project, repository, artifact)
		}
		return nil
	default:
		return client.API.Artifact().DeleteTag(project, repository, artifact, tag)
	}

}

func promoteArtifact(project, repository, artifactName, sourceTag string, targetTags []string, client *harbor.Client) error {

	for _, targetTag := range targetTags {

		// Check if tag is already use on other artifact
		// If tag already use and is the only tag, we delete the artifact
		// Else we only delete the tag
		artifact, err := client.API.Artifact().GetFromTag(project, repository, targetTag)
		if err != nil {
			return err
		}
		if artifact != nil {
			log.Debugf("Found artifact %d with tag %s", artifact.ID, targetTag)
			tags, err := client.API.Artifact().GetTags(project, repository, artifact.Digest)
			if err != nil {
				return err
			}
			if len(tags) == 1 {
				if err := client.API.Artifact().Delete(project, repository, artifact.Digest); err != nil {
					return err
				}
				log.Infof("We delete artifact %d that use target tag %s", artifact.ID, targetTag)
			} else {
				if err := client.API.Artifact().DeleteTag(project, repository, artifact.Digest, targetTag); err != nil {
					return err
				}
				log.Infof("We delete tag %s on artifact %d", targetTag, artifact.ID)
			}
		}

		// Add the target tag
		if err := client.API.Artifact().AddTag(project, repository, artifactName, targetTag); err != nil {
			return err
		}
	}

	// Delete the source tag
	if err := client.API.Artifact().DeleteTag(project, repository, artifactName, sourceTag); err != nil {
		return err
	}

	return nil
}
