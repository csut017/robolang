package robolang

import (
	"encoding/json"
	"testing"
)

func TestVariableMapToJSON(t *testing.T) {
	fm := VariableMap{}
	fm["wait"] = NewVariable("wait")
	fm["hello"] = NewVariable("hello").Set("world")
	out, err := json.Marshal(fm)
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	expected := "[{\"name\":\"wait\"},{\"name\":\"hello\",\"value\":\"world\"}]"
	if string(out) != expected {
		t.Errorf("Unexpected JSON output: expected %s, actual %s", expected, out)
	}
}

func TestVariableMapFromJSON(t *testing.T) {
	fm := VariableMap{}
	in := "[{\"name\":\"wait\"},{\"name\":\"hello\",\"value\":\"world\"}]"
	err := json.Unmarshal([]byte(in), &fm)
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	val, ok := fm["hello"]
	if !ok {
		t.Errorf("Missing value 'hello', actual is %+v", fm)
	}
	expected := "world"
	if val.Value == nil {
		t.Errorf("Incorrect value for 'hello': expected %s, actual <nil>", expected)
	} else if *val.Value != expected {
		t.Errorf("Incorrect value for 'hello': expected %s, actual %s", expected, *val.Value)
	}
}

func TestVariableGet(t *testing.T) {
	base := NewVariableTable(
		NewVariable("test"),
		NewVariable("parent").Set("only"),
		NewVariable("hello").Set("world"))
	child := NewVariableTable(
		NewVariable("hello").Set("there"),
		NewVariable("me").Set("too"))
	child.Parent = base

	vals := []string{"world", "there", "only"}
	tests := []struct {
		Key      string
		Exists   bool
		Expected *string
		Table    *VariableTable
	}{
		{Key: "nothing", Exists: false, Expected: nil, Table: base},
		{Key: "test", Exists: true, Expected: nil, Table: base},
		{Key: "hello", Exists: true, Expected: &vals[0], Table: base},
		{Key: "hello", Exists: true, Expected: &vals[1], Table: child},
		{Key: "parent", Exists: true, Expected: &vals[2], Table: child},
		{Key: "neither", Exists: false, Expected: nil, Table: child},
	}
	for count, test := range tests {
		t.Logf("Running test #%d", count+1)
		value, found := test.Table.Get(test.Key)
		if test.Exists != found {
			t.Errorf("Exists condition does not match: expected %v, actual %v", test.Exists, found)
		}
		if value != nil {
			if test.Expected == nil || value.Value == nil {
				if test.Expected != value.Value {
					t.Errorf("Value condition does not match: expected %v, actual %v", test.Expected, value.Value)
				}
			} else if *test.Expected != *value.Value {
				t.Errorf("Value condition does not match: expected %s, actual %s", *test.Expected, *value.Value)
			}
		}
	}
}

func TestVariableAdd(t *testing.T) {
	table := NewVariableTable(NewVariable("Exists"))
	tests := []struct {
		Key     string
		Success bool
	}{
		{"Missing", true},
		{"Exists", false},
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
