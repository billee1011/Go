package abs

/**
抽象的任务池
 */
type TaskPool interface{
	GetName() string
	Start()
	Execute(task Task)
	StopNow()
	Stop()
	OnTaskError(task Task,err error)
}
