package harborapi

type API interface {
	Artifact() ArtifactAPI
}

type ArtifactAPI interface {
	Scan(project, repositoryName, artifactName string) error
	Get(project, repositoryName, artifactName string) (*Artifact, error)
	GetVulnerabilities(project, repositoryName, artifactName string) (VulnerabilityReportResponse, error)
	Delete(project, repositoryName, artifactName string) error
	AddTag(project, repository, artifact, tag string) error
	DeleteTag(project, repository, artifact, tag string) error
	GetTags(project, repository, artifact string) (listTags []Tag, err error)
}
