package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type WaitForTaskConditionFunc func(*Task) error

type WaitForTaskOptions struct {
	// Condition returns nil when the task is considered finished.
	// [DefaultWaitForTaskCondition] is the default implementation.
	Condition WaitForTaskConditionFunc

	// Maximum amount of time to wait. Defaults to one hour.
	MaxElapsedTime time.Duration
}

func (o WaitForTaskOptions) makeBackOff() backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 5 * time.Second
	b.MaxInterval = time.Minute

	if b.MaxElapsedTime = o.MaxElapsedTime; b.MaxElapsedTime == 0 {
		b.MaxElapsedTime = time.Hour
	}

	return b
}

// DefaultWaitForTaskCondition waits for a terminal status (success, failure
// or revoked).
func DefaultWaitForTaskCondition(task *Task) error {
	if task.Status.Terminal() {
		return nil
	}

	return fmt.Errorf("task %q has non-terminal status %q", task.TaskID, task.Status)
}

type taskWaiter struct {
	logger Logger
	b      backoff.BackOff
	get    func(context.Context) (*Task, error)
	cond   WaitForTaskConditionFunc
}

func (w taskWaiter) wait(ctx context.Context) (*Task, error) {
	task, err := backoff.RetryNotifyWithData(func() (*Task, error) {
		task, err := w.get(ctx)

		if err != nil {
			var reqErr *RequestError

			if errors.As(err, &reqErr) && (reqErr.StatusCode/100) == 5 {
				return nil, err
			}

			return nil, backoff.Permanent(err)
		}

		return task, w.cond(task)
	}, backoff.WithContext(w.b, ctx), func(err error, delay time.Duration) {
		w.logger.Debugf("Condition not met, retry in %s: %v", delay.String(), err)
	})

	if err == nil {
		err = task.statusError()
	}

	return task, err
}

// WaitForTask polls the status of a task until it reaches a terminal status
// (success, failure or revoked). Task failures are reported as an error of
// type [TaskError].
func (c *Client) WaitForTask(ctx context.Context, taskID string, opts WaitForTaskOptions) (*Task, error) {
	w := taskWaiter{
		logger: &prefixLogger{
			wrapped: c.logger,
			prefix:  fmt.Sprintf("task %s: ", taskID),
		},
		b: opts.makeBackOff(),
		get: func(ctx context.Context) (*Task, error) {
			task, _, err := c.GetTask(ctx, taskID)
			return task, err
		},
		cond: opts.Condition,
	}

	if w.cond == nil {
		w.cond = DefaultWaitForTaskCondition
	}

	return w.wait(ctx)
}
