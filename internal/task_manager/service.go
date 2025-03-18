package task_manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/task_executor"
)

type Service interface {
	ExecuteTasks(ctx context.Context, messages []types.Message) map[string]task_executor.Result
}

type Manager struct {
	numWorkers int
	processFn  func(ctx context.Context, input string) error
}

func NewManager(numWorkers int, processFn func(ctx context.Context, input string) error) Manager {
	return Manager{
		numWorkers: numWorkers,
		processFn:  processFn,
	}
}

func (m Manager) ExecuteTasks(ctx context.Context, messages []types.Message) map[string]task_executor.Result {
	tasks := make(map[string]task_executor.Tasker)

	for _, msg := range messages {
		message := msg
		tasks[*message.MessageId] = task_executor.Task[string, error]{
			Func: func(ctx context.Context, input string) (error, error) {
				err := m.processFn(ctx, input)
				return err, err
			},
			Args: *message.Body,
		}
	}

	return task_executor.WorkerPool(ctx, tasks, m.numWorkers)
}
