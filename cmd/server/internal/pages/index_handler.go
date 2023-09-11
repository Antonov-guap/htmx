package pages

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

var counter int

func GegIndex(c *fiber.Ctx) error {
	block := c.Params("block")

	if block == "counter" {
		if c.Method() == fiber.MethodPost {
			counter++
		} else if c.Method() == fiber.MethodDelete {
			counter = 0
		}
	}

	if block == "" || block == "counter" {
		c.Locals("Counter", fiber.Map{
			"Value":         counter,
			"IsValueDanger": counter >= 10,
		})
	}

	if block == "" || block == "campaigns" {
		campaigns := []fiber.Map{
			{"ID": rand.Intn(100), "Name": lo.RandomString(4, lo.LowerCaseLettersCharset)},
			{"ID": rand.Intn(100), "Name": lo.RandomString(4, lo.LowerCaseLettersCharset)},
			{"ID": rand.Intn(100), "Name": lo.RandomString(4, lo.LowerCaseLettersCharset)},
			{"ID": rand.Intn(100), "Name": lo.RandomString(4, lo.LowerCaseLettersCharset)},
		}
		c.Locals("Campaigns", campaigns)
	}

	switch block {
	case "campaigns", "counter":
		return c.Render("pages/index/"+block, fiber.Map{})
	default:
		return c.Render("pages/index", fiber.Map{}, "layouts/default")
	}
}
