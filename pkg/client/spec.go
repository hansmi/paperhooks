package client

import (
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// OrderingSpec controls the sorting order for lists.
type OrderingSpec struct {
	// Field name, e.g. "created".
	Field string

	// Set to true for descending order. Ascending is the default.
	Desc bool
}

var _ query.Encoder = (*OrderingSpec)(nil)

func (o OrderingSpec) EncodeValues(key string, v *url.Values) error {
	if o.Field != "" {
		v.Set(key, map[bool]string{
			false: "",
			true:  "-",
		}[o.Desc]+o.Field)
	}

	return nil
}

// CharFilterSpec contains filters available on character/string fields. All
// comparison are case-insensitive.
type CharFilterSpec struct {
	EqualsIgnoringCase     *string
	StartsWithIgnoringCase *string
	EndsWithIgnoringCase   *string
	ContainsIgnoringCase   *string
}

var _ query.Encoder = (*CharFilterSpec)(nil)

func (s CharFilterSpec) EncodeValues(key string, v *url.Values) error {
	for suffix, value := range map[string]*string{
		"iexact":      s.EqualsIgnoringCase,
		"istartswith": s.StartsWithIgnoringCase,
		"iendswith":   s.EndsWithIgnoringCase,
		"icontains":   s.ContainsIgnoringCase,
	} {
		if !(value == nil || *value == "") {
			v.Set(key+"__"+suffix, *value)
		}
	}

	return nil
}

// IntFilterSpec contains filters available on numeric fields.
type IntFilterSpec struct {
	Equals *int64
	Gt     *int64
	Gte    *int64
	Lt     *int64
	Lte    *int64
	IsNull *bool
}

var _ query.Encoder = (*IntFilterSpec)(nil)

func (s IntFilterSpec) EncodeValues(key string, v *url.Values) error {
	for suffix, value := range map[string]*int64{
		"exact": s.Equals,
		"gt":    s.Gt,
		"gte":   s.Gte,
		"lt":    s.Lt,
		"lte":   s.Lte,
	} {
		if value != nil {
			v.Set(key+"__"+suffix, strconv.FormatInt(*value, 10))
		}
	}

	if s.IsNull != nil {
		v.Set(key+"__isnull", strconv.FormatBool(*s.IsNull))
	}

	return nil
}

type ForeignKeyFilterSpec struct {
	IsNull *bool
	ID     *int64
	Name   CharFilterSpec
}

var _ query.Encoder = (*ForeignKeyFilterSpec)(nil)

func (s ForeignKeyFilterSpec) EncodeValues(key string, v *url.Values) error {
	if s.IsNull != nil {
		v.Set(key+"__isnull", strconv.FormatBool(*s.IsNull))
	}

	if s.ID != nil {
		v.Set(key+"__id", strconv.FormatInt(*s.ID, 10))
	}

	return s.Name.EncodeValues(key+"__name", v)
}

type DateTimeFilterSpec struct {
	// Set to a non-nil value to only include newer items.
	Gt *time.Time

	// Set to a non-nil value to only include older items.
	Lt *time.Time
}

var _ query.Encoder = (*DateTimeFilterSpec)(nil)

func (s DateTimeFilterSpec) EncodeValues(key string, v *url.Values) error {
	for suffix, value := range map[string]*time.Time{
		"gt": s.Gt,
		"lt": s.Lt,
	} {
		if value != nil {
			v.Set(key+"__"+suffix, value.Format(time.RFC3339))
		}
	}

	return nil
}
