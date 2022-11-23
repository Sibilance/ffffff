package yamlreader

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

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

func nameFromPath(fileName string, isFile bool) string {
	name := path.Base(fileName)

	// Remove .yaml file extension if it exists.
	if isFile && strings.HasSuffix(name, ".yaml") {
		name = name[:len(name)-5]
	}

	return name
}
