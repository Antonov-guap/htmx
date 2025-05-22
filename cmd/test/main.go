package main

import (
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// test commit
	app := fiber.New()

	app.Get("/event_stream", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/event-stream")
		c.Set(fiber.HeaderCacheControl, "no-cache")
		c.Set(fiber.HeaderConnection, fiber.HeaderKeepAlive)

		pr, pw := io.Pipe()
		defer func() { pw.Close() }()

		var bb []byte
		_, err := pr.Read(bb)
		fmt.Println(err)

		return c.SendStream(pr)
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
