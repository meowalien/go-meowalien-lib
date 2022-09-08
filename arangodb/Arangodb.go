package arangodb

import (
	"context"
	"crypto/tls"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/cluster"
	"github.com/arangodb/go-driver/http"
	"github.com/meowalien/go-meowalien-lib/errs"
	"golang.org/x/net/http2"
	"net"
	defaulthttp "net/http"
	"net/url"
	"strings"
	"time"
)

type HTTP_PROTOCOL string

const (
	HTTP_1_1_PROTOCOL HTTP_PROTOCOL = "1.1"
	HTTP_2_0_PROTOCOL HTTP_PROTOCOL = "2.0"
)

type ArangoDBConnectionConfig struct {
	// the scheme could be http:// or https://, if scheme not set, default as http://
	Address          []string
	UserName         string
	Password         string
	HTTPProtocol     HTTP_PROTOCOL
	ConnectionLimit  int
	DefaultDoTimeout time.Duration
	// json 0, velocypack 1
	ContentType        driver.ContentType
	DontFollowRedirect bool
	FailOnRedirect     bool
	InsecureSkipVerify bool
}

func NewClient(ctx context.Context, config ArangoDBConnectionConfig) (c driver.Client, err error) {
	err = checkFormat(config)
	if err != nil {
		return
	}
	conn, err := createConnection(config)
	if err != nil {
		err = errs.New(err)
		return
	}

	c, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.UserName, config.Password),
	})
	if err != nil {
		err = errs.New(err)
		return
	}
	if _, err1 := c.Version(ctx); err1 != nil {
		err = errs.New(err1)
		return
	}
	return
}

func createConnection(config ArangoDBConnectionConfig) (conn driver.Connection, err error) {
	var tlsConfig *tls.Config = nil
	dbIPs, err := formatAsUrls(config.Address)
	if err != nil {
		err = errs.New(err)
		return
	}

	for _, p := range dbIPs {
		if strings.HasPrefix(p, "https://") {
			tlsConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}
		}
	}

	var transport defaulthttp.RoundTripper
	switch config.HTTPProtocol {
	case HTTP_1_1_PROTOCOL:
	case HTTP_2_0_PROTOCOL:
		transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
	default:
		err = errs.New("Unsupported HTTP protocol")
		return
	}
	conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints:          dbIPs,
		TLSConfig:          tlsConfig,
		Transport:          transport,
		DontFollowRedirect: config.DontFollowRedirect,
		FailOnRedirect:     config.FailOnRedirect,
		ConnectionConfig:   cluster.ConnectionConfig{DefaultTimeout: config.DefaultDoTimeout},
		ContentType:        config.ContentType,
		ConnLimit:          config.ConnectionLimit,
	})
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func formatAsUrls(host []string) ([]string, error) {
	for i, h := range host {
		var u *url.URL
		var err error
		if strings.HasPrefix(h, "http://") {
			u, err = url.Parse(h)
		} else if strings.HasPrefix(h, "https://") {
			u, err = url.Parse(h)
		} else {
			u, err = url.Parse("http://" + h)
		}
		if err != nil {
			err = errs.New(err)
			return nil, err
		}
		host[i] = u.String()
	}
	return host, nil
}

func checkFormat(config ArangoDBConnectionConfig) error {
	if config.Address == nil {
		return errs.New("Address is empty")
	}
	switch config.HTTPProtocol {
	case HTTP_1_1_PROTOCOL:
	case HTTP_2_0_PROTOCOL:
	default:
		return errs.New("Unsupported HTTP protocol")
	}

	return nil
}
