package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/colorstring"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := main0(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main0(osArgs []string) error {
	app := App{}

	executionPlanLegend := colorstring.Color(`
Resource actions are indicated with the following symbols:
   [green]+[reset] create
   [yellow]~[reset] update in-place
   [red]-[reset] destroy
 [red]-[reset]/[green]+[reset] destroy and then create replacement
 [green]+[reset]/[red]-[reset] create replacement and then destroy
`)

	cliapp := &cli.App{
		Name:   "tfreveal",
		Usage:  "Show an execution plan with all sensitive values revealed.",
		Action: app.Reveal,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-color",
				Usage:       "Disable colorized output",
				Destination: &app.noColor,
			},
			&cli.BoolFlag{
				Name:        "show-unchanged",
				Usage:       "Show unchanged attributes",
				Destination: &app.showUnchanged,
			},
		},
		CustomAppHelpTemplate: cli.AppHelpTemplate + executionPlanLegend,

		Copyright: "© 2024, Lucas Bremgartner",
		Version:   fmt.Sprintf("%s (%s, %s)", version, commit, date),
	}
	return cliapp.Run(osArgs)
}

type App struct {
	noColor       bool
	showUnchanged bool
}

func (a *App) Reveal(c *cli.Context) error {
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

	fmt.Println("The provided execution plan contains the following changes.")
	fmt.Println()
	fmt.Print(a.resourceChanges(plan))
	fmt.Print(a.outputChanges(plan))

	return nil
}
