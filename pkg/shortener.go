package shortener

import (
	"log"
	"os"
	"strings"

	"github.com/aemery-cb/shortener/pkg/server"
	"github.com/aemery-cb/shortener/pkg/store"
	"github.com/couchbase/gocb/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run() {
	var zlog *zap.Logger
	var err error

	if strings.ToLower(os.Getenv("BUILD")) == "true" {
		zlog, err = zap.NewDevelopment()
	} else {
		zlog, err = zap.NewProduction()

	}

	if err != nil {
		log.Fatal(err)
	}

	sugar := zlog.Sugar()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		sugar.Fatal(err)
	}

	endpoint := viper.GetString("DB_URL")
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASS")

	cluster, err := gocb.Connect(endpoint, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})

	if err != nil {
		sugar.Fatal(err)
	}

	store := store.New(cluster)

	srv := server.NewServer(viper.GetString("HOST"), store, sugar)

	log.Fatal(srv.Run(":8090"))
}
