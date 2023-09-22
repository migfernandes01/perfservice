package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"rinha-backend/gen/rinha/public/model"
	"rinha-backend/gen/rinha/public/table"
	"strings"

	jetPg "github.com/go-jet/jet/v2/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Person struct {
	Nickname  string   `json:"apelido"`
	Name      string   `json:"nome"`
	BirthDate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
}

func (p Person) validate() error {
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

	db, err := sql.Open("postgres", "postgresql://postgres:postgres@postgres:5432/rinha?sslmode=disable")
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

func createPerson(c *fiber.Ctx, db *sql.DB) error {
	body := c.Body()

	person := Person{}
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

	is := table.People.INSERT(
		table.People.ID,
		table.People.Nickname,
		table.People.Name,
		table.People.Birthdate,
		table.People.Stack,
		table.People.Search,
	).VALUES(
		uuid,
		person.Nickname,
		person.Name,
		person.BirthDate,
		pq.Array(person.Stack),
		search,
	)

	_, err = is.Exec(db)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"people_nickname_key\"" {
			return c.Status(422).SendString("Nickname already exists")
		}
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}

	c.Location(fmt.Sprintf("/pessoas/%s", uuid))
	return c.Status(201).SendString("Successfully created person")
}

func getPerson(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	person := model.People{}

	ss := table.People.SELECT(
		table.People.AllColumns,
	).FROM(
		table.People,
	).WHERE(
		table.People.ID.EQ(jetPg.UUID(uuid.MustParse(id))),
	)

	err := ss.Query(db, &person)
	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}

	return c.Status(200).JSON(person)
}

func searchPerson(c *fiber.Ctx, db *sql.DB) error {
	t := c.Query("t")
	if t == "" {
		return c.Status(400).SendString("t query param is required")
	}
	t = strings.ToLower(t)

	people := []model.People{}

	ss := table.People.SELECT(
		table.People.ID,
		table.People.Nickname,
		table.People.Name,
		table.People.Birthdate,
		table.People.Stack,
	).FROM(
		table.People,
	).WHERE(
		table.People.Search.LIKE(jetPg.String(fmt.Sprintf("%%%s%%", t))),
	)

	err := ss.Query(db, &people)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}

	return c.Status(200).JSON(people)
}

func getPeopleCount(c *fiber.Ctx, db *sql.DB) error {
	// initialize a int count
	// count := jetPg.COUNT

	// ss := table.People.SELECT(table.People.ID).FROM(table.People)
	ss := table.People.SELECT(jetPg.COUNT(table.People.ID)).FROM(table.People)
	res, err := ss.Exec(db)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	}
	ra, _ := res.RowsAffected()

	// err := ss.Query(db, &count)
	// if err != nil {
	// 	return c.Status(500).SendString(fmt.Sprintf("err executing query: %v", err.Error()))
	// }
	return c.SendString(fmt.Sprintf("%d", ra))
}
