package yamlreader

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Kind int8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	BooleanNode
	IntegerNode
	FloatNode
	StringNode
	NullNode

	SeqTag   string = "!!seq"
	MapTag   string = "!!map"
	BoolTag  string = "!!bool"
	IntTag   string = "!!int"
	FloatTag string = "!!float"
	StrTag   string = "!!str"
	NullTag  string = "!!null"
)

func (kind Kind) String() string {
	switch kind {
	case UndefinedNode:
		return "UndefinedNode"
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case BooleanNode:
		return "BooleanNode"
	case IntegerNode:
		return "IntegerNode"
	case FloatNode:
		return "FloatNode"
	case StringNode:
		return "StringNode"
	case NullNode:
		return "NullNode"
	}
	return "unknown"
}

type Value interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
}

type DefaultValue struct{}

func (v DefaultValue) Bool() (r bool) {
	return
}

func (v DefaultValue) Int() (r int64) {
	return
}

func (v DefaultValue) Float() (r float64) {
	return
}

func (v DefaultValue) Str() (r string) {
	return
}

type BoolValue struct {
	DefaultValue
	Value bool
}

func (v BoolValue) Bool() bool {
	return v.Value
}

type IntValue struct {
	DefaultValue
	Value int64
}

func (v IntValue) Int() int64 {
	return v.Value
}

type FloatValue struct {
	DefaultValue
	Value float64
}

func (v FloatValue) Float() float64 {
	return v.Value
}

type StringValue struct {
	DefaultValue
	Value string
}

func (v StringValue) String() string {
	return v.Value
}

type Node struct {
	Name     string
	FileName string
	Line     int
	Column   int
	Comment  string
	Kind     Kind
	Tag      string
	Sequence []Node
	Mapping  map[string]Node
	Raw      string
	Value
}

func (n *Node) String() string {
	return fmt.Sprintf("%s (%s:%d:%d)", n.Name, n.FileName, n.Line, n.Column)
}

func (n *Node) ReadDirectory(dirName string) error {
	n.FileName = dirName
	n.Kind = MappingNode
	n.Mapping = make(map[string]Node)

	if n.Name == "" {
		n.Name = nameFromPath(dirName, false)
	}

	dirEntries, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		childName := dirEntry.Name()

		// Ignore hidden files.
		if strings.HasPrefix(childName, ".") {
			continue
		}

		childPath := path.Join(dirName, childName)
		child := Node{}

		// If the child is a directory, load it recursively.
		if dirEntry.IsDir() {
			child.Name = n.Name + "." + childName
			err := child.ReadDirectory(childPath)
			if err != nil {
				return err
			}

			// Prune empty directories.
			if len(child.Mapping) == 0 {
				continue
			}
		} else if strings.HasSuffix(childPath, ".yaml") {
			childName = childName[:len(childName)-5]
			child.Name = n.Name + "." + childName

			err := child.ReadFile(childPath)
			if err != nil {
				return err
			}
		} else {
			continue // Nothing to do here.
		}

		n.Mapping[childName] = child
	}

	return nil
}

func (n *Node) ReadFile(fileName string) error {
	n.FileName = fileName

	if n.Name == "" {
		n.Name = nameFromPath(fileName, true)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	return n.ReadStream(file)
}

func (n *Node) ReadStream(file io.Reader) error {
	// Read all the documents from the yaml file.
	var documents []*yaml.Node
	decoder := yaml.NewDecoder(file)
	for {
		var document yaml.Node
		if err := decoder.Decode(&document); err != nil {
			if err != io.EOF {
				return fmt.Errorf("%s: %s", n, err)
			}
			break
		}
		documents = append(documents, &document)
	}

	if len(documents) == 1 {
		if err := n.UnmarshalYAML(documents[0]); err != nil {
			return err
		}
	} else if len(documents) > 1 {
		n.Kind = SequenceNode
		n.Tag = SeqTag
		for i, document := range documents {
			innerNode := Node{
				Name:     fmt.Sprintf("%s[%d]", n.Name, i),
				FileName: n.FileName,
			}
			if err := innerNode.UnmarshalYAML(document); err != nil {
				return err
			}
			n.Sequence = append(n.Sequence, innerNode)
		}
	} else {
		n.Kind = NullNode
		n.Tag = NullTag
	}

	return nil
}

func (n *Node) UnmarshalYAML(yamlNode *yaml.Node) error {
	recursionDetection := make(map[*yaml.Node]struct{})
	return n.unmarshalYAML(yamlNode, recursionDetection)
}

func (n *Node) unmarshalYAML(yamlNode *yaml.Node, recursionDetection map[*yaml.Node]struct{}) error {
	n.Tag = yamlNode.ShortTag()
	n.Line = yamlNode.Line
	n.Column = yamlNode.Column
	n.Comment = strings.Trim(yamlNode.HeadComment+"\n\n"+yamlNode.LineComment, "\n")

	switch yamlNode.Kind {
	case yaml.DocumentNode:
		if len(yamlNode.Content) != 1 {
			return fmt.Errorf("%s: expected one node in document, got %d", n, len(yamlNode.Content))
		}
		return n.unmarshalYAML(yamlNode.Content[0], recursionDetection)

	case yaml.AliasNode:
		// If we've encountered this particular AliasNode before, it must be a recursive alias.
		// Note that we track the AliasNode itself, not the node it is aliased to, as multiple
		// pointers to the same node are valid and don't indicate recursion.
		if _, contained := recursionDetection[yamlNode]; contained {
			return fmt.Errorf("%s: recursive alias detected", n)
		}
		recursionDetection[yamlNode] = struct{}{}
		// Follow the alias, but keep these attributes from the alias node.
		line, column, comment := n.Line, n.Column, n.Comment
		err := n.unmarshalYAML(yamlNode.Alias, recursionDetection)
		n.Line, n.Column, n.Comment = line, column, comment
		return err

	case yaml.SequenceNode:
		n.Kind = SequenceNode
		for i, innerYamlNode := range yamlNode.Content {
			innerNode := Node{
				Name:     fmt.Sprintf("%s[%d]", n.Name, i),
				FileName: n.FileName,
			}
			if err := innerNode.unmarshalYAML(innerYamlNode, recursionDetection); err != nil {
				return err
			}
			n.Sequence = append(n.Sequence, innerNode)
		}

	case yaml.MappingNode:
		n.Kind = MappingNode
		n.Mapping = make(map[string]Node)
		mapping := map[string]yaml.Node{}
		if err := yamlNode.Decode(mapping); err != nil {
			return fmt.Errorf("%s: %s", n, err)
		}
		for k, innerYamlNode := range mapping {
			innerNode := Node{
				Name:     fmt.Sprintf("%s.%s", n.Name, k),
				FileName: n.FileName,
			}
			if err := innerNode.unmarshalYAML(&innerYamlNode, recursionDetection); err != nil {
				return err
			}
			n.Mapping[k] = innerNode
		}

	case yaml.ScalarNode:
		n.Raw = yamlNode.Value
		switch n.Tag {
		case BoolTag:
			n.Kind = BooleanNode
			value := BoolValue{}
			yamlNode.Decode(&value.Value)
			n.Value = value
		case IntTag:
			n.Kind = IntegerNode
			// We have to check for overflows of positive integers because the yaml
			// library detects anything that fits in int64 or uint64 as !!int.
			var overflowCheck uint64
			yamlNode.Decode(&overflowCheck)
			if overflowCheck > ^uint64(0)>>1 {
				n.Kind = FloatNode
				n.Tag = FloatTag
				value := FloatValue{}
				yamlNode.Decode(&value.Value)
				n.Value = value
			} else {
				value := IntValue{}
				yamlNode.Decode(&value.Value)
				n.Value = value
			}
		case FloatTag:
			n.Kind = FloatNode
			value := FloatValue{}
			yamlNode.Decode(&value.Value)
			n.Value = value
		case NullTag:
			n.Kind = NullNode
		default:
			n.Kind = StringNode
			n.Value = StringValue{Value: n.Raw}
		}

	default:
		return fmt.Errorf("%s: unknown node kind", n)
	}

	return nil
}

func nameFromPath(fileName string, isFile bool) string {
	name := path.Base(fileName)

	// Remove .yaml file extension if it exists.
	if isFile && strings.HasSuffix(strings.ToLower(name), ".yaml") {
		name = name[:len(name)-5]
	}

	return name
}
