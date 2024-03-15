package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	plan := tfjson.Plan{}
	err = json.Unmarshal(data, &plan)
	if err != nil {
		panic(err)
	}

	fmt.Println("Changes to Resources:")
	for _, v := range plan.ResourceChanges {
		if v.Change.Actions.NoOp() {
			continue
		}

		diffString := diff(v.Change)

		fmt.Printf("%- 3s %s =%s\n", marker(v.Change.Actions), v.Address, indent(diffString, 2))
	}
	fmt.Println("")

	outputs := make([]string, 0, len(plan.OutputChanges))
	for k := range plan.OutputChanges {
		outputs = append(outputs, k)
	}
	sort.Strings(outputs)

	fmt.Println("Changes to Outputs:")
	for _, k := range outputs {
		v := plan.OutputChanges[k]
		if v.Actions.NoOp() {
			continue
		}

		diffString := diff(v)

		fmt.Printf("%- 3s %s =%s\n", marker(v.Actions), k, indent(diffString, 2))
	}
}

func diff(change *tfjson.Change) string {
	if change.Actions.NoOp() {
		return ""
	}

	maputil.Walk(change.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
		if val, _ := value.(bool); val {
			maputil.DeepSet(change.After, path, "<known after>")
		}
		return nil
	})

	before := []interface{}{change.Before}
	after := []interface{}{change.After}

	d := gojsondiff.New().CompareArrays(before, after)

	formatter := formatter.NewAsciiFormatter(before, formatter.AsciiFormatterConfig{
		ShowArrayIndex: false,
		Coloring:       true,
	})
	diffString, err := formatter.Format(d)
	if err != nil {
		// No error can occur
		panic(err)
	}

	diffString = removeArray(diffString)

	return diffString
}

func removeArray(in string) string {
	start := strings.Index(in, "[")
	end := strings.LastIndex(in, "]")
	// return strings.Trim(in[start+1:end-1], " \n")
	return in[start+1 : end-1]
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

func indent(in string, indent int) string {
	lines := strings.Split(in, "\n")
	for i := range lines {
		lines[i] = strings.Repeat(" ", indent) + lines[i]
	}
	return strings.Join(lines, "\n")
}
