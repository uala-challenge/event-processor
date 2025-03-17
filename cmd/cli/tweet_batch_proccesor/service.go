package tweet_batch_proccesor

import (
	"context"
	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/task_manager"
	"log"
	"time"
)

type Runner struct {
	processor   consume_tweet_event_sqs.Service
	taskManager task_manager.Manager
	queueURL    string
	batchSize   int
	numWorkers  int
	retries     int
}

func NewRunner(processor consume_tweet_event_sqs.Service, taskManager task_manager.Manager, queueURL string) *Runner {
	return &Runner{
		processor:   processor,
		taskManager: taskManager,
		queueURL:    queueURL,
		batchSize:   10,
		numWorkers:  5,
		retries:     3,
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
