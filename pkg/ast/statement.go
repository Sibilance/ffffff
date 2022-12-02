package ast

type Statement interface{}

type BreakStatement struct {
	AstNode
}

type ContinueStatement struct {
	AstNode
}

type ExpressionStatement struct {
	AstNode

	Expression Expression
}

type ForLoopStatement struct {
	AstNode

	Target   LValue
	Iterable Expression
	Body     CodeBlock
}

type IfStatement struct {
	AstNode

	Conditionals []IfStatementConditional
	Else         CodeBlock
}

type IfStatementConditional struct {
	AstNode

	If   Expression
	Then CodeBlock
}

type LetStatement struct {
	AstNode

	VariableName string
	Type         TypeDefinition
	Value        Expression
}

type ReturnStatement struct {
	AstNode

	Value Expression
}

type SwitchStatement struct {
	AstNode

	InputValue Expression
	Cases      []SwitchStatementCase
	Default    CodeBlock
}

type SwitchStatementCase struct {
	AstNode

	ComparisonValue Expression
	Then            CodeBlock
}

type ThrowStatement struct {
	AstNode

	Value Expression
}

type TryStatement struct {
	AstNode

	Body    CodeBlock
	Catches []TryStatementCatch
	Else    CodeBlock
	Finally CodeBlock
}

type TryStatementCatch struct {
	AstNode

	Exceptions []Expression
	Body       CodeBlock
}

type WhileLoopStatement struct {
	AstNode

	Condition Expression
	Body      CodeBlock
	Else      CodeBlock
}

type WithStatement struct {
	AstNode

	Value  Expression
	Target LValue
	Body   CodeBlock
}
