package cmd

import (
	"time"

	"github.com/disaster37/harbor-cli/harbor"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getClientWrapper(c *cli.Context) (*harbor.Client, error) {

	url := c.String("url")
	username := c.String("username")
	password := c.String("password")
	disableVerifySSL := c.Bool("self-signed-certificate")
	listCAPath := c.StringSlice("ca-path")
	timeout := c.Duration("timeout")

	log.Debugf("URL: %s", url)
	log.Debugf("Username: %s", username)
	log.Debugf("Self signed certificate: %t", disableVerifySSL)
	log.Debugf("List CA path: %+v", listCAPath)

	return getClient(url, username, password, disableVerifySSL, listCAPath, timeout)

}

func getClient(url string, username string, password string, disableVerifySSL bool, listCAPath []string, timeout time.Duration) (*harbor.Client, error) {

	if url == "" {
		return nil, errors.New("You need to set url")
	}

	config := harbor.Config{
		Address:          url,
		Username:         username,
		Password:         password,
		DisableVerifySSL: disableVerifySSL,
		CAs:              listCAPath,
		Timeout:          timeout,
	}

	return harbor.NewClient(config)
}
