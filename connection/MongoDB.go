package connection

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MangoDBConfiguration struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	//Database string `json:"Database"`
}

func (dbconf MangoDBConfiguration) Create() (DB *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	option := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/"  , dbconf.User , dbconf.Password, dbconf.Host, dbconf.Port ))//.SetAuth(credential)
	DB, err = mongo.Connect(ctx, option)
	if err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = DB.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	//DB = DB.Database(dbconf.Database)
	return
}
