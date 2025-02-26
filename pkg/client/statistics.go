package client

import "context"

type Statistics struct {
	DocumentsTotal         int64                        `json:"documents_total"`
	DocumentsInbox         int64                        `json:"documents_inbox"`
	InboxTag               int64                        `json:"inbox_tag"`
	InboxTags              []int64                      `json:"inbox_tags"`
	DocumentFileTypeCounts []StatisticsDocumentFileType `json:"document_file_type_counts"`
	CharacterCount         int64                        `json:"character_count"`
	TagCount               int64                        `json:"tag_count"`
	CorrespondentCount     int64                        `json:"correspondent_count"`
	DocumentTypeCount      int64                        `json:"document_type_count"`
	StoragePathCount       int64                        `json:"storage_path_count"`
	CurrentAsn             int64                        `json:"current_asn"`
}

type StatisticsDocumentFileType struct {
	MimeType      string `json:"mime_type"`
	MimeTypeCount int64  `json:"mime_type_count"`
}

func (c *Client) GetStatistics(ctx context.Context) (*Statistics, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(&Statistics{}).
		Get("api/statistics/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*Statistics), wrapResponse(resp), nil
}
