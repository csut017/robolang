package robolang

// Script defines the execution environment for a script
type Script struct {
	Nodes []*Node
	State ScriptState

	current *scriptNode
	nodeMap map[int]*scriptNode
}

// Start begins executing the script
func (s *Script) Start() error {
	s.nodeMap = map[int]*scriptNode{}
	s.initialiseNodes(&s.Nodes, 0, nil)
	first, ok := s.nodeMap[0]
	if !ok {
		s.State = ScriptStateFinished
		return nil
	}

	s.current = first
	return s.executeLoop()
}

func (s *Script) executeLoop() error {
	err := s.executeNode()
	if err != nil {
		s.State = ScriptStateFailed
		return err
	}

	// next := s.current.next

	return nil
}

func (s *Script) executeNode() error {
	return nil
}

func (s *Script) initialiseNodes(nodes *[]*Node, id int, parent *scriptNode) (*scriptNode, int) {
	var last, first *scriptNode
	for _, node := range *nodes {
		this := &scriptNode{
			id:     id,
			node:   *node,
			parent: parent,
		}
		if first == nil {
			first = this
		}
		id++
		s.nodeMap[id] = this
		if last != nil {
			last.next = this
		}
		last = this

		this.firstArg, id = s.initialiseNodes(&this.node.Args, id, this)
		this.firstChild, id = s.initialiseNodes(&this.node.Children, id, this)
	}

	return first, id
}

type scriptNode struct {
	firstArg   *scriptNode
	firstChild *scriptNode
	id         int
	node       Node
	next       *scriptNode
	parent     *scriptNode
}

// ScriptState defines the current state of the script
type ScriptState int

//go:generate stringer -type=ScriptState

const (
	// ScriptStatePending means the script has not started execution yet
	ScriptStatePending ScriptState = iota

	// ScriptStateWaiting means the script is waiting for external input
	ScriptStateWaiting

	// ScriptStateFinished means the script has finished execution successfully
	ScriptStateFinished

	// ScriptStateFailed means the script has failed at some point in its execution
	ScriptStateFailed
)
