package server

import (
	"fmt"
	"net/url"

	"github.com/aemery-cb/shortener/pkg/store"
	"github.com/aemery-cb/shortener/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	jwtware "github.com/gofiber/jwt/v3"
)

type Server struct {
	store    *store.Store
	hostname string
	logger   *zap.SugaredLogger
}

func NewServer(hostname string, store *store.Store, logger *zap.SugaredLogger) Server {
	return Server{
		hostname: hostname,
		store:    store,
		logger:   logger,
	}
}

type NewURLRequester struct {
	Url string
}

func (s *Server) Run(port string) error {
	app := fiber.New()

	app.Use(logger.New())
	app.Use("/", filesystem.New(filesystem.Config{
		Root: ui.BuildHTTPFS(),
	}))

	app.Use("/:code/stats", s.GetURLStats)
	app.Use("/:code", s.GetURL)

	successHandler := func(c *fiber.Ctx) error {
		token := c.Locals("user")
		if parsed, ok := token.(*jwt.Token); ok {
			if mapClaims, moreOk := parsed.Claims.(jwt.MapClaims); moreOk {
				c.Locals("user_id", mapClaims["sub"].(string))
			}
		}
		return c.Next()
	}

	app.Use("/api/shorten", func(c *fiber.Ctx) error {
		if _, ok := c.GetReqHeaders()["Authorization"]; !ok {
			return s.ShortenURL(c)
		}
		return c.Next()
	})

	app.Use(jwtware.New(jwtware.Config{
		KeySetURL:      fmt.Sprintf("https://%s/.well-known/jwks.json", viper.GetString("AUTH0_DOMAIN")),
		SuccessHandler: successHandler,
	}))
	app.Use("/api/shorten", s.ShortenURL)

	app.Use("/", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	return app.Listen(port)
}

func (s *Server) ShortenURL(c *fiber.Ctx) error {
	urlReq := &NewURLRequester{}

	if err := c.BodyParser(urlReq); err != nil {
		return err
	}

	if urlReq.Url == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("No URL provided")
	}

	if len(urlReq.Url) > 64000 {
		return c.Status(fiber.StatusBadRequest).SendString("URL too long")
	}

	if _, err := url.ParseRequestURI(urlReq.Url); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid URL")
	}

	key := s.GenerateURLKey()
	var userId string = ""

	if maybeUserId, ok := c.Locals("user_id").(string); ok {
		userId = maybeUserId
	}

	if err := s.store.StoreUrl(key, urlReq.Url, userId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store URL")
	}

	final := fmt.Sprintf("%s/%s", s.hostname, key)

	response := map[string]string{"url": final}
	return c.JSON(response)
}

func (s *Server) GetURL(c *fiber.Ctx) error {
	code := c.Params("code")
	if code != "" {
		if url := s.store.GetUrl(code); url != "" {
			s.store.UpdateHitCounter(code)
			return c.Redirect(url)
		}
	}
	return c.Next()
}

func (s *Server) GetURLStats(c *fiber.Ctx) error {
	code := c.Params("code")
	return c.SendString(fmt.Sprintf("%d", s.store.GetStats(code)))
}
