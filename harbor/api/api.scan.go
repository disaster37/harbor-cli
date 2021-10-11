package harborapi

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	basePathScan = "/projects/%s/repositories/%s/artifacts/%s/scan"
)

type ScanOverview map[string]*NativeReportSummary
type VulnerabilityReportResponse map[string]*VulnerabilityReport

type NativeReportSummary struct {
	Description     string                `json:"description,omitempty"`
	StartTime       time.Time             `json:"start_time,omitempty"`
	ScanStatus      string                `json:"scan_status,omitempty"`
	CompletePercent int64                 `json:"complete_percent,omitempty"`
	EndTime         time.Time             `json:"end_time,omitempty"`
	ReportID        string                `json:"report_id,omitempty"`
	Severity        string                `json:"severity,omitempty"`
	Duration        int64                 `json:"duration,omitempty"`
	Scanner         *Scanner              `json:"scanner,omitempty"`
	Summary         *VulnerabilitySummary `json:"summary,omitempty"`
}

type Scanner struct {
	Version string `json:"version,omitempty"`
	Vendor  string `json:"vendor,omitempty"`
	Name    string `json:"name,omitempty"`
}

type VulnerabilitySummary struct {
	Description string           `json:"description,omitempty"`
	Fixable     int64            `json:"fixable,omitempty"`
	Total       int64            `json:"total,omitempty"`
	Summary     map[string]int64 `json:"summary,omitempty"`
}

type VulnerabilityReport struct {
	GeneratedAt     time.Time        `json:"generated_at,omitempty"`
	Scanner         *Scanner         `json:"scanner,omitempty"`
	Severity        string           `json:"severity,omitempty"`
	Vulnerabilities []*Vulnerability `json:"vulnerabilities,omitempty"`
}

type Vulnerability struct {
	ID          string   `json:"id,omitempty"`
	Package     string   `json:"package,omitempty"`
	Version     string   `json:"version,omitempty"`
	FixVersion  string   `json:"fix_version,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Links       []string `json:"links,omitempty"`
}

func (api *ArtifactAPIImpl) Scan(project, repositoryName, artifactName string) error {

	if project == "" {
		return errors.New("You must need provide project")
	}
	if repositoryName == "" {
		return errors.New("You must need provide repository name")
	}
	if artifactName == "" {
		return errors.New("You must need provide artifact name")
	}

	path := fmt.Sprintf(basePathScan, project, repositoryName, artifactName)

	resp, err := api.client.R().Post(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() >= 300 {
		return errors.Errorf("Error when run scan on %s/%s/%s: %s", project, repositoryName, artifactName, resp.String())
	}

	return nil
}
