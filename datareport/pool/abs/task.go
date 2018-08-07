package abs

/**
表示一个任务的接口
 */

type Task interface {
	DoTask() error
}