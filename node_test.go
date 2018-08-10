package robolang

import "testing"

func TestNodeBasicString(t *testing.T) {
	node := makeNode(NodeFunction, TokenIdentifier, "test")
	actual, expected := node.String(), "[NodeFunction]test()"
	if actual != expected {
		t.Errorf("Node.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}

func TestNodeStringWithArg(t *testing.T) {
	node := makeNode(NodeFunction, TokenIdentifier, "test").
		AddArgument(makeNode(NodeArgument, TokenIdentifier, "value"))
	actual, expected := node.String(), "[NodeFunction]test([NodeArgument]value())"
	if actual != expected {
		t.Errorf("Node.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}

func TestNodeStringWithArgs(t *testing.T) {
	node := makeNode(NodeFunction, TokenIdentifier, "test").
		AddArgument(makeNode(NodeArgument, TokenIdentifier, "first")).
		AddArgument(makeNode(NodeArgument, TokenIdentifier, "second"))
	actual := node.String()
	expected := "[NodeFunction]test([NodeArgument]first(),[NodeArgument]second())"
	if actual != expected {
		t.Errorf("Node.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}

func TestNodeStringWithChild(t *testing.T) {
	node := makeNode(NodeFunction, TokenIdentifier, "test").
		AddChild(makeNode(NodeFunction, TokenIdentifier, "child"))
	actual, expected := node.String(), "[NodeFunction]test()->([NodeFunction]child())"
	if actual != expected {
		t.Errorf("Node.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}

func TestNodeStringWithChildren(t *testing.T) {
	node := makeNode(NodeFunction, TokenIdentifier, "test").
		AddChild(makeNode(NodeFunction, TokenIdentifier, "first")).
		AddChild(makeNode(NodeFunction, TokenIdentifier, "second"))
	actual, expected := node.String(), "[NodeFunction]test()->([NodeFunction]first(),[NodeFunction]second())"
	if actual != expected {
		t.Errorf("Node.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}

func makeNode(nodeType NodeType, tokenType TokenType, value string) *Node {
	return &Node{
		Type:  nodeType,
		Token: &Token{Type: tokenType, Value: value},
	}
}
