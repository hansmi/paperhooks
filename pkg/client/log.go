package client

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// ListLogs retrieves the names of available log files.
func (c *Client) ListLogs(ctx context.Context) ([]string, *Response, error) {
	req := c.newRequest(ctx).SetResult([]string(nil))

	resp, err := req.Get("api/logs/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return *resp.Result().(*[]string), wrapResponse(resp), nil
}

type LogEntry struct {
	Time    time.Time
	Level   string
	Module  string
	Message string
}

func (e *LogEntry) appendLine(line string) {
	e.Message += "\n" + line
}

// Regular expression matching a log message. Example:
// [2023-02-28 00:28:37,604] [INFO] [paperless.consumer] Consuming xyz.pdf"
var logEntryRe = regexp.MustCompile(`^` +
	`\[(?P<time>\d\d\d\d-\d\d-\d\d\s+\d\d:\d\d:\d\d(?:[.,]\d{1,6})?)\]\s+` +
	`\[(?P<level>[A-Z]{1,20})\]\s+` +
	`\[(?P<module>[^\]]{1,64})\]\s?` +
	`(?P<message>.*)` +
	`$`)

type logParser struct {
	loc *time.Location
}

func (p *logParser) parseTime(value string) time.Time {
	for _, layout := range []string{
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05",
	} {
		if ts, err := time.ParseInLocation(layout, value, p.loc); err == nil {
			return ts
		}
	}

	return time.Time{}
}

func (p *logParser) detectStart(line string) *LogEntry {
	groups := logEntryRe.FindStringSubmatch(line)
	if len(groups) < 4 {
		return nil
	}

	return &LogEntry{
		Time:    p.parseTime(groups[1]),
		Level:   groups[2],
		Module:  groups[3],
		Message: groups[4],
	}
}

func (p *logParser) parse(lines []string) []LogEntry {
	var result []LogEntry
	var current *LogEntry

	for _, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)

		if entry := p.detectStart(line); entry != nil {
			if current != nil {
				result = append(result, *current)
			}

			current = entry
		} else if current != nil {
			current.appendLine(line)
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result
}

// GetLog retrieves all entries of the named log file.
func (c *Client) GetLog(ctx context.Context, name string) ([]LogEntry, *Response, error) {
	req := c.newRequest(ctx).SetResult([]string(nil))

	resp, err := req.Get(fmt.Sprintf("api/logs/%s/", url.PathEscape(name)))

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	lines := *resp.Result().(*[]string)

	p := logParser{
		loc: c.loc,
	}

	return p.parse(lines), wrapResponse(resp), nil
}
