package robolang

import (
	"encoding/json"
	"fmt"
	"sort"
)

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

// NewFunctionTable starts a new function table
func NewFunctionTable(functions ...*FunctionDefinition) *FunctionTable {
	table := &FunctionTable{
		Functions: FunctionMap{},
	}
	if functions != nil {
		for _, function := range functions {
			table.Functions[function.Name] = function
		}
	}
	return table
}

// Get attempts to retrieve a value from table
func (table *FunctionTable) Get(name string) (*FunctionDefinition, bool) {
	value, exists := table.Functions[name]
	if exists {
		return value, true
	}
	if table.Parent != nil {
		return table.Parent.Get(name)
	}
	return nil, false
}

// Add adds a new function to the table
func (table *FunctionTable) Add(name string) (*FunctionDefinition, error) {
	_, exists := table.Functions[name]
	if exists {
		return nil, fmt.Errorf("Function %s already exists", name)
	}
	value := NewFunction(name)
	table.Functions[name] = value
	return value, nil
}

// FunctionMap is a convience wrapper to simplify the marshalling and unmarshalling of function definitions
type FunctionMap map[string]*FunctionDefinition

// MarshalJSON converts a FunctionMap to JSON
func (fm FunctionMap) MarshalJSON() ([]byte, error) {
	out, pos := make(FunctionDefinitions, len(fm)), 0
	for key, val := range fm {
		val.Name = key
		out[pos] = *val
		pos++
	}
	sort.Sort(out)
	return json.Marshal(out)
}

// UnmarshalJSON converts JSON to a FunctionMap
func (fm *FunctionMap) UnmarshalJSON(data []byte) error {
	in := FunctionDefinitions{}
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	for _, val := range in {
		function := val
		(*fm)[val.Name] = &function
	}
	return nil
}

// FunctionDefinitions is a convience wrapper that allows sorting a list of FunctionDefinition instances
type FunctionDefinitions []FunctionDefinition

// Len returns the number of instances in this list
func (vd FunctionDefinitions) Len() int {
	return len(vd)
}

// Less checks which name is the lower of two
func (vd FunctionDefinitions) Less(i, j int) bool {
	return vd[i].Name < vd[j].Name
}

// Swap changes the location of two Function definitions
func (vd FunctionDefinitions) Swap(i, j int) {
	vd[i], vd[j] = vd[j], vd[i]
}

// FunctionDefinition defines a function that can be executed in a block
type FunctionDefinition struct {
	Name       string   `json:"name"`
	Definition *Node    `json:"definition,omitempty"`
	Function   Function `json:"-"`
}

// NewFunction starts a new function definition
func NewFunction(name string) *FunctionDefinition {
	return &FunctionDefinition{Name: name}
}
