package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type listResult[T any] struct {
	// Total item count.
	Count int64 `json:"count"`

	// URL for next page (if any).
	Next string `json:"next"`

	// URL for previous page (if any).
	Previous string `json:"previous"`

	Items []T `json:"results"`
}

type crudOptions struct {
	newRequest func(context.Context) *resty.Request
	base       string
}

func crudList[T any](ctx context.Context, opts crudOptions, listOpts any) ([]T, *Response, error) {
	req := opts.newRequest(ctx).SetResult(new(listResult[T]))

	if values, err := query.Values(listOpts); err != nil {
		return nil, nil, err
	} else {
		req.SetQueryParamsFromValues(values)
	}

	resp, err := req.Get(opts.base)

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	results := resp.Result().(*listResult[T])

	w := wrapResponse(resp)

	if w.NextPage, err = pageTokenFromURL(results.Next); err != nil {
		return nil, nil, err
	}

	if w.PrevPage, err = pageTokenFromURL(results.Previous); err != nil {
		return nil, nil, err
	}

	return results.Items, w, nil
}

func crudGet[T any](ctx context.Context, opts crudOptions, id int64) (*T, *Response, error) {
	resp, err := opts.newRequest(ctx).
		SetResult(new(T)).
		Get(fmt.Sprintf("%s%d/", opts.base, id))

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*T), wrapResponse(resp), nil
}

func crudCreate[T any](ctx context.Context, opts crudOptions, data *T) (*T, *Response, error) {
	resp, err := opts.newRequest(ctx).
		SetResult(new(T)).
		SetBody(*data).
		Post(opts.base)

	err = convertError(err, resp)

	if detail, ok := err.(*RequestError); ok && detail.StatusCode == http.StatusCreated {
		return resp.Result().(*T), wrapResponse(resp), nil
	}

	if err == nil {
		err = &RequestError{
			StatusCode: resp.StatusCode(),
			Message:    fmt.Sprintf("unexpected status %s", resp.Status()),
		}
	}

	return nil, wrapResponse(resp), err
}

func crudUpdate[T any](ctx context.Context, opts crudOptions, id int64, data *T) (*T, *Response, error) {
	resp, err := opts.newRequest(ctx).
		SetResult(new(T)).
		SetBody(*data).
		Put(fmt.Sprintf("%s%d/", opts.base, id))

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*T), wrapResponse(resp), nil
}

func crudDelete[T any](ctx context.Context, opts crudOptions, id int64) (*Response, error) {
	resp, err := opts.newRequest(ctx).
		Delete(fmt.Sprintf("%s%d/", opts.base, id))

	return wrapResponse(resp), convertError(err, resp)
}
