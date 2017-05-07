package operation

type DeleteOperation struct {
	Content string
}

func (op1 DeleteOperation) Equals(op2 Operation) bool {
	switch op2 := op2.(type) {
	case DeleteOperation:
		return op1.Content == op2.Content
	}
	return false
}
