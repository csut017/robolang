package robolang

import (
	"strings"
)

// Node is a node in the AST
type Node struct {
	Token    *Token   `json:"token"`
	Type     NodeType `json:"-"`
	TypeText string   `json:"type"`
	Args     []*Node  `json:"args,omitempty"`
	Children []*Node  `json:"children,omitempty"`
}

// String converts the node to a human-readable form.
func (n *Node) String() string {
	out := "[" + n.Type.String() + "]"
	if n.Token == nil {
		return out
	}

	out += n.Token.Value + "("
	if len(n.Args) > 0 {
		args := make([]string, len(n.Args))
		for pos, arg := range n.Args {
			args[pos] = arg.String()
		}
		out += strings.Join(args, ",")
	}
	out += ")"

	if len(n.Children) > 0 {
		children := make([]string, len(n.Children))
		for pos, child := range n.Children {
			children[pos] = child.String()
		}
		out += "->(" + strings.Join(children, ",") + ")"
	}

	return out
}

// AddArgument adds a new argument to the node
func (n *Node) AddArgument(arg *Node) *Node {
	n.Args = append(n.Args, arg)
	return n
}

// AddChild adds a new child to the node
func (n *Node) AddChild(arg *Node) *Node {
	n.Children = append(n.Children, arg)
	return n
}

// NodeType defines the type of node
type NodeType int

//go:generate stringer -type=NodeType

const (
	// NodeInvalid means the node is not valid for its location in the AST
	NodeInvalid NodeType = iota

	// NodeFunction means this node should be executed as a function
	NodeFunction

	// NodeArgument means this node is an argument to a function
	NodeArgument

	// NodeConstant means this node is a constant value
	NodeConstant

	// NodeResource means this node points to a resource
	NodeResource
)
