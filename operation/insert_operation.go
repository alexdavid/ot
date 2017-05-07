package operation

type InsertOperation struct {
	Content string
}

func (op1 InsertOperation) Equals(op2 Operation) bool {
	switch op2 := op2.(type) {
	case InsertOperation:
		return op1.Content == op2.Content
	}
	return false
}
