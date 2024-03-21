package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/breml/tfreveal/formatter"
	"github.com/breml/tfreveal/gojsondiff"
	origformatter "github.com/breml/tfreveal/gojsondiff/formatter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "provide terraform plan in JSON format as first argument")
		os.Exit(1)
	}
	err := run(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run(args string) error {
	data, err := os.ReadFile(args)
	if err != nil {
		return err
	}

	plan := tfjson.Plan{}
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return err
	}

	if hasResourceChanged(plan.ResourceChanges) {
		fmt.Println("Changes to Resources:")
		for _, v := range plan.ResourceChanges {
			if v.Change.Actions.NoOp() {
				continue
			}

			maputil.Walk(v.Change.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
				if val, _ := value.(bool); val {
					maputil.DeepSet(v.Change.After, path, "<known after>")
				}
				return nil
			})

			key := fmt.Sprintf("resource %q %q", v.Type, v.Name)

			before := map[string]interface{}{key: v.Change.Before}
			after := map[string]interface{}{key: v.Change.After}

			diffString, err := formatDiff(before, after)
			if err != nil {
				return err
			}

			fmt.Printf("%- 3s %s\n", marker(v.Change.Actions), diffString)
		}
		fmt.Println("")
	}

	if hasOutputChanged(plan.OutputChanges) {
		fmt.Println("Changes to Outputs:")
		outputsBefore := map[string]interface{}{}
		outputsAfter := map[string]interface{}{}
		for k, v := range plan.OutputChanges {
			if v.Actions.NoOp() {
				continue
			}

			maputil.Walk(v.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
				if val, _ := value.(bool); val {
					maputil.DeepSet(v.After, path, "<known after>")
				}
				return nil
			})

			outputsBefore[k] = v.Before
			outputsAfter[k] = v.After
		}

		diffString, err := formatDiff(outputsBefore, outputsAfter)
		if err != nil {
			return err
		}

		fmt.Println("  ", diffString)
		fmt.Println("")
	}

	return nil
}

func hasResourceChanged(changes []*tfjson.ResourceChange) bool {
	for _, change := range changes {
		if !change.Change.Actions.NoOp() {
			return true
		}
	}
	return false
}

func hasOutputChanged(changes map[string]*tfjson.Change) bool {
	for _, change := range changes {
		if !change.Actions.NoOp() {
			return true
		}
	}
	return false
}

func formatDiff(before, after map[string]interface{}) (string, error) {
	d := gojsondiff.New().CompareObjects(before, after)

	var f interface {
		Format(diff gojsondiff.Diff) (result string, err error)
	}
	if true {
		f = formatter.NewTerraformFormatter(before, formatter.TerraformFormatterConfig{
			Coloring: true,
		})
	} else {
		f = origformatter.NewAsciiFormatter(before, origformatter.AsciiFormatterConfig{
			Coloring: true,
		})
	}
	diffString, err := f.Format(d)
	if err != nil {
		return "", err
	}

	diffString = formatCleanup(diffString)

	return diffString, nil
}

func formatCleanup(in string) string {
	start := strings.Index(in, "{")
	end := strings.LastIndex(in, "}")

	return strings.Trim(in[start+1:end-1], " \n")
}

func marker(actions tfjson.Actions) string {
	switch {
	case actions.Create():
		return "+"
	case actions.CreateBeforeDestroy():
		return "+/-"
	case actions.Delete():
		return "-"
	case actions.DestroyBeforeCreate():
		return "-/+"
	case actions.NoOp():
		return ""
	case actions.Read():
		return ""
	case actions.Replace():
		return "-/+"
	case actions.Update():
		return "~"
	default:
		return "!  "
	}
}
