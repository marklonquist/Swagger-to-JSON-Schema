package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "JSON Schema Generator CLI"
	app.Usage = ""
	app.UsageText = "A CLI used to generate JSON Schema based on swagger:model comments"
	app.Author = "Mark Bøg Lønquist"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "",
			Action: func(c *cli.Context) {
				inputFile := c.String("input-file")
				outputFolder := c.String("output-folder")
				pretty := c.Bool("pretty")
				if inputFile == "" {
					log.Println("Missing input file - should point to a swagger file")
					return
				}
				if outputFolder == "" {
					log.Println("Missing output folder (should point to a folder, overwrites existing files)")
					return
				}

				generateJSONSchema(inputFile, outputFolder, pretty)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input-file",
					Value: "",
				},
				cli.StringFlag{
					Name:  "output-folder",
					Value: "",
				},
				cli.BoolFlag{
					Name: "pretty",
				},
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
