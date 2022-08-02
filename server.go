package main

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store    *Store
	hostname string
}

func NewServer(hostname string, store *Store) Server {
	return Server{
		hostname: hostname,
		store:    store,
	}

}

type NewURLRequester struct {
	Url string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s *Server) GenerateURLKey() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	key := string(b)

	return key
}

func (s *Server) shorten(c *fiber.Ctx) error {

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

	if err := s.store.StoreUrl(key, urlReq.Url); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store URL")
	}

	final := fmt.Sprintf("%s/%s", s.hostname, key)

	response := map[string]string{"url": final}
	return c.JSON(response)

}

func (s *Server) redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	if code != "" {
		if url := s.store.GetUrl(code); url != "" {
			s.store.UpdateHitCounter(code)
			return c.Redirect(url)
		}
	}
	return c.Next()
}

func (s *Server) stats(c *fiber.Ctx) error {
	code := c.Params("code")
	return c.SendString(fmt.Sprintf("%d", s.store.GetStats(code)))
}
