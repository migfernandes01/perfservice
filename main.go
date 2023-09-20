package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", root)
	app.Post("/pessoas", createPerson)
	app.Get("/pessoas/:id", getPerson)
	app.Get("/pessoas", searchPerson)
	app.Get("/contagem-pessoas", getPeopleCount)

	app.Listen(":3000")
}

func root(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋 from docker!")
}

func createPerson(c *fiber.Ctx) error {
	return c.SendString("Create person route")
}

func getPerson(c *fiber.Ctx) error {
	return c.SendString("Get person route")
}

func searchPerson(c *fiber.Ctx) error {
	return c.SendString("Search person route")
}

func getPeopleCount(c *fiber.Ctx) error {
	return c.SendString("Get people count route")
}
