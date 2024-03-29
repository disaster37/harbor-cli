package harbor

import (
	"crypto/tls"
	"time"

	"github.com/disaster37/harbor-cli/harbor/api"
	"github.com/go-resty/resty/v2"
)

// Config contain the value to access on Kibana API
type Config struct {
	Address          string
	Username         string
	Password         string
	DisableVerifySSL bool
	CAs              []string
	Timeout          time.Duration
	Debug            bool
}

// Client contain the REST client and the API specification
type Client struct {
	API harborapi.API
}

// NewDefaultClient init client with empty config
func NewDefaultClient() (*Client, error) {
	return NewClient(Config{})
}

// NewClient init client with custom config
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = "http://localhost/api/v2.0"
	}

	restyClient := resty.New().
		SetBaseURL(cfg.Address).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetHeader("Content-Type", "application/json").
		SetTimeout(cfg.Timeout).
		SetDebug(cfg.Debug).
		SetCookieJar(nil)

	for _, path := range cfg.CAs {
		restyClient.SetRootCertificate(path)
	}

	if cfg.DisableVerifySSL {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client := &Client{
		API: harborapi.New(restyClient),
	}

	return client, nil

}
