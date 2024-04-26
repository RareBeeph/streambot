package tasks

type Task struct {
	Spec    string
	Handler func()
}

var All = []Task{}
