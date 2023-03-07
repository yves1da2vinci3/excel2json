package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tealeg/xlsx"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func main() {

	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
		document, err := c.FormFile("document")
		if err != nil {
			return err
		}

		// Open the Excel file
		file, err := xlsx.OpenFile(document.Filename)
		if err != nil {
			log.Fatal(err)
		}

		// Get the first sheet in the Excel file
		sheet := file.Sheets[0]

		// Extract data from the sheet and convert to list of Person structs
		var people []Person
		for _, row := range sheet.Rows[1:] { // skip header row
			age, err := row.Cells[1].Int()

			if err != nil {
				log.Fatal(err)
			}
			person := Person{
				Name:    row.Cells[0].Value,
				Age:     age,
				Country: row.Cells[2].Value,
			}
			people = append(people, person)
		}

		// Convert slice of Person structs to JSON string

		jsonData, err := json.Marshal(people)
		if err != nil {
			log.Fatal(err)
		}

		return c.Status(fiber.StatusAccepted).JSON(string(jsonData))

	})

	log.Fatal(app.Listen(":3000"))

}
