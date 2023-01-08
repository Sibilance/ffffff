package ast

type Expression interface{}

type AssignmentExpression struct {
	ASTNode

	Target VariableExpression
	Value  Expression
}

type CallExpression struct {
	ASTNode

	Callable            Expression
	PositionalArguments []Expression
	KeywordArguments    map[string]Expression
}

type VariableExpression struct {
	ASTNode

	Name string
}
