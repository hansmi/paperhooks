package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestTaskStatus(t *testing.T) {
	type test struct {
		status TaskStatus
	}

	tests := []test{
		{},
	}

	for i := range taskStatusText {
		tests = append(tests, test{
			status: i,
		})
	}

	for _, tc := range tests {
		t.Run(tc.status.String(), func(t *testing.T) {
			buf, err := json.Marshal(tc.status)
			if err != nil {
				t.Fatalf("Marshal(%v) failed: %v", tc.status, err)
			}

			var got TaskStatus

			if err := json.Unmarshal(buf, &got); err != nil {
				t.Fatalf("Unmarshal(%q) failed: %v", buf, err)
			}

			if diff := cmp.Diff(tc.status, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("StatusText diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
		want    []Task
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "tasks",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 11,
							"task_id": "69b9bf86-97d5-4e84-bf56-fbfcf5797f3b",
							"task_file_name": "example.pdf",
							"date_created": "2023-02-28T00:40:40.991412+00:00",
							"date_done": "2023-02-28T00:40:43.836833-00:00",
							"type": "file",
							"status": "SUCCESS",
							"result": "Success. New document id 26150 created",
							"acknowledged": false,
							"related_document": "26150"
						}, {
							"id": 22,
							"task_id": "minimal"
						}
					]`))
			},
			want: []Task{
				{
					ID:           11,
					TaskID:       "69b9bf86-97d5-4e84-bf56-fbfcf5797f3b",
					TaskFileName: String("example.pdf"),
					Created:      Time(time.Date(2023, time.February, 28, 0, 40, 40, 991412000, time.UTC)),
					Done:         Time(time.Date(2023, time.February, 28, 0, 40, 43, 836833000, time.UTC)),
					Type:         "file",
					Status:       TaskSuccess,
					Result:       String("Success. New document id 26150 created"),
					Acknowledged: false,
				},
				{
					ID:     22,
					TaskID: "minimal",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.ListTasks(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("ListTasks() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ListTasks() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		id      string
		wantErr error
		want    Task
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
			wantErr: &RequestError{
				StatusCode: http.StatusNotFound,
				Message:    `task "" not found`,
			},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "multiple",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tasks/",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 27781,
							"task_id": "c329f0f3-4d11-4be2-b6d0-16fd3c551215"
						}, {
							"id": 17165,
							"task_id": "2b83ba79-1fe5-4cf7-9b1e-69fd52494f86"
						}
					]`))
			},
			id: "req-id",
			wantErr: &RequestError{
				StatusCode: http.StatusMultipleChoices,
				Message:    `received 2 tasks for ID "req-id"`,
			},
		},
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tasks/",
					"task_id=7b0e2c15e928",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 7959,
							"task_id": "608a5b03-dd7a-4888-9cbc-14000d4655ba",
							"task_file_name": "task.pdf"
						}
					]`))
			},
			id: "7b0e2c15e928",
			want: Task{
				ID:           7959,
				TaskID:       "608a5b03-dd7a-4888-9cbc-14000d4655ba",
				TaskFileName: String("task.pdf"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.GetTask(context.Background(), tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetTask() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, *got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetTask() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
