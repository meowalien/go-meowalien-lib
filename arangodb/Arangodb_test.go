package arangodb

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/meowalien/go-meowalien-lib/arangodb/arangodb_wrapper"
	"github.com/meowalien/go-meowalien-lib/wrapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDatabaseConnection(t *testing.T) {
	cli, err := NewClient(context.TODO(), ArangoDBConnectionConfig{
		Address:         []string{"http://localhost:8529"},
		UserName:        "root",
		Password:        "",
		HTTPProtocol:    HTTP_1_1_PROTOCOL,
		ConnectionLimit: 40,
		//DefaultDoTimeout:   0,
		//ContentType:        0,
		//DontFollowRedirect: false,
		//FailOnRedirect:     false,
		//InsecureSkipVerify: false,
	})

	if !assert.NoError(t, err) {
		return
	}
	db, err := cli.Database(context.TODO(), "Database")
	if !assert.NoError(t, err) {
		return
	}
	db = wrapper.Wrap[driver.Database](db, &arangodb_wrapper.RetryWrapper{
		RetryCount:    5,
		RetryInterval: time.Millisecond * 300,
	})

	name := db.Name()
	fmt.Println(name)
	cs, err := db.Collections(context.TODO())
	if !assert.NoError(t, err) {
		return
	}
	for _, c := range cs {
		fmt.Println(c.Name())
	}
}
