package yamlreader

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(fileName string) ([]yaml.Node, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var documents []yaml.Node
	decoder := yaml.NewDecoder(file)
	for {
		var document yaml.Node
		if err := decoder.Decode(&document); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		documents = append(documents, document)
	}
	return documents, nil
}
