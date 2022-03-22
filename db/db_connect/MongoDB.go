package db_connect

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/db/config_modules"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func CreateMongoDBConnection(dbconf config_modules.MangoDBConfiguration) (client *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: dbconf.User,
		Password: dbconf.Password,
	}

	option := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", dbconf.Host, dbconf.Port)).SetAuth(credential)
	client, err = mongo.Connect(ctx, option)
	if err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	return
}
