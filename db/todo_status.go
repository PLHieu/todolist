package db

// TodoStatus contains the different types of Todo status.
type TodoStatus int

const (
	Done TodoStatus = iota
	Undone
)

func (o TodoStatus) String() string {
	return [...]string{"DONE", "UNDONE"}[o]
}
