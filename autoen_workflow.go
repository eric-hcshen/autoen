package autoen

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

/**
 * This sample workflow executes multiple branches in parallel using workflow.Go() method.
 */
type (
	// TripEvent passed in as signal to TripWorkflow
	TripEvent struct {
		ID    string
		Total int
	}
)

const (
	// TripSignalName is the signal name for trip completion event
	TripSignalName = "trip_event"

	// QueryName is the query type name
	QueryName = "counter"
)

// SampleParallelWorkflow workflow definition
func AutoEnWorkflow(ctx workflow.Context) ([]string, error) {
	logger := workflow.GetLogger(ctx)
	defer logger.Info("Workflow completed.")
	var results []string
	var result string
	// Register query handler to return trip count
	err := workflow.SetQueryHandler(ctx, QueryName, func(input []byte) ([]string, error) {
		return results, nil
	})
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10000 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err = workflow.ExecuteActivity(ctx, SampleActivity, "Main Start Point").Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	results = append(results, "start..")
	results = append(results, result)

	futureA := runWorkflowA(ctx)
	futureB := runWorkflowB(ctx)

	/*
			future1, settable1 := workflow.NewFuture(ctx)
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer logger.Info("A goroutine completed.")

				var results []string
				var result string
				err := workflow.ExecuteActivity(ctx, SampleActivity, "branch1.1").Get(ctx, &result)
				if err != nil {
					settable1.SetError(err)
					return
				}
				results = append(results, result)
				err = workflow.ExecuteActivity(ctx, SampleActivity, "branch1.2").Get(ctx, &result)
				if err != nil {
					settable1.SetError(err)
					return
				}
				results = append(results, result)
				settable1.SetValue(results)
			})
		future2, settable2 := workflow.NewFuture(ctx)
		workflow.Go(ctx, func(ctx workflow.Context) {
			defer logger.Info("Second goroutine completed.")

			var result string
			err := workflow.ExecuteActivity(ctx, SampleActivity, "branch2").Get(ctx, &result)
			settable2.Set(result, err)
		})*/

	// Future.Get returns error from Settable.SetError
	// Note that the first goroutine puts a slice into the settable while the second a string value
	var resultsA []string
	var resultsB []string
	err = futureA.Get(ctx, &resultsA)
	if err != nil {
		return nil, err
	}
	results = append(results, resultsA[0], resultsA[1])
	err = futureB.Get(ctx, &resultsB)
	if err != nil {
		return nil, err
	}
	results = append(results, resultsB[0])
	results = append(results, "end...")

	return results, nil
}

func runWorkflowA(ctx workflow.Context) workflow.Future {
	logger := workflow.GetLogger(ctx)
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer logger.Info("A goroutine completed.")

		var results []string
		var result string
		err := workflow.ExecuteActivity(ctx, SampleActivity, "branch_A.1").Get(ctx, &result)
		if err != nil {
			settable.SetError(err)
			return
		}
		results = append(results, result)
		err = workflow.ExecuteActivity(ctx, SampleActivity, "branch_A.2").Get(ctx, &result)
		if err != nil {
			settable.SetError(err)
			return
		}
		results = append(results, result)
		settable.SetValue(results)
	})
	return future
}

func runWorkflowB(ctx workflow.Context) workflow.Future {
	logger := workflow.GetLogger(ctx)
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer logger.Info("B goroutine completed.")

		var results []string
		var result string

		ch := workflow.GetSignalChannel(ctx, TripSignalName)
		ch.Receive(ctx, &results)
		err := workflow.ExecuteActivity(ctx, SampleActivity, "branch_B.1").Get(ctx, &result)
		if err != nil {
			settable.SetError(err)
			return
		}
		results = append(results, result)
		settable.SetValue(results)
	})
	return future
}

func SampleActivity(input string) (string, error) {
	name := "sampleActivity"
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + input, nil
}
