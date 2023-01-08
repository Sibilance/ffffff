package ast

type Statement interface{}

type BreakStatement struct {
	Node

	Label string
}

type ContinueStatement struct {
	Node

	Label string
}

type ExpressionStatement struct {
	Node

	Expression Expression
}

type ForLoopStatement struct {
	Node

	Label    string
	Target   VariableExpression
	Iterable Expression
	Body     CodeBlock
	Else     CodeBlock
}

type IfStatement struct {
	Node

	Conditionals []IfStatementConditional
	Else         CodeBlock
}

type IfStatementConditional struct {
	Node

	If   Expression
	Then CodeBlock
}

type LetStatement struct {
	Node

	VariableName string
	Type         TypeDefinition
	Value        Expression
}

type ReturnStatement struct {
	Node

	Value Expression
}

type SwitchStatement struct {
	Node

	InputValue Expression
	Cases      []SwitchStatementCase
	Default    CodeBlock
}

type SwitchStatementCase struct {
	Node

	ComparisonValue Expression
	Then            CodeBlock
}

type ThrowStatement struct {
	Node

	Value Expression
}

type TryStatement struct {
	Node

	Body    CodeBlock
	Catches []TryStatementCatch
	Else    CodeBlock
	Finally CodeBlock
}

type TryStatementCatch struct {
	Node

	Exceptions []Expression
	Body       CodeBlock
}

type WhileLoopStatement struct {
	Node

	Label     string
	Condition Expression
	Body      CodeBlock
	Else      CodeBlock
}

type WithStatement struct {
	Node

	Value  Expression
	Target Expression
	Body   CodeBlock
}
