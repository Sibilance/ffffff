package languagemodels

type Statement interface{}

type BreakStatement struct {
	Name string
}

type ContinueStatement struct {
	Name string
}

type ExpressionStatement struct {
	Name string

	Expression Expression
}

type ForLoopStatement struct {
	Name string

	Target   LValue
	Iterable Expression
	Body     CodeBlock
}

type IfStatement struct {
	Name string

	Conditionals []IfStatementConditional
	Else         CodeBlock
}

type IfStatementConditional struct {
	If   Expression
	Then CodeBlock
}

type LetStatement struct {
	Name string

	VariableName string
	Type         TypeDefinition
	Value        Expression
}

type ReturnStatement struct {
	Name string

	Value Expression
}

type SwitchStatement struct {
	Name string

	InputValue Expression
	Cases      []SwitchStatementCase
	Default    CodeBlock
}

type SwitchStatementCase struct {
	Name string

	ComparisonValue Expression
	Then            CodeBlock
}

type ThrowStatement struct {
	Name string

	Value Expression
}

type TryStatement struct {
	Name string

	Body    CodeBlock
	Catches []TryStatementCatch
	Else    CodeBlock
	Finally CodeBlock
}

type TryStatementCatch struct {
	Name string

	Exceptions []Expression
	Body       CodeBlock
}

type WhileLoopStatement struct {
	Name string

	Condition Expression
	Body      CodeBlock
	Else      CodeBlock
}

type WithStatement struct {
	Name string

	Value  Expression
	Target LValue
	Body   CodeBlock
}
