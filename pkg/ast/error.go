package ast

type ASTError interface {
	ASTNode

	Error() string
	InnerErrors() []ASTError
}

func NewError(astNode ASTNode, message string, innerErrors []ASTError) ASTError {
	return &astError{
		ASTNode:     astNode,
		message:     message,
		innerErrors: innerErrors,
	}
}

type astError struct {
	ASTNode
	message     string
	innerErrors []ASTError
}

func (err *astError) Error() string {
	return err.message
}

func (err *astError) InnerErrors() []ASTError {
	return err.innerErrors
}
