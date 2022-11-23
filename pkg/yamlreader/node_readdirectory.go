package yamlreader

import (
	"os"
	"path"
	"strings"
)

func (n *Node) ReadDirectory(dirName string) error {
	n.FileName = dirName
	n.Kind = MappingNode
	n.Tag = MapTag
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
