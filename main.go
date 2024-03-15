package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghetzel/go-stockutil/maputil"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/yudai/gojsondiff"

	"github.com/breml/tfreveal/formatter"
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

	for _, v := range plan.ResourceChanges {
		if v.Change.Actions.NoOp() {
			continue
		}

		fmt.Println("Resource: ", v.Address, v.Change.Actions)

		maputil.Walk(v.Change.AfterUnknown, func(value interface{}, path []string, isLeaf bool) error {
			if val, _ := value.(bool); val {
				maputil.DeepSet(v.Change.After, path, "<known after>")
			}
			return nil
		})

		differ := gojsondiff.New()
		// FIXME: support different types, check before compare
		d := differ.CompareObjects(v.Change.Before.(map[string]interface{}), v.Change.After.(map[string]interface{}))

		formatter := formatter.NewTerraformFormatter(v.Change.Before, formatter.TerraformFormatterConfig{
			ShowArrayIndex: false,
			Coloring:       true,
		})
		diffString, err := formatter.Format(d)
		if err != nil {
			// No error can occur
			panic(err)
		}

		fmt.Println(diffString)
	}
}
