package robolang

import (
	"encoding/json"
	"testing"
)

func TestFunctionMapToJSON(t *testing.T) {
	fm := FunctionMap{}
	fm["wait"] = &FunctionDefinition{}
	fm["show"] = &FunctionDefinition{}
	fm["repeat"] = &FunctionDefinition{}
	out, err := json.Marshal(fm)
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	expected := "[{\"name\":\"repeat\"},{\"name\":\"show\"},{\"name\":\"wait\"}]"
	if string(out) != expected {
		t.Errorf("Unexpected JSON output: expected %s, actual %s", expected, out)
	}
}

func TestFunctionMapFromJSON(t *testing.T) {
	fm := FunctionMap{}
	in := "[{\"name\":\"wait\"}]"
	err := json.Unmarshal([]byte(in), &fm)
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	if _, ok := fm["wait"]; !ok {
		t.Errorf("Missing value 'wait', actual is %+v", fm)
	}
}

func TestFunctionGet(t *testing.T) {
	base := NewFunctionTable(
		NewFunction("wait"))
	child := NewFunctionTable(
		NewFunction("doCheck"))
	child.Parent = base

	tests := []struct {
		Key    string
		Exists bool
		Table  *FunctionTable
	}{
		{Key: "nothing", Exists: false, Table: base},
		{Key: "wait", Exists: true, Table: base},
		{Key: "wait", Exists: true, Table: child},
		{Key: "doCheck", Exists: true, Table: child},
		{Key: "neither", Exists: false, Table: child},
	}
	for count, test := range tests {
		t.Logf("Running test #%d", count+1)
		_, found := test.Table.Get(test.Key)
		if test.Exists != found {
			t.Errorf("Exists condition does not match: expected %v, actual %v", test.Exists, found)
		}
	}
}

func TestFunctionAdd(t *testing.T) {
	table := NewFunctionTable(NewFunction("wait"))
	tests := []struct {
		Key     string
		Success bool
	}{
		{"doCheck", true},
		{"wait", false},
	}

	for _, test := range tests {
		_, err := table.Add(test.Key)
		if err == nil {
			if !test.Success {
				t.Errorf("Expected an error, not nil")
			}
		} else {
			if test.Success {
				t.Errorf("Unexpected error: %v", err)
			}
		}
	}
}
