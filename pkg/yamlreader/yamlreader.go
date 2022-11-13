package yamlreader

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(fileName string) (node Node, err error) {
	node.FileName = fileName
	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	// Read all the documents from the yaml file.
	var documents []Node
	decoder := yaml.NewDecoder(file)
	for {
		var document Node
		if err = decoder.Decode(&document); err != nil {
			if err != io.EOF {
				return
			}
			break
		}
		documents = append(documents, document)
	}
	return
}
