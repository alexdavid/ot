package operation

type Operation interface {
	Equals(Operation) bool
}
