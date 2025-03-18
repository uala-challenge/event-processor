package tweet_batch_proccesor

import (
	"context"
	"log"
	"time"

	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/task_manager"
)

type Runner struct {
	processor   consume_tweet_event_sqs.Service
	taskManager task_manager.Manager
	queueURL    string
	batchSize   int
	numWorkers  int
	retries     int
}

func NewRunner(d Dependencies) *Runner {
	return &Runner{
		processor:   d.Processor,
		taskManager: d.TaskManager,
		queueURL:    d.Config.Endpoint,
		batchSize:   10,
		numWorkers:  2,
		retries:     1,
	}
}

func (r *Runner) Run(ctx context.Context) {
	log.Println("🚀 Iniciando procesamiento batch...")

	for {
		messages, err := r.processor.ReceiveMessages(ctx, r.queueURL, r.batchSize, r.retries)
		if err != nil {
			log.Printf("Error recibiendo mensajes: %v", err)
			continue
		}

		if len(messages) == 0 {
			log.Println("No hay mensajes en la cola.")
			time.Sleep(5 * time.Second)
			continue
		}

		results := r.taskManager.ExecuteTasks(ctx, messages)

		for msgID, result := range results {
			if result.Err == nil {
				err := r.processor.DeleteMessage(ctx, r.queueURL, msgID)
				if err != nil {
					log.Printf("Error eliminando mensaje: %v", err)
				}
			} else {
				log.Printf("Error procesando mensaje %s: %v", msgID, result.Err)
			}
		}
	}
}
