package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"golang.org/x/sync/errgroup"
)

func defaultListOpts[T any](opts *T) *T {
	if opts == nil {
		var zero T
		opts = &zero
	}

	return opts
}

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
	getID      func(any) int64
	setPage    func(any, *PageToken)
}

func crudList[T, O any](ctx context.Context, opts crudOptions, listOpts *O) ([]T, *Response, error) {
	listOpts = defaultListOpts(listOpts)

	req := opts.newRequest(ctx).SetResult(new(listResult[T]))

	var pageNumber int

	if values, err := query.Values(listOpts); err != nil {
		return nil, nil, err
	} else {
		req.SetQueryParamsFromValues(values)

		// listOpts has a property with the page number. Getting its value
		// directly would require reflection.
		if page := values.Get("page"); page != "" {
			if pageNumber, err = strconv.Atoi(page); err != nil {
				return nil, nil, err
			}
		}
	}

	resp, err := req.Get(opts.base)

	// Items modified or deleted during iteration can cause a received page
	// number to become unavailable. Treat the situation as an empty result
	// set.
	if err == nil && resp.StatusCode() == http.StatusNotFound && pageNumber > 1 {
		return nil, wrapResponse(resp), nil
	}

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

func crudListAll[T, O any](ctx context.Context, opts crudOptions, listOpts *O, handler func(context.Context, T) error) error {
	listOpts = defaultListOpts(listOpts)

	queue := make(chan []T, 2)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(queue)

		for {
			items, resp, err := crudList[T](ctx, opts, listOpts)
			if err != nil {
				return err
			}

			queue <- items

			if resp == nil || resp.NextPage == nil {
				break
			}

			opts.setPage(listOpts, resp.NextPage)
		}

		return nil
	})

	g.Go(func() error {
		seen := map[int64]struct{}{}

		for items := range queue {
			for _, i := range items {
				key := opts.getID(i)

				if _, ok := seen[key]; ok {
					// Duplicates may be returned when items are added,
					// modified or deleted during iteration.
					continue
				}

				seen[key] = struct{}{}

				if err := handler(ctx, i); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return g.Wait()
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
