package ast

type Statement interface{}

type BreakStatement struct {
	ASTNode

	Label string
}

type ContinueStatement struct {
	ASTNode

	Label string
}

type ExpressionStatement struct {
	ASTNode

	Expression Expression
}

type ForLoopStatement struct {
	ASTNode

	Label    string
	Target   VariableExpression
	Iterable Expression
	Body     CodeBlock
	Else     CodeBlock
}

type IfStatement struct {
	ASTNode

	Conditionals []IfStatementConditional
	Else         CodeBlock
}

type IfStatementConditional struct {
	ASTNode

	If   Expression
	Then CodeBlock
}

type LetStatement struct {
	ASTNode

	VariableName string
	Type         TypeDefinition
	Value        Expression
}

type ReturnStatement struct {
	ASTNode

	Value Expression
}

type SwitchStatement struct {
	ASTNode

	InputValue Expression
	Cases      []SwitchStatementCase
	Default    CodeBlock
}

type SwitchStatementCase struct {
	ASTNode

	ComparisonValue Expression
	Then            CodeBlock
}

type ThrowStatement struct {
	ASTNode

	Value Expression
}

type TryStatement struct {
	ASTNode

	Body    CodeBlock
	Catches []TryStatementCatch
	Else    CodeBlock
	Finally CodeBlock
}

type TryStatementCatch struct {
	ASTNode

	Exceptions []Expression
	Body       CodeBlock
}

type WhileLoopStatement struct {
	ASTNode

	Label     string
	Condition Expression
	Body      CodeBlock
	Else      CodeBlock
}

type WithStatement struct {
	ASTNode

	Value  Expression
	Target Expression
	Body   CodeBlock
}
