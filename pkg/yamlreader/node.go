package yamlreader

import "fmt"

type Node struct {
	Name     string
	FileName string
	Line     int
	Column   int
	Comment  string
	Raw      string
	Kind     Kind
	Tag      string
	Sequence []Node
	Mapping  map[string]Node
	Scalar
}

func (n *Node) String() string {
	return fmt.Sprintf("%s (%s:%d:%d)", n.Name, n.FileName, n.Line, n.Column)
}
