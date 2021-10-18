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
	debug := c.Bool("debug")

	log.Debugf("URL: %s", url)
	log.Debugf("Username: %s", username)
	log.Debugf("Self signed certificate: %t", disableVerifySSL)
	log.Debugf("List CA path: %+v", listCAPath)
	log.Debugf("Debug: %t", debug)

	return getClient(url, username, password, disableVerifySSL, listCAPath, timeout, debug)

}

func getClient(url string, username string, password string, disableVerifySSL bool, listCAPath []string, timeout time.Duration, debug bool) (*harbor.Client, error) {

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
		Debug:            debug,
	}

	return harbor.NewClient(config)
}
