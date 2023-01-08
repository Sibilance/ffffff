package ast

type Error interface {
	Node

	Error() string
	InnerErrors() []Error
}

func NewError(astNode Node, message string, innerErrors []Error) Error {
	return &astError{
		Node:        astNode,
		message:     message,
		innerErrors: innerErrors,
	}
}

type astError struct {
	Node
	message     string
	innerErrors []Error
}

func (err *astError) Error() string {
	return err.message
}

func (err *astError) InnerErrors() []Error {
	return err.innerErrors
}
