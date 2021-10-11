package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/disaster37/harbor-cli/harbor"
	harborapi "github.com/disaster37/harbor-cli/harbor/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func CheckScanVulnerability(c *cli.Context) error {
	client, err := getClientWrapper(c)
	if err != nil {
		return err
	}

	return checkScanVulnerability(c.String("project"), c.String("repository"), c.String("artifact"), c.String("severity"), c.Duration("timeout"), c.Bool("force-scan"), client)
}

func checkScanVulnerability(project string, repositoryName string, artifactName string, severity string, timeout time.Duration, forceScan bool, client *harbor.Client) error {

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if forceScan {
		if err := client.API.Artifact().Scan(project, repositoryName, artifactName); err != nil {
			return err
		}
		log.Infof("Scan running")
	}

	// Wait the end of scan
	isScan := true
	var report *harborapi.NativeReportSummary
	for isScan {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			artifact, err := client.API.Artifact().Get(project, repositoryName, artifactName)
			if err != nil {
				return err
			}
			if artifact == nil {
				return errors.Errorf("Can't find %s/%s/%s", project, repositoryName, artifactName)
			}

			for _, report = range artifact.ScanOverview {
				if !report.EndTime.IsZero() {
					// Scan finished
					isScan = false
					break
				}
				time.Sleep(1 * time.Second)
			}
		}

	}

	// Check result
	if report.ScanStatus != "Success" {
		return errors.Errorf("Error during scan, status: %s", report.ScanStatus)
	}

	log.Infof("Scan successfully finished with %s / %s / %s", report.Scanner.Name, report.Scanner.Vendor, report.Scanner.Version)

	// Display summary
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Scan date", report.EndTime.String()})
	table.Append([]string{"Severity", report.Severity})
	table.Append([]string{"Totale", fmt.Sprintf("%d", report.Summary.Total)})
	table.Append([]string{"Fixable", fmt.Sprintf("%d", report.Summary.Fixable)})
	table.Append([]string{"Severity High", fmt.Sprintf("%d", report.Summary.Summary["High"])})
	table.Append([]string{"Severity Medium", fmt.Sprintf("%d", report.Summary.Summary["Medium"])})
	table.Append([]string{"Severity Low", fmt.Sprintf("%d", report.Summary.Summary["Low"])})
	table.Render()

	// Get and display vulnerabilities
	vulnReport, err := client.API.Artifact().GetVulnerabilities(project, repositoryName, artifactName)
	if err != nil {
		return err
	}
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"LIBRARY", "VULNERABILITY ID", "SEVERITY", "INSTALLED VERSION", "FIXED VERSION", "Links"})
	for _, report := range vulnReport {
		for _, vuln := range report.Vulnerabilities {
			table.Append([]string{vuln.Package, vuln.ID, vuln.Severity, vuln.Version, vuln.FixVersion, strings.Join(vuln.Links, "\n")})
		}
	}

	table.Render()

	// Compute if error base on current Severity
	switch strings.ToLower(severity) {
	case "high":
		if report.Severity == "High" {
			return errors.Errorf("Current severity is %s", report.Severity)
		}
		break
	case "medium":
		if report.Severity == "High" || report.Severity == "Medium" {
			return errors.Errorf("Current severity is %s", report.Severity)
		}
		break
	case "low":
		if report.Severity == "High" || report.Severity == "Medium" || report.Severity == "Low" {
			return errors.Errorf("Current severity is %s", report.Severity)
		}
		break
	case "":
		break
	default:
		return errors.Errorf("Severity must be high, medium or low")
	}

	return nil

}
