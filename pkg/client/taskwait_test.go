package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func makeShortBackOff() backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = time.Millisecond / 10
	b.MaxInterval = time.Millisecond
	b.MaxElapsedTime = 10 * time.Millisecond
	return b
}

func TestTaskWaiter(t *testing.T) {
	for _, tc := range []struct {
		name    string
		b       backoff.BackOff
		get     func(*testing.T) (*Task, error)
		cond    WaitForTaskConditionFunc
		wantErr error
		want    *Task
	}{
		{
			name: "immediate success",
			get: func(t *testing.T) (*Task, error) {
				return &Task{
					ID:     13900,
					Status: TaskSuccess,
				}, nil
			},
			cond: DefaultWaitForTaskCondition,
			want: &Task{
				ID:     13900,
				Status: TaskSuccess,
			},
		},
		{
			name: "success after a few requests",
			get: func() func(t *testing.T) (*Task, error) {
				var counter int

				return func(_ *testing.T) (*Task, error) {
					defer func() { counter++ }()

					if counter == 0 {
						return nil, &RequestError{
							StatusCode: http.StatusInternalServerError,
						}
					}

					t := &Task{
						ID:     94,
						Status: TaskPending,
					}

					if counter > 3 {
						t.Status = TaskSuccess
					}

					return t, nil
				}
			}(),
			cond: DefaultWaitForTaskCondition,
			want: &Task{
				ID:     94,
				Status: TaskSuccess,
			},
		},
		{
			name: "request always fails",
			get: func(t *testing.T) (*Task, error) {
				return nil, &RequestError{
					StatusCode: http.StatusTeapot,
				}
			},
			cond: DefaultWaitForTaskCondition,
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
			},
		},
		{
			name: "task failed",
			get: func(t *testing.T) (*Task, error) {
				return &Task{
					ID:     281,
					TaskID: "abc",
					Status: TaskFailure,
					Result: String("foobar"),
				}, nil
			},
			cond: DefaultWaitForTaskCondition,
			wantErr: &TaskError{
				TaskID: "abc",
				Status: TaskFailure,
			},
			want: &Task{
				ID:     281,
				TaskID: "abc",
				Status: TaskFailure,
				Result: String("foobar"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			t.Cleanup(cancel)

			w := taskWaiter{
				logger: &discardLogger{},
				b:      tc.b,
				get: func(_ context.Context) (*Task, error) {
					return tc.get(t)
				},
				cond: tc.cond,
			}

			if w.b == nil {
				w.b = makeShortBackOff()
			}

			got, err := w.wait(ctx)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("taskWaiter error diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("taskWaiter result diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWaitForTask(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		taskID  string
		opts    WaitForTaskOptions
		wantErr error
		want    *Task
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tasks/",
					"task_id=xyz404",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
			taskID: "xyz404",
			wantErr: &RequestError{
				StatusCode: http.StatusNotFound,
				Message:    `task "xyz404" not found`,
			},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tasks/",
					"task_id=badjson",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			taskID:  "badjson",
			wantErr: cmpopts.AnyError,
		},
		{
			name: "task success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tasks/",
					"task_id=successtask",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 27502,
							"task_id": "2e84c704-8762-4499-b144-29673844a2c1",
							"task_file_name": "success.pdf",
							"status": "SUCCESS",
							"result": "Success. New document id 26150 created"
						}
					]`))
			},
			taskID: "successtask",
			want: &Task{
				ID:           27502,
				TaskID:       "2e84c704-8762-4499-b144-29673844a2c1",
				TaskFileName: String("success.pdf"),
				Status:       TaskSuccess,
				Result:       String("Success. New document id 26150 created"),
			},
		},
		{
			name: "task failure",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tasks/",
					"task_id=failtask",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 4766,
							"task_id": "f59d61da-bd2f-46b9-a4e4-7ac5cfd5e606",
							"task_file_name": "fail.pdf",
							"status": "FAILURE",
							"result": "Something went wrong"
						}
					]`))
			},
			taskID: "failtask",
			wantErr: &TaskError{
				TaskID: "f59d61da-bd2f-46b9-a4e4-7ac5cfd5e606",
				Status: TaskFailure,
			},
			want: &Task{
				ID:           4766,
				TaskID:       "f59d61da-bd2f-46b9-a4e4-7ac5cfd5e606",
				TaskFileName: String("fail.pdf"),
				Status:       TaskFailure,
				Result:       String("Something went wrong"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			tc.opts.MaxElapsedTime = time.Second

			got, err := c.WaitForTask(context.Background(), tc.taskID, tc.opts)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("WaitForTask() error diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("WaitForTask() diff (-want +got):\n%s", diff)
			}
		})
	}
}
