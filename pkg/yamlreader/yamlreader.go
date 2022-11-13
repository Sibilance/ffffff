package yamlreader

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(name string, fileName string) (n *Node, err error) {
	n = &Node{
		Name:     name,
		FileName: fileName,
	}

	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	// Read all the documents from the yaml file.
	var documents []*yaml.Node
	decoder := yaml.NewDecoder(file)
	for {
		var document yaml.Node
		if err = decoder.Decode(&document); err != nil {
			if err != io.EOF {
				err = fmt.Errorf("%s: %s", n, err)
				return
			}
			break
		}
		documents = append(documents, &document)
	}

	if len(documents) == 1 {
		err = n.UnmarshalYAML(documents[0])
	} else if len(documents) > 1 {
		n.Kind = SequenceNode
		for i, document := range documents {
			innerNode := Node{
				Name:     fmt.Sprintf("%s[%d]", name, i),
				FileName: fileName,
			}
			if err = innerNode.UnmarshalYAML(document); err != nil {
				return
			}
			n.Sequence = append(n.Sequence, innerNode)
		}
	} else {
		n.Kind = ScalarNode
		n.Tag = "!!null"
	}

	return
}
