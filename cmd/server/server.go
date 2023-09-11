package main

import (
	"html/template"
	"log"

	"htmx/cmd/server/internal/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/samber/lo"
)

func main() {
	engine := html.New("./cmd/server/internal", ".gohtml")
	engine.Reload(true)
	engine.AddFunc("componentScope", func() any {
		return template.HTMLAttr("data-gotmpl-" + lo.RandomString(6, lo.LowerCaseLettersCharset))
	})

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())
	app.All("/:block?", pages.GegIndex)

	log.Fatal(app.Listen(":3000"))
}
