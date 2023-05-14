package client

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-querystring/query"
)

func TestSpecs(t *testing.T) {
	type FakeOrderBy struct {
		Order OrderingSpec `url:"order_by"`
	}

	type FakeChar struct {
		Title CharFilterSpec `url:"title"`
	}

	type FakeInt struct {
		Number IntFilterSpec `url:"number"`
	}

	type FakeForeignKey struct {
		Kind ForeignKeyFilterSpec `url:"kind"`
	}

	type FakeDateTime struct {
		Created DateTimeFilterSpec `url:"created"`
	}

	for _, tc := range []struct {
		name    string
		value   any
		wantErr error
		want    url.Values
	}{
		{
			name: "empty",
		},
		{
			name: "ordering asc",
			value: FakeOrderBy{
				OrderingSpec{
					Field: "title",
				},
			},
			want: url.Values{
				"order_by": []string{"title"},
			},
		},
		{
			name: "ordering desc",
			value: FakeOrderBy{
				OrderingSpec{
					Field: "title",
					Desc:  true,
				},
			},
			want: url.Values{
				"order_by": []string{"-title"},
			},
		},
		{
			name: "char iexact",
			value: FakeChar{
				Title: CharFilterSpec{
					EqualsIgnoringCase: String("xyz"),
				},
			},
			want: url.Values{
				"title__iexact": []string{"xyz"},
			},
		},
		{
			name: "char all",
			value: FakeChar{
				Title: CharFilterSpec{
					EqualsIgnoringCase:     String("equals"),
					StartsWithIgnoringCase: String("startswith"),
					EndsWithIgnoringCase:   String("endswith"),
					ContainsIgnoringCase:   String("contains"),
				},
			},
			want: url.Values{
				"title__iexact":      []string{"equals"},
				"title__istartswith": []string{"startswith"},
				"title__iendswith":   []string{"endswith"},
				"title__icontains":   []string{"contains"},
			},
		},
		{
			name: "int",
			value: FakeInt{
				Number: IntFilterSpec{
					Equals: Int64(300),
					Gt:     Int64(400),
					Gte:    Int64(401),
					Lt:     Int64(500),
					Lte:    Int64(501),
					IsNull: Bool(false),
				},
			},
			want: url.Values{
				"number__exact":  []string{"300"},
				"number__gt":     []string{"400"},
				"number__gte":    []string{"401"},
				"number__lt":     []string{"500"},
				"number__lte":    []string{"501"},
				"number__isnull": []string{"false"},
			},
		},
		{
			name: "foreign key",
			value: FakeForeignKey{
				Kind: ForeignKeyFilterSpec{
					ID:     Int64(123),
					IsNull: Bool(true),
					Name: CharFilterSpec{
						EqualsIgnoringCase:     String("equals"),
						StartsWithIgnoringCase: String("startswith"),
						EndsWithIgnoringCase:   String("endswith"),
						ContainsIgnoringCase:   String("contains"),
					},
				},
			},
			want: url.Values{
				"kind__id":                []string{"123"},
				"kind__isnull":            []string{"true"},
				"kind__name__iexact":      []string{"equals"},
				"kind__name__istartswith": []string{"startswith"},
				"kind__name__iendswith":   []string{"endswith"},
				"kind__name__icontains":   []string{"contains"},
			},
		},
		{
			name: "datetime",
			value: FakeDateTime{
				Created: DateTimeFilterSpec{
					Lt: Time(time.Date(2015, time.March, 7, 1, 2, 3, 0, time.UTC)),
					Gt: Time(time.Date(2018, time.July, 9, 4, 5, 6, 0, time.UTC)),
				},
			},
			want: url.Values{
				"created__lt": []string{"2015-03-07T01:02:03Z"},
				"created__gt": []string{"2018-07-09T04:05:06Z"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := query.Values(tc.value)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("query.Values() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("Encoded query diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
