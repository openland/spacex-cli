package main

import (
	"github.com/openland/spacex-cli/codegen"
	"github.com/openland/spacex-cli/il"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func collectFiles(path string) ([]string, error) {
	var res []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".graphql") {
			return nil
		}
		res = append(res, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "spacex-cli"
	app.HelpName = "spacex-cli"
	app.Version = "1.0.0"

	var output string
	var source string
	var schema string
	var target string
	var pkg string
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate SpaceX client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "output, o",
					Usage:       "Output file name",
					Destination: &output,
				},
				cli.StringFlag{
					Name:        "queries, q",
					Usage:       "Query directory",
					Destination: &source,
				},
				cli.StringFlag{
					Name:        "schema, s",
					Usage:       "Schema path",
					Destination: &schema,
				},
				cli.StringFlag{
					Name:        "target, t",
					Usage:       "Generator target",
					Destination: &target,
				},
				cli.StringFlag{
					Name:        "package, p",
					Usage:       "Package name",
					Destination: &pkg,
				},
			},
			Action: func(c *cli.Context) error {
				if target != "kotlin" && target != "swift" && target != "typescript" && target != "client" {
					log.Panic("Only kotlin, swift, typescript or client target is supported")
				}

				sourcePath, err := filepath.Abs(source)
				if err != nil {
					panic(err)
				}
				schemaPath, err := filepath.Abs(schema)
				if err != nil {
					panic(err)
				}
				files, err := collectFiles(sourcePath)
				if err != nil {
					panic(err)
				}
				model := il.LoadModel(schemaPath, files)
				if target == "kotlin" {
					codegen.GenerateKotlin(model, output, pkg)
				} else if target == "swift" {
					codegen.GenerateSwift(model, output)
				} else if target == "typescript" {
					codegen.GenerateTypescript(model, output)
				} else if target == "client" {
					codegen.GenerateClient(model, output)
				}
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
