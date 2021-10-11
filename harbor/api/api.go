package harborapi

type API interface {
	Artifact() ArtifactAPI
}

type ArtifactAPI interface {
	Scan(project string, repositoryName string, artifactName string) error
	Get(project string, repositoryName string, artifactName string) (*Artifact, error)
	GetVulnerabilities(project string, repositoryName string, artifactName string) (VulnerabilityReportResponse, error)
	Delete(project string, repositoryName string, artifactName string) error
}
