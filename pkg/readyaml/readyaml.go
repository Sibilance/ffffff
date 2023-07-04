package readyaml

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(fileName string) ([]*yaml.Node, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadStream(file)
}

func ReadStream(stream io.Reader) ([]*yaml.Node, error) {
	var documents []*yaml.Node
	decoder := yaml.NewDecoder(stream)
	for {
		var document yaml.Node
		if err := decoder.Decode(&document); err != nil {
			if err != io.EOF {
				return documents, fmt.Errorf("error decoding document: %w", err)
			}
			break
		}
		documents = append(documents, &document)
	}
	return documents, nil
}
