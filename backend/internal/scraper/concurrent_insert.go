package scraper

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type ConcurrentInserter[T any] struct {
	expectedTotalRecords int
	insertFn             func([]T) error
	ResC                 chan T
	breathingTime        time.Duration
	batchSize            int
}

func NewConcurrentInserter[T any](
	expectedTotalRecords int,
	insertFn func([]T) error,
	breathingTime time.Duration,
	batchSize int,
) *ConcurrentInserter[T] {
	return &ConcurrentInserter[T]{
		expectedTotalRecords: expectedTotalRecords,
		insertFn:             insertFn,
		ResC:                 make(chan T, 10_000),
		breathingTime:        breathingTime,
		batchSize:            min(batchSize, expectedTotalRecords),
	}
}

func (ci *ConcurrentInserter[T]) logProgress(buffer []T, totalProcessedRecords int) {
	fmt.Println(
		"Type",
		reflect.TypeOf(ci),
		"buffer :",
		len(buffer),
		"processed :",
		totalProcessedRecords,
		"expected :",
		ci.expectedTotalRecords,
		"batch",
		ci.batchSize,
		"rows left :",
		ci.expectedTotalRecords-totalProcessedRecords,
	)
}

func (ci *ConcurrentInserter[T]) Serve(ctx context.Context) error {
	buffer := make([]T, 0, 10_000)

	totalProcessedRecords := 0

	ticker := time.NewTicker(ci.breathingTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ci.logProgress(buffer, totalProcessedRecords)

			// Process the buffer if it has records or if we are done processing
			if len(buffer) > 0 && (len(buffer) >= ci.batchSize || totalProcessedRecords < ci.expectedTotalRecords) {
				if err := ci.insertFn(buffer); err != nil {
					return err
				}
				totalProcessedRecords += len(buffer)
				ci.logProgress(buffer, totalProcessedRecords)
				buffer = nil
			}

			// Check if we have processed all expected records
			if totalProcessedRecords >= ci.expectedTotalRecords && len(buffer) == 0 {
				return nil
			}

		case res := <-ci.ResC:
			buffer = append(buffer, res)
		}
	}
}
