package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/colorstring"
	"github.com/wI2L/jsondiff"

	"github.com/breml/jsondiffprinter"
)

func (a *App) resourceChanges(plan tfjson.Plan) string {
	buf := strings.Builder{}

	for _, v := range plan.ResourceChanges {
		if v.Change.Actions.NoOp() {
			continue
		}

		diffString := a.diff(v.Change)

		buf.WriteString(fmt.Sprintf("%s %s = %s\n", a.marker(v.Change.Actions), v.Address, indent(diffString, 2)))
	}

	if buf.Len() == 0 {
		return ""
	}

	return fmt.Sprintf("Changes to Resources:\n%s", buf.String())
}

func (a *App) outputChanges(plan tfjson.Plan) string {
	buf := strings.Builder{}

	outputs := make([]string, 0, len(plan.OutputChanges))
	for k := range plan.OutputChanges {
		outputs = append(outputs, k)
	}
	sort.Strings(outputs)

	for _, k := range outputs {
		v := plan.OutputChanges[k]
		if v.Actions.NoOp() {
			continue
		}

		diffString := a.diff(v)

		buf.WriteString(fmt.Sprintf("%s %s = %s\n", a.marker(v.Actions), k, indent(diffString, 2)))
	}

	if buf.Len() == 0 {
		return ""
	}

	return fmt.Sprintf("Changes to Outputs:\n%s", buf.String())
}

func (a *App) diff(change *tfjson.Change) string {
	if change.Actions.NoOp() {
		return ""
	}

	err := maputil.Walk(change.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
		if val, _ := value.(bool); val {
			maputil.DeepSet(change.After, path, "<known after>")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	patch, err := jsondiff.Compare(change.Before, change.After)
	if err != nil {
		panic(err)
	}

	buf := &strings.Builder{}
	formatter := jsondiffprinter.NewTerraformFormatter(buf,
		jsondiffprinter.WithIndentation("    "),
		jsondiffprinter.WithHideUnchanged(true),
		jsondiffprinter.WithJSONinJSONCompare(compare),
		jsondiffprinter.WithColor(!a.noColor),
	)
	err = formatter.Format(change.Before, patch)
	if err != nil {
		panic(err)
	}

	return strings.TrimLeft(buf.String(), " ")
}

func compare(before, after any) ([]byte, error) {
	patch, err := jsondiff.Compare(before, after)
	if err != nil {
		return nil, err
	}
	return json.Marshal(patch)
}

func indent(in string, indent int) string {
	lines := strings.Split(in, "\n")
	for i := range lines[1:] {
		lines[i] = strings.Repeat(" ", indent) + lines[i]
	}
	lines[0] = strings.TrimLeft(lines[0], " ")
	return strings.Join(lines, "\n")
}

func (a *App) marker(actions tfjson.Actions) string {
	colorize := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: a.noColor,
	}
	switch {
	case actions.Create():
		return colorize.Color("  [green][bold]+[reset]")
	case actions.CreateBeforeDestroy():
		return colorize.Color("[green][bold]+[reset]/[red][bold]-[reset]")
	case actions.Delete():
		return colorize.Color("  [red][bold]-[reset]")
	case actions.DestroyBeforeCreate():
		return colorize.Color("[red][bold]-[reset]/[green][bold]+[reset]")
	case actions.NoOp():
		return colorize.Color("")
	case actions.Read():
		return colorize.Color("")
	case actions.Replace():
		return colorize.Color("[red][bold]-[reset]/[green][bold]+[reset]")
	case actions.Update():
		return colorize.Color("  [yellow][bold]~[reset]")
	default:
		return colorize.Color("[red][bold]!  [reset]")
	}
}
