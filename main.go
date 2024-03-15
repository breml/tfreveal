package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := main0(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main0(osArgs []string) error {
	app := App{}

	cliapp := &cli.App{
		Name:   "tfreveal",
		Usage:  "Show Terraform plan file with all sensitive values revealed.",
		Action: app.TerraformReveal,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-color",
				Usage:       "Disable colorized output",
				Destination: &app.noColor,
			},
		},
	}
	return cliapp.Run(osArgs)
}

type App struct {
	noColor bool
}

func (a *App) TerraformReveal(c *cli.Context) error {
	var source io.Reader = os.Stdin

	if c.Args().Present() {
		file, err := os.Open(c.Args().First())
		if err != nil {
			return fmt.Errorf(`open source: %w`, err)
		}
		defer file.Close()
		source = file
	}

	data, err := io.ReadAll(source)
	if err != nil {
		return fmt.Errorf(`read from source: %w`, err)
	}

	plan := tfjson.Plan{}
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return fmt.Errorf(`unmarshal Terraform plan json: %w`, err)
	}

	fmt.Print(a.resourceChanges(plan))

	fmt.Print(a.outputChanges(plan))

	return nil
}
