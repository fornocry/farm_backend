package constant

type Task string

const (
	FRIENDS   Task = "FRIENDS"
	SUBSCRIBE Task = "SUBSCRIBE"
	INVENTORY Task = "INVENTORY"
)

type TaskCompleteStatus string

const (
	TASK_COMPLETE_NULL     TaskCompleteStatus = "TASK_COMPLETE_NULL"
	TASK_COMPLETE_DONE     TaskCompleteStatus = "TASK_COMPLETE_DONE"
	TASK_COMPLETE_FINISHED TaskCompleteStatus = "TASK_COMPLETE_FINISHED"
)
