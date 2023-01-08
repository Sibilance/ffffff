package ast

type Import struct {
	ASTNode

	Name string
	Path ImportPath
}

type ImportPath []string
