package ast

type Statement[N Node] interface {
	Node
}

type BreakStatement[N Node] struct {
	Node N

	Label string
}

type ContinueStatement[N Node] struct {
	Node N

	Label string
}

type ExpressionStatement[N Node] struct {
	Node N

	Expression Expression[N]
}

type ForLoopStatement[N Node] struct {
	Node N

	Label    string
	Target   VariableExpression[N]
	Iterable Expression[N]
	Body     CodeBlock[N]
	Else     CodeBlock[N]
}

type IfStatement[N Node] struct {
	Node N

	Conditionals []IfStatementConditional[N]
	Else         CodeBlock[N]
}

type IfStatementConditional[N Node] struct {
	Node N

	If   Expression[N]
	Then CodeBlock[N]
}

type LetStatement[N Node] struct {
	Node N

	VariableName string
	Type         TypeDefinition[N]
	Value        Expression[N]
}

type ReturnStatement[N Node] struct {
	Node N

	Value Expression[N]
}

type SwitchStatement[N Node] struct {
	Node N

	InputValue Expression[N]
	Cases      []SwitchStatementCase[N]
	Default    CodeBlock[N]
}

type SwitchStatementCase[N Node] struct {
	Node N

	ComparisonValue Expression[N]
	Then            CodeBlock[N]
}

type ThrowStatement[N Node] struct {
	Node N

	Value Expression[N]
}

type TryStatement[N Node] struct {
	Node N

	Body    CodeBlock[N]
	Catches []TryStatementCatch[N]
	Else    CodeBlock[N]
	Finally CodeBlock[N]
}

type TryStatementCatch[N Node] struct {
	Node N

	Exceptions []Expression[N]
	Body       CodeBlock[N]
}

type WhileLoopStatement[N Node] struct {
	Node N

	Label     string
	Condition Expression[N]
	Body      CodeBlock[N]
	Else      CodeBlock[N]
}

type WithStatement[N Node] struct {
	Node N

	Value  Expression[N]
	Target Expression[N]
	Body   CodeBlock[N]
}
