package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"htmx/cmd/tictactoe/internal/broadcast"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/samber/lo"
)

//go:embed internal/views
var viewsFS embed.FS

//go:embed internal/assets/*
var assetsFS embed.FS

func main() {
	// infra specific
	views := html.NewFileSystem(http.FS(viewsFS), ".html")
	views.Reload(true)
	views.Directory = "internal/views"
	app := fiber.New(fiber.Config{
		Views:             views,
		PassLocalsToViews: true,
	})
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(assetsFS),
		PathPrefix: "assets",
	}))

	// business specific
	currentPlayer := "X"

	var gameMu sync.RWMutex

	gameField := [][]string{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}

	// infra common
	type event struct {
		name string
		data string
	}

	var bus broadcast.Broadcast[event]
	bus.Logger = log.Default()

	app.Get("/events", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderConnection, fiber.HeaderKeepAlive)
		c.Set(fiber.HeaderContentType, "text/event-stream")
		c.Set(fiber.HeaderCacheControl, "no-cache")

		events := bus.Subscribe()

		r, w := io.Pipe()
		go func() {
			defer func() { _ = w.Close() }()

			_, _ = fmt.Fprintln(w, "event: init-stream")
			_, _ = fmt.Fprintln(w, "data:")
			_, _ = fmt.Fprintln(w)

			for e := range events {
				_, _ = fmt.Fprintln(w, "event:", e.name)
				_, _ = fmt.Fprintln(w, "data:", e.data)
				_, _ = fmt.Fprintln(w)
			}
		}()

		err := c.SendStream(r)
		return err
	})

	// infra-business specific
	app.Get("/", func(c *fiber.Ctx) error {
		gameMu.RLock()
		defer gameMu.RUnlock()

		c.Locals("CurrentPlayer", currentPlayer)
		c.Locals("GameField", gameField)

		return c.Render("index", fiber.Map{}, "default-layout")
	})

	// infra-business specific
	app.Post("/make-turn", func(c *fiber.Ctx) error {
		gameMu.Lock()
		defer gameMu.Unlock()

		position := c.Request().PostArgs().PeekMulti("position")
		row := lo.Must(strconv.Atoi(string(position[0])))
		col := lo.Must(strconv.Atoi(string(position[1])))
		if gameField[row][col] == " " {
			gameField[row][col] = currentPlayer
			if currentPlayer == "X" {
				currentPlayer = "O"
			} else {
				currentPlayer = "X"
			}
			bus.Send(event{fmt.Sprintf("cell-updated-%d-%d", row, col), gameField[row][col]})
			bus.Send(event{"player-updated", currentPlayer})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// infra-business specific
	app.Post("/new-game", func(c *fiber.Ctx) error {
		gameMu.Lock()
		defer gameMu.Unlock()

		currentPlayer = "X"
		lo.ForEach(gameField, func(cells []string, i int) {
			lo.ForEach(cells, func(_ string, j int) {
				cells[j] = " "
				bus.Send(event{fmt.Sprintf("cell-updated-%d-%d", i, j), " "})
			})
		})

		bus.Send(event{"player-updated", currentPlayer})
		return c.SendStatus(fiber.StatusNoContent)
	})

	// infra specific
	lo.Must0(app.ListenTLS(
		"0.0.0.0:8080",
		"cert/conf/live/htmx.space/fullchain.pem",
		"cert/conf/live/htmx.space/privkey.pem",
	))
	// lo.Must0(app.Listen(":8080"))
}
