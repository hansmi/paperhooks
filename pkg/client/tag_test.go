package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestListTags(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		opts    ListTagsOptions
		wantErr error
		want    []Tag
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "entries",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"results": [
							{ "id": 100, "name": "first" },
							{ "id": 200, "name": "second" }
						]
					}`))
			},
			want: []Tag{
				{ID: 100, Name: "first"},
				{ID: 200, Name: "second"},
			},
		},
		{
			name: "options",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"ordering=name&name__istartswith=hello&page=1&page_size=25",
					httpmock.NewStringResponder(http.StatusOK, `{
						"results": [
							{ "id": 400, "name": "four" },
							{ "id": 500, "name": "five" }
						]
					}`))
			},
			opts: ListTagsOptions{
				Ordering: OrderingSpec{
					Field: "name",
				},
				Name: CharFilterSpec{
					StartsWithIgnoringCase: String("hello"),
				},
			},
			want: []Tag{
				{ID: 400, Name: "four"},
				{ID: 500, Name: "five"},
			},
		},
		{
			name: "third page",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=3&page_size=25",
					httpmock.NewStringResponder(http.StatusOK, `{
						"results": [
							{ "id": 300, "name": "third" }
						]
					}`))
			},
			opts: ListTagsOptions{
				ListOptions: ListOptions{
					Page: &PageToken{number: 3},
				},
			},
			want: []Tag{
				{ID: 300, Name: "third"},
			},
		},
		{
			name: "first page not found",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=1&page_size=25",
					httpmock.NewStringResponder(http.StatusNotFound, `{}`))
			},
			wantErr: &RequestError{
				StatusCode: http.StatusNotFound,
				Message:    "{}",
			},
		},
		{
			name: "third page not found",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=3&page_size=25",
					httpmock.NewStringResponder(http.StatusNotFound, `{}`))
			},
			opts: ListTagsOptions{
				ListOptions: ListOptions{
					Page: &PageToken{number: 3},
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

			got, _, err := c.ListTags(context.Background(), tc.opts)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("ListTags() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ListTags() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestListAllTags(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		opts    ListTagsOptions
		wantErr error
		want    []Tag
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
		},
		{
			name: "three pages",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=1&page_size=25",
					httpmock.NewStringResponder(http.StatusOK, `{
						"next": "?page=2",
						"results": [
							{ "id": 10, "name": "first" },
							{ "id": 20, "name": "second" }
						]
					}`))
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=2&page_size=25",
					httpmock.NewStringResponder(http.StatusOK, `{
						"next": "?page=3"
					}`))
				transport.RegisterResponderWithQuery(http.MethodGet, "/api/tags/",
					"page=3&page_size=25",
					httpmock.NewStringResponder(http.StatusOK, `{
						"results": [
							{ "id": 90, "name": "last" }
						]
					}`))
			},
			want: []Tag{
				{ID: 10, Name: "first"},
				{ID: 20, Name: "second"},
				{ID: 90, Name: "last"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			var got []Tag

			err := c.ListAllTags(context.Background(), tc.opts, func(_ context.Context, v Tag) error {
				got = append(got, v)
				return nil
			})

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("ListAllTags() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ListAllTags() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestListAllTagsHandlerCancelsContext(t *testing.T) {
	var nextID atomic.Int64

	transport := newMockTransport(t)
	transport.RegisterResponder(http.MethodGet, "/api/tags/",
		httpmock.Responder(func(req *http.Request) (*http.Response, error) {
			var pageNumber int

			if str := req.FormValue("page"); str != "" {
				if value, err := strconv.Atoi(str); err != nil {
					t.Error(err)
					return nil, err
				} else {
					pageNumber = value
				}
			}

			result := &listResult[Tag]{}
			result.Next = fmt.Sprintf("?page=%d", pageNumber+1)

			for idx := 0; idx < 5; idx++ {
				result.Items = append(result.Items, Tag{
					ID: nextID.Add(1),
				})
			}

			return httpmock.NewJsonResponse(http.StatusOK, result)
		}))

	c := New(Options{
		transport: transport,
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var count atomic.Int64

	err := c.ListAllTags(ctx, ListTagsOptions{}, func(_ context.Context, v Tag) error {
		if count.Add(1) > 20 {
			cancel()
		}

		return nil
	})

	wantErr := context.Canceled

	if diff := cmp.Diff(wantErr, err, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("ListAllTags() error diff (-want +got):\n%s", diff)
	}
}

func TestGetTag(t *testing.T) {
	for _, tc := range []struct {
		name      string
		setup     func(*testing.T, *httpmock.MockTransport)
		id        int64
		responder httpmock.Responder
		wantErr   error
		want      *Tag
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/0/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			want: &Tag{},
		},
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/tags/0/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 123,
						"name": "inbox",
						"color": "#ff00ff",
						"matching_algorithm": 2,
						"is_inbox_tag": true
					}`))
			},
			want: &Tag{
				ID:                123,
				Name:              "inbox",
				Color:             Color{R: 255, B: 255},
				MatchingAlgorithm: MatchAll,
				IsInboxTag:        true,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.GetTag(context.Background(), tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetTag() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetTag() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		input   *Tag
		wantErr error
		want    *Tag
	}{
		{
			name:  "empty",
			input: &Tag{},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/tags/",
					httpmock.NewStringResponder(http.StatusCreated, `{}`))
			},
			want: &Tag{},
		},
		{
			name: "success",
			input: &Tag{
				Name: "foo",
			},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodPost, "/api/tags/",
					httpmock.BodyContainsString(`"foo"`),
					httpmock.NewStringResponder(http.StatusCreated, `{
						"id": 999,
						"name": "created"
					}`))
			},
			want: &Tag{
				ID:   999,
				Name: "created",
			},
		},
		{
			name:  "unexpected HTTP 200",
			input: &Tag{},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/tags/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			wantErr: &RequestError{
				StatusCode: http.StatusOK,
				Message:    "unexpected status 200",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.CreateTag(context.Background(), tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("CreateTag() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("CreateTag() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestUpdateTag(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		id      int64
		input   *Tag
		wantErr error
		want    *Tag
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPut, "/api/tags/14830/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			id:    14830,
			input: &Tag{},
			want:  &Tag{},
		},
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodPut, "/api/tags/123/",
					httpmock.BodyContainsString(`"newname"`),
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 999,
						"name": "created"
					}`))
			},
			id: 123,
			input: &Tag{
				Name:       "newname",
				IsInboxTag: true,
			},
			want: &Tag{
				ID:   999,
				Name: "created",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.UpdateTag(context.Background(), tc.id, tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("UpdateTag() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("UpdateTag() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestDeleteTag(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		id      int64
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodDelete, "/api/tags/7816/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			id: 7816,
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodDelete, "/api/tags/12234/",
					httpmock.NewStringResponder(http.StatusTeapot, `{ "detail": "error" }`))
			},
			id: 12234,
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    `{"detail":"error"}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			_, err := c.DeleteTag(context.Background(), tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteTag() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
