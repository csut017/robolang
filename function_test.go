package robolang

import (
	"encoding/json"
	"testing"
)

func TestFunctionMapToJSON(t *testing.T) {
	fm := FunctionMap{}
	fm["wait"] = FunctionDefinition{}
	out, err := json.Marshal(fm)
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	expected := "[{\"name\":\"wait\"}]"
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
