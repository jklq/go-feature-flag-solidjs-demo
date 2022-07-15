package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

type user struct {
	userId string
	value  string
	ffuser ffuser.User
}

var users [10000]user

func main() {
	app := fiber.New()

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/", "./public")

	err := ffclient.Init(ffclient.Config{
		PollingInterval: 500 * time.Microsecond,
		Context:         context.Background(),
		Retriever: &fileretriever.Retriever{
			Path: "flags.yaml",
		},
	})

	for i := 0; i < 10000; i++ {
		id := i
		u := ffuser.NewUser(fmt.Sprint(id))

		users[i].ffuser = u
	}

	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		var mapToRender [10000]string
		for i, user := range users {
			color, _ := ffclient.StringVariation("color-flag", user.ffuser, "grey")
			mapToRender[i] = color
		}
		jsonMapToRender, err := json.Marshal(mapToRender)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendString(string(jsonMapToRender))
	})

	app.Listen(":3030")
}
