package db_connect

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MangoDBConfiguration struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Database string `json:"Database"`
}

func CreateMongoDBConnection(dbconf MangoDBConfiguration) (DB *mongo.Database, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//credential := options.Credential{
	//	Username: dbconf.User,
	//	Password: dbconf.Password,
	//}

	option := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s"  , dbconf.User , dbconf.Password, dbconf.Host, dbconf.Port , dbconf.Database))//.SetAuth(credential)
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}
	DB = client.Database(dbconf.Database)
	return
}
