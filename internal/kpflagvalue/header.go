package kpflagvalue

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

type httpHeader http.Header

var _ kingpin.Value = (*httpHeader)(nil)

func (h *httpHeader) String() string {
	return fmt.Sprintf("%q", *(*http.Header)(h))
}

func (h *httpHeader) IsCumulative() bool {
	return true
}

func (h *httpHeader) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected header:value, got %q", value)
	}

	if *h == nil {
		*h = httpHeader{}
	}

	p := (*http.Header)(h)
	p.Add(parts[0], parts[1])

	return nil
}

func HTTPHeaderVar(t kingpin.Settings, target *http.Header) {
	t.SetValue((*httpHeader)(target))
}
