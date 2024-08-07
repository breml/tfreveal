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
	colorize := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: a.noColor,
	}

	buf := strings.Builder{}

	for _, v := range plan.ResourceChanges {
		if v.Change.Actions.NoOp() {
			continue
		}

		diffString := a.diff(v.Change)

		if len(v.Change.ReplacePaths) > 0 {
			buf.WriteString(colorize.Color(fmt.Sprintf("  [white]# %s[reset] must be [light_red]replaced[reset]\n", v.Address)))
		}
		buf.WriteString(fmt.Sprintf("%s %s = %s\n", a.marker(v.Change.Actions), v.Address, indent(diffString, 2)))
	}

	if buf.Len() == 0 {
		return ""
	}

	return fmt.Sprintf("Changes to Resources:\n\n%s", buf.String())
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

		buf.WriteString(fmt.Sprintf("%s %s = %s", a.marker(v.Actions), k, indent(diffString, 2)))
	}

	if buf.Len() == 0 {
		return ""
	}

	return fmt.Sprintf("Changes to Outputs:\n\n%s", buf.String())
}

func (a *App) diff(change *tfjson.Change) string {
	if change.Actions.NoOp() {
		return ""
	}

	err := maputil.Walk(change.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
		if val, _ := value.(bool); val {
			maputil.DeepSet(change.After, path, "(known after apply)")
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

	replacePaths := make([]string, 0, len(change.ReplacePaths))
	for _, replacePath := range change.ReplacePaths {
		rp := replacePath.([]any)
		var path string
		for _, pathSegment := range rp {
			path += "/" + fmt.Sprint(pathSegment)
		}
		replacePaths = append(replacePaths, path)
	}

	buf := &strings.Builder{}
	options := []jsondiffprinter.Option{
		jsondiffprinter.WithTerraformDefaults(),
		jsondiffprinter.WithWriter(buf),
		jsondiffprinter.WithIndentation("    "),
		jsondiffprinter.WithHideUnchanged(!a.showUnchanged),
		jsondiffprinter.WithJSONinJSONCompare(compare),
		jsondiffprinter.WithColor(!a.noColor),
		jsondiffprinter.WithPatchSeriesPostProcess(func(diff jsondiffprinter.Patch) jsondiffprinter.Patch {
			colorize := colorstring.Colorize{
				Colors:  colorstring.DefaultColors,
				Disable: a.noColor,
			}

		outerLoop:
			for _, path := range replacePaths {
				for i := range diff {
					if diff[i].Path.String() == path {
						if diff[i].Metadata == nil {
							diff[i].Metadata = make(map[string]string, 0)
						}
						diff[i].Metadata["note"] = colorize.Color(" [light_red]# forces replacement[reset]")
						continue outerLoop
					}
				}
			}
			return diff
		}),
	}
	err = jsondiffprinter.Format(change.Before, patch, options...)
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
		return colorize.Color("  [green]+[reset]")
	case actions.CreateBeforeDestroy():
		return colorize.Color("[green]+[reset]/[red]-[reset]")
	case actions.Delete():
		return colorize.Color("  [red]-[reset]")
	case actions.DestroyBeforeCreate():
		return colorize.Color("[red]-[reset]/[green]+[reset]")
	case actions.NoOp():
		return colorize.Color("")
	case actions.Read():
		return colorize.Color("")
	case actions.Replace():
		return colorize.Color("[red]-[reset]/[green]+[reset]")
	case actions.Update():
		return colorize.Color("  [yellow]~[reset]")
	default:
		return colorize.Color("[red]!  [reset]")
	}
}
