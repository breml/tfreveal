package main

import "testing"

func TestRun_outputs(t *testing.T) {
	err := run("./testdata/outputs/plan.json")
	if err != nil {
		t.Fatal(err)
	}
}
