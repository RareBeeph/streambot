package tasks

type Task struct {
	CronSpec string
	CronFunc func()
}

var All = []Task{}
