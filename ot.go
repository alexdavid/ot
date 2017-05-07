package ot

import . "github.com/alexdavid/ot/operation"

type Transform struct {
	operations   []Operation
	BaseLength   int
	TargetLength int
}

////////////////////////////////////////////////////
// Static methods
func Equals(t1 *Transform, t2 *Transform) bool {
	if t1.BaseLength != t2.BaseLength {
		return false
	}
	if t1.TargetLength != t2.TargetLength {
		return false
	}
	for i, op := range t1.operations {
		if !op.Equals(t2.operations[i]) {
			return false
		}
	}
	return true
}

////////////////////////////////////////////////////
// Accessor methods
func (t *Transform) Operations() []Operation {
	if t.operations != nil {
		return t.operations
	}
	return []Operation{}
}

func (t *Transform) Last() Operation {
	index := len(t.operations) - 1
	if index < 0 {
		return nil
	}
	return t.operations[index]
}

func (t *Transform) Penultimate() Operation {
	index := len(t.operations) - 2
	if index < 0 {
		return nil
	}
	return t.operations[index]
}

////////////////////////////////////////////////////
// Internal helpers
func (t *Transform) appendOp(op Operation) {
	t.operations = append(t.operations, op)
}

func (t *Transform) replaceLastOp(op Operation) {
	if len(t.operations) > 0 {
		t.operations[len(t.operations)-1] = op
	}
}

////////////////////////////////////////////////////
// Operation methods
func (t *Transform) Retain(length int) *Transform {
	if length < 1 {
		return t
	}

	t.BaseLength += length
	t.TargetLength += length

	switch lastOp := t.Last().(type) {
	case RetainOperation:
		t.replaceLastOp(RetainOperation{Length: lastOp.Length + length})
	default:
		t.appendOp(RetainOperation{Length: length})
	}

	return t
}

func (t *Transform) Insert(content string) *Transform {
	if content == "" {
		return t
	}
	t.TargetLength += len(content)

	switch lastOp := t.Last().(type) {
	case InsertOperation:
		// If the last operation is an insert we append to the operation
		t.replaceLastOp(InsertOperation{Content: lastOp.Content + content})

	case DeleteOperation:
		// Insert("foo").Delete("bar") is the same as Delete("bar").Insert("foo")
		// To prevent these operations from differing we enforce that inserts
		// always come before deletes
		switch penultimateOp := t.Penultimate().(type) {
		case InsertOperation:
			t.operations[len(t.operations)-2] = InsertOperation{Content: penultimateOp.Content + content}
		default:
			t.replaceLastOp(InsertOperation{Content: content})
			t.appendOp(lastOp)
		}

	default:
		t.appendOp(InsertOperation{Content: content})
	}

	return t
}

func (t *Transform) Delete(content string) *Transform {
	if content == "" {
		return t
	}
	t.BaseLength += len(content)

	switch lastOp := t.Last().(type) {
	case DeleteOperation:
		t.replaceLastOp(DeleteOperation{Content: lastOp.Content + content})
	default:
		t.appendOp(DeleteOperation{Content: content})
	}

	return t
}
