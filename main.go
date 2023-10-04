package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PersonInput struct {
	Nickname  string   `json:"apelido"`
	Name      string   `json:"nome"`
	BirthDate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
}

type Person struct {
	ID        uuid.UUID `json:"id"`
	Nickname  string    `json:"apelido"`
	Name      string    `json:"nome"`
	Birthdate time.Time `json:"nascimento"`
	Stack     *string   `json:"stack"`
	Search    string    `json:"search"`
}

func (p PersonInput) validate() error {
	if p.Nickname == "" {
		return errors.New("Nickname is required")
	}

	if p.Name == "" {
		return errors.New("Name is required")
	}

	if p.BirthDate == "" {
		return errors.New("BirthDate is required")
	}
	// TODO: validate that BirthDate is a valid date

	return nil
}

func main() {
	app := fiber.New()

	db, err := sqlx.Connect("postgres", "postgresql://postgres:postgres@postgres:5432/rinha?sslmode=disable")
	if err != nil {
		fmt.Printf("err opening db conn: %v", err.Error())
		return
	}

	app.Get("/", root)
	app.Post("/pessoas", func(c *fiber.Ctx) error {
		return createPerson(c, db)
	})
	app.Get("/pessoas/:id", func(c *fiber.Ctx) error {
		return getPerson(c, db)
	})
	app.Get("/pessoas", func(c *fiber.Ctx) error {
		return searchPerson(c, db)
	})
	app.Get("/contagem-pessoas", func(c *fiber.Ctx) error {
		return getPeopleCount(c, db)
	})

	app.Listen(":3000")

	defer db.Close()
}

func root(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋 from docker!")
}

func createPerson(c *fiber.Ctx, db *sqlx.DB) error {
	body := c.Body()

	person := PersonInput{}
	err := json.Unmarshal(body, &person)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err = person.validate()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	uuid := uuid.New()

	// lower case stack
	lowerCaseStack := make([]string, len(person.Stack))
	for i, s := range person.Stack {
		lowerCaseStack[i] = strings.ToLower(s)
	}

	// convert everything to lower-case and concat all fields to create a search field
	search := fmt.Sprintf(
		"%s %s %s %s",
		strings.ToLower(person.Nickname),
		strings.ToLower(person.Name),
		strings.ToLower(person.BirthDate),
		lowerCaseStack,
	)

	_, err = db.Query(
		`INSERT INTO people 
		(id, nickname, name, birthdate, stack, search) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		uuid, person.Nickname, person.Name, person.BirthDate, pq.Array(person.Stack), search,
	)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"people_nickname_key\"" {
			return c.Status(422).SendString("Nickname already exists")
		}
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}

	c.Location(fmt.Sprintf("/pessoas/%s", uuid))
	return c.Status(201).SendString("Successfully created person")
}

func getPerson(c *fiber.Ctx, db *sqlx.DB) error {
	id := c.Params("id")

	person := Person{}

	err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}

	return c.Status(200).JSON(person)
}

func searchPerson(c *fiber.Ctx, db *sqlx.DB) error {
	t := c.Query("t")
	if t == "" {
		return c.Status(400).SendString("t query param is required")
	}
	t = strings.ToLower(t)

	people := []Person{}

	db.Select(&people, "SELECT * FROM people WHERE search LIKE $1", fmt.Sprintf("%%%s%%", t))

	return c.Status(200).JSON(people)
}

func getPeopleCount(c *fiber.Ctx, db *sqlx.DB) error {
	count := 0

	db.Get(&count, "SELECT COUNT(*) FROM people")

	return c.SendString(fmt.Sprintf("%d", count))
}
