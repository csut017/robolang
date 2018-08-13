package robolang

import "testing"

func TestFromPraseResult(t *testing.T) {
	parser := NewParser("clear()")
	script := parser.Parse().Script()
	if len(script.Nodes) != 1 {
		t.Errorf("Unexpected number of nodes in the script: expected 1, actual %d", len(script.Nodes))
	}
}

func TestStartTerminatesWithNoNodes(t *testing.T) {
	script := &Script{}
	script.Start()
	expected := ScriptStateFinished
	if script.State != expected {
		t.Errorf("Unexpected script state: expected %s, actual %s", expected.String(), script.State.String())
	}
}
