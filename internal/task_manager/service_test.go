package task_manager

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
)

// Función auxiliar para crear un puntero a string
func strPtr(s string) *string {
	return &s
}

// Mock para simular ejecución exitosa de una tarea
func mockSuccessProcessFn(ctx context.Context, input string) error {
	return nil
}

// Mock para simular ejecución fallida de una tarea
func mockFailProcessFn(ctx context.Context, input string) error {
	return errors.New("mock error")
}

func TestExecuteTasks_Success(t *testing.T) {
	// Crear manager con una función de procesamiento exitosa
	manager := NewManager(2, mockSuccessProcessFn)

	// Simular mensajes SQS
	messages := []types.Message{
		{MessageId: strPtr("msg-1"), Body: strPtr("Task 1")},
		{MessageId: strPtr("msg-2"), Body: strPtr("Task 2")},
	}

	// Ejecutar tareas
	results := manager.ExecuteTasks(context.TODO(), messages)

	// Validar que todas las tareas se ejecutaron sin error
	assert.Len(t, results, 2)
	for _, result := range results {
		assert.NoError(t, result.Err)
	}
}

func TestExecuteTasks_Failure(t *testing.T) {
	// Crear manager con una función de procesamiento que falla
	manager := NewManager(2, mockFailProcessFn)

	// Simular mensajes SQS
	messages := []types.Message{
		{MessageId: strPtr("msg-1"), Body: strPtr("Task 1")},
	}

	// Ejecutar tareas
	results := manager.ExecuteTasks(context.TODO(), messages)

	// Validar que hubo un error en la ejecución
	assert.Len(t, results, 1)
	assert.Error(t, results["msg-1"].Err)
	assert.EqualError(t, results["msg-1"].Err, "mock error")
}

func TestExecuteTasks_EmptyMessages(t *testing.T) {
	// Crear manager con una función de procesamiento (éxito o fallo no importa)
	manager := NewManager(2, mockSuccessProcessFn)

	// Ejecutar con una lista vacía de mensajes
	results := manager.ExecuteTasks(context.TODO(), []types.Message{})

	// Validar que no se procesaron tareas
	assert.Len(t, results, 0)
}
