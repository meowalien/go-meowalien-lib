package arangodb

import (
	"context"
	"crypto/tls"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/meowalien/go-meowalien-lib/errs"
	"golang.org/x/net/http2"
	"net"
	defaulthttp "net/http"
	"net/url"
	"time"
)

type _HTTP_PROTOCOL int

const (
	HTTP_1_1_PROTOCOL _HTTP_PROTOCOL = iota
	HTTP_2_0_PROTOCOL _HTTP_PROTOCOL = iota
)

type ArangoDBConnectionConfig struct {
	Address         []string
	UserName        string
	Password        string
	HTTPProtocol    _HTTP_PROTOCOL
	SSLCert         bool
	ConnectionLimit int
	Database        string
}

func NewDatabaseConnection(ctx context.Context, config ArangoDBConnectionConfig) (dbConn driver.Database, err error) {
	err = checkFormat(config)
	if err != nil {
		return
	}
	conn, err := createConnectionByHTTPProtocol(config)
	if err != nil {
		err = errs.New(err)
		return
	}

	c, err := driver.NewClient(driver.ClientConfig{
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
	dbConn, err = c.Database(ctx, config.Database)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func createConnectionByHTTPProtocol(config ArangoDBConnectionConfig) (conn driver.Connection, err error) {
	dbIPs, err := formatAsUrls(config.Address)
	if err != nil {
		err = errs.New(err)
		return
	}
	switch config.HTTPProtocol {
	case HTTP_1_1_PROTOCOL:
		transport := &defaulthttp.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:          0,
			IdleConnTimeout:       30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		conn, err = http.NewConnection(http.ConnectionConfig{
			Endpoints: dbIPs,
			ConnLimit: config.ConnectionLimit,
			Transport: transport,
		})
		if err != nil {
			err = errs.New(err)
			return
		}
	case HTTP_2_0_PROTOCOL:
		transport := &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
		conn, err = http.NewConnection(http.ConnectionConfig{
			Endpoints: dbIPs,
			ConnLimit: config.ConnectionLimit,
			Transport: transport,
		})
		if err != nil {
			err = errs.New(err)
			return
		}
	default:
		err = errs.New("Unsupported HTTP protocol")
		return
	}

	return
}

func formatAsUrls(host []string) ([]string, error) {
	for i, h := range host {
		var u *url.URL
		u, err := url.Parse("http://" + h)
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
	if config.Database == "" {
		return errs.New("Database is empty")
	}
	return nil
}
