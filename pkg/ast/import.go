package ast

type Import struct {
	Node

	Name string
	Path ImportPath
}

type ImportPath []string
