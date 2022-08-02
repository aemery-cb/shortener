package main

import (
	"log"

	"github.com/aemery-cb/shortener/ui"
	"github.com/couchbase/gocb/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// const HOST_URL = "https://yock.to"

func main() {
	l, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
	}

	sugar := l.Sugar()

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

	store := NewStore(cluster)

	server := NewServer(viper.GetString("HOST"), store)

	app := fiber.New()

	app.Use(logger.New())
	app.Use("/", filesystem.New(filesystem.Config{
		Root: ui.BuildHTTPFS(),
	}))

	app.Use("/:code/stats", server.stats)
	app.Use("/:code", server.redirect)
	app.Post("/api/shorten", server.shorten)

	app.Use("/", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})
	log.Fatal(app.Listen(":8090"))
}
