package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestListLogs(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
		want    []string
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "logs",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/",
					httpmock.NewStringResponder(http.StatusOK, `[
						"first",
						"second"
					]`))
			},
			want: []string{"first", "second"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.ListLogs(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("ListLogs() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ListLogs() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestLogParserParse(t *testing.T) {
	for _, tc := range []struct {
		name  string
		loc   *time.Location
		input []string
		want  []LogEntry
	}{
		{
			name: "empty",
		},
		{
			name: "one line",
			input: []string{
				"[2023-02-28 00:28:37,604] [INFO] [paperless.consumer] Consuming xyz.pdf",
			},
			want: []LogEntry{
				{
					Time:    time.Date(2023, time.February, 28, 0, 28, 37, 604000000, time.UTC),
					Level:   "INFO",
					Module:  "paperless.consumer",
					Message: "Consuming xyz.pdf",
				},
			},
		},
		{
			name: "joined lines",
			input: []string{
				"[2020-01-01 01:02:03.123] [INFO] [foo] Command xyz failed:\t",
				"  Command not found\t",
				"[2020-01-01 03:04:05] [ERROR] [bar] Something bad happened",
			},
			want: []LogEntry{
				{
					Time:    time.Date(2020, time.January, 1, 1, 2, 3, 123000000, time.UTC),
					Level:   "INFO",
					Module:  "foo",
					Message: "Command xyz failed:\n  Command not found",
				},
				{
					Time:    time.Date(2020, time.January, 1, 3, 4, 5, 0, time.UTC),
					Level:   "ERROR",
					Module:  "bar",
					Message: "Something bad happened",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.loc == nil {
				tc.loc = time.UTC
			}

			p := logParser{
				loc: tc.loc,
			}

			got := p.parse(tc.input)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("logParser.parse() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetLog(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		logName string
		wantErr error
		want    []LogEntry
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/mail/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
			logName: "mail",
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/system/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			logName: "system",
			wantErr: cmpopts.AnyError,
		},
		{
			name: "parsed",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/logs/entries/",
					httpmock.NewStringResponder(http.StatusOK, `[
						"ignored",
						"more ignored",
						"[2023-02-28 00:28:37,604] [INFO] [paperless.consumer] Consuming xyz.pdf",
						"[2023-02-28 01:00:12,931] [INFO] [paperless] Another message"
					]`))
			},
			logName: "entries",
			want: []LogEntry{
				{
					Time:    time.Date(2023, time.February, 28, 0, 28, 37, 604000000, time.UTC),
					Level:   "INFO",
					Module:  "paperless.consumer",
					Message: "Consuming xyz.pdf",
				},
				{
					Time:    time.Date(2023, time.February, 28, 1, 0, 12, 931000000, time.UTC),
					Level:   "INFO",
					Module:  "paperless",
					Message: "Another message",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				ServerLocation: time.UTC,
				transport:      transport,
			})

			got, _, err := c.GetLog(context.Background(), tc.logName)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetLog() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetLog() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
