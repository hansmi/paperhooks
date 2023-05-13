package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:generate stringer -linecomment -type=TaskStatus -trimprefix=Task -output=task_string.go
type TaskStatus int

// Celery task states
// (https://docs.celeryq.dev/en/latest/userguide/tasks.html#built-in-states).
const (
	TaskStatusUnspecified TaskStatus = iota

	// Task is waiting for execution.
	TaskPending

	// Task has been started.
	TaskStarted

	// Task has been successfully executed.
	TaskSuccess

	// Task execution resulted in failure.
	TaskFailure

	// Task is being retried.
	TaskRetry

	// Task has been revoked.
	TaskRevoked
)

var taskStatusText = map[TaskStatus]string{
	TaskPending: "PENDING",
	TaskStarted: "STARTED",
	TaskSuccess: "SUCCESS",
	TaskFailure: "FAILURE",
	TaskRetry:   "RETRY",
	TaskRevoked: "REVOKED",
}

var _ json.Marshaler = (*TaskStatus)(nil)
var _ json.Unmarshaler = (*TaskStatus)(nil)

func (s TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(taskStatusText[s])
}

func (s *TaskStatus) UnmarshalJSON(data []byte) error {
	var str *string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == nil || *str == "" {
		*s = TaskStatusUnspecified
		return nil
	}

	for key, value := range taskStatusText {
		if strings.EqualFold(*str, value) {
			*s = key
			return nil
		}
	}

	return fmt.Errorf("unrecognized task status %q", *str)
}

type Task struct {
	ID           int64      `json:"id"`
	TaskID       string     `json:"task_id"`
	TaskFileName *string    `json:"task_file_name"`
	Created      *time.Time `json:"date_created"`
	Done         *time.Time `json:"date_done"`
	Type         string     `json:"type"`
	Status       TaskStatus `json:"status"`
	Result       *string    `json:"result"`
	Acknowledged bool       `json:"acknowledged"`
}

func (c *Client) ListTasks(ctx context.Context) ([]Task, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult([]Task(nil)).
		Get("api/tasks/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return *resp.Result().(*[]Task), wrapResponse(resp), nil
}

func (c *Client) GetTask(ctx context.Context, taskID string) (*Task, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult([]*Task(nil)).
		SetQueryParam("task_id", taskID).
		Get("api/tasks/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	switch tasks := *resp.Result().(*[]*Task); len(tasks) {
	case 0:
		return nil, wrapResponse(resp), &RequestError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("task %q not found", taskID),
		}

	case 1:
		return tasks[0], wrapResponse(resp), nil

	default:
		return nil, wrapResponse(resp), &RequestError{
			StatusCode: http.StatusMultipleChoices,
			Message:    fmt.Sprintf("received %d tasks for ID %q", len(tasks), taskID),
		}
	}
}
