package cmd

import (
	"errors"
	"time"

	harborapi "github.com/disaster37/harbor-cli/harbor/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func (t CmdTestSuite) TestCheckScanVulnerability() {

	t.mockClient.EXPECT().Artifact().AnyTimes().Return(t.mockArtifact)

	// Normale use case when no severity provided and no force scan
	artifact := &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport := harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err := checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "", 60*time.Second, false, t.client)
	assert.NoError(t.T(), err)

	// Normale use case when no severity provided and force scan
	t.mockArtifact.
		EXPECT().
		Scan(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil)

	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport = harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "", 60*time.Second, true, t.client)
	assert.NoError(t.T(), err)

	// Normale use case when severity High
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport = harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "high", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// Normale use case when severity Medium
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport = harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "medium", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// Normale use case when severity Low
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport = harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "low", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// When unknow severity
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	scanReport = harborapi.VulnerabilityReportResponse{
		"test": &harborapi.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Scanner: &harborapi.Scanner{
				Name:    "test",
				Version: "0.0.0",
				Vendor:  "test",
			},
			Severity: "High",
			Vulnerabilities: []*harborapi.Vulnerability{
				{
					ID:         "test",
					Package:    "test",
					Version:    "1.0.0",
					FixVersion: "1.0.0",
					Severity:   "High",
				},
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(scanReport, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "fake", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// When error on get artifact
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil, errors.New("Fake error"))

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "fake", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// When error on get scan result
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "Success",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	t.mockArtifact.
		EXPECT().
		GetVulnerabilities(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(nil, errors.New("fake errors"))

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "fake", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

	// When error on force scan
	t.mockArtifact.
		EXPECT().
		Scan(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(errors.New("Fake errors"))

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "fake", 60*time.Second, true, t.client)
	assert.Error(t.T(), err)

	// When bad scan status
	artifact = &harborapi.Artifact{
		ScanOverview: harborapi.ScanOverview{
			"test": &harborapi.NativeReportSummary{
				StartTime:  time.Now(),
				ScanStatus: "fake",
				EndTime:    time.Now(),
				Scanner: &harborapi.Scanner{
					Name:    "test",
					Version: "0.0.0",
					Vendor:  "test",
				},
				Summary: &harborapi.VulnerabilitySummary{
					Fixable: 0,
					Total:   0,
				},
				Severity: "High",
			},
		},
	}
	t.mockArtifact.
		EXPECT().
		Get(gomock.Eq("projectTest"), gomock.Eq("repositoryTest"), gomock.Eq("artifactTest")).
		Return(artifact, nil)

	err = checkScanVulnerability("projectTest", "repositoryTest", "artifactTest", "fake", 60*time.Second, false, t.client)
	assert.Error(t.T(), err)

}
