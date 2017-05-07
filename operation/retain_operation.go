package operation

type RetainOperation struct {
	Length int
}

func (op1 RetainOperation) Equals(op2 Operation) bool {
	switch op2 := op2.(type) {
	case RetainOperation:
		return op1.Length == op2.Length
	}
	return false
}
