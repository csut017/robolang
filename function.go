package robolang

import "encoding/json"

// Function is a block of functionality that can be executed
type Function interface {
	Start()
	Resume()
}

// FunctionTable defines all the available functions for a block
type FunctionTable struct {
	Parent    *FunctionTable `json:"parent,omitempty"`
	Functions FunctionMap    `json:"functions"`
}

// FunctionMap is a convience wrapper to simplify the marshalling and unmarshalling of function definitions
type FunctionMap map[string]FunctionDefinition

// MarshalJSON converts a FunctionMap to JSON
func (fm FunctionMap) MarshalJSON() ([]byte, error) {
	out, pos := make([]FunctionDefinition, len(fm)), 0
	for key, val := range fm {
		val.Name = key
		out[pos] = val
		pos++
	}
	return json.Marshal(out)
}

// UnmarshalJSON converts JSON to a FunctionMap
func (fm *FunctionMap) UnmarshalJSON(data []byte) error {
	in := []FunctionDefinition{}
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	for _, val := range in {
		(*fm)[val.Name] = val
	}
	return nil
}

// FunctionDefinition defines a function that can be executed in a block
type FunctionDefinition struct {
	Name       string   `json:"name"`
	Definition *Node    `json:"definition,omitempty"`
	Function   Function `json:"-"`
}
