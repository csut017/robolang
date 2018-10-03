package robolang

import (
	"encoding/json"
	"fmt"
	"sort"
)

// VariableTable defines all the available variables for a block
type VariableTable struct {
	Parent    *VariableTable `json:"-"`
	Variables VariableMap    `json:"variables"`
}

// NewVariableTable starts a new variable table
func NewVariableTable(variables ...*VariableDefinition) *VariableTable {
	table := &VariableTable{
		Variables: VariableMap{},
	}
	if variables != nil {
		for _, variable := range variables {
			table.Variables[variable.Name] = variable
		}
	}
	return table
}

// Get attempts to retrieve a value from table
func (table *VariableTable) Get(name string) (*VariableDefinition, bool) {
	value, exists := table.Variables[name]
	if exists {
		return value, true
	}
	if table.Parent != nil {
		return table.Parent.Get(name)
	}
	return nil, false
}

// Add adds a new variable to the table
func (table *VariableTable) Add(name string) (*VariableDefinition, error) {
	_, exists := table.Variables[name]
	if exists {
		return nil, fmt.Errorf("Variable %s already exists", name)
	}
	value := NewVariable(name)
	table.Variables[name] = value
	return value, nil
}

// VariableMap is a convience wrapper to simplify the marshalling and unmarshalling of variable definitions
type VariableMap map[string]*VariableDefinition

// MarshalJSON converts a VariableMap to JSON
func (fm VariableMap) MarshalJSON() ([]byte, error) {
	out, pos := make(VariableDefinitions, len(fm)), 0
	for key, val := range fm {
		val.Name = key
		out[pos] = *val
		pos++
	}
	sort.Sort(out)
	return json.Marshal(out)
}

// UnmarshalJSON converts JSON to a VariableMap
func (fm *VariableMap) UnmarshalJSON(data []byte) error {
	in := VariableDefinitions{}
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	for _, val := range in {
		out := val
		(*fm)[val.Name] = &out
	}
	return nil
}

// VariableDefinitions is a convience wrapper that allows sorting a list of VariableDefinition instances
type VariableDefinitions []VariableDefinition

// Len returns the number of instances in this list
func (vd VariableDefinitions) Len() int {
	return len(vd)
}

// Less checks which name is the lower of two
func (vd VariableDefinitions) Less(i, j int) bool {
	return vd[i].Name < vd[j].Name
}

// Swap changes the location of two variable definitions
func (vd VariableDefinitions) Swap(i, j int) {
	vd[i], vd[j] = vd[j], vd[i]
}

// VariableDefinition defines a variable that holds a value
type VariableDefinition struct {
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}

// NewVariable starts a new variable definition
func NewVariable(name string) *VariableDefinition {
	return &VariableDefinition{Name: name}
}

// Set sets the value of the variable
func (variable *VariableDefinition) Set(value string) *VariableDefinition {
	variable.Value = &value
	return variable
}
