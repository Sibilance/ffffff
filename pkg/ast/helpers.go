package ast

func deleteMissingKeys[T interface{}](n Node, m map[string]T, t Tag) {
	nodeMapping := n.AsMapping()
	for key := range m {
		if innerNode, ok := nodeMapping[key]; !ok || (t != "" && innerNode.Tag() != t) {
			delete(m, key)
		}
	}
}

type parsable interface {
	Parse(n Node)
}

func parseByTag[T any, TPtr interface {
	*T
	parsable
}](n Node, m map[string]T, t Tag) {
	deleteMissingKeys(n, m, t)
	for key, innerNode := range n.AsMapping() {
		if innerNode.Tag() == t {
			value := m[key]
			if value == nil {
				value = new(T)
			}
			value.Parse(innerNode)
			m[key] = value
		}
	}
}
