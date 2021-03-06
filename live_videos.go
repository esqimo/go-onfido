package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LiveVideo represents a live video object in Onfido API
// https://documentation.onfido.com/#live-video-object
type LiveVideo struct {
	ID           string     `json:"id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Href         string     `json:"href,omitempty"`
	DownloadHref string     `json:"download_href,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	FileType     string     `json:"file_type,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
}

type LiveVideoDownload struct {
	// Data is the binary data of the live video
	Data []byte
}

// DownloadLiveVideo returns the binary data representing the video.
// see https://documentation.onfido.com/#download-live-video
func (c *client) DownloadLiveVideo(ctx context.Context, id string) (*LiveVideoDownload, error) {
	req, err := c.newRequest(http.MethodGet, "/live_videos/"+id+"/download", nil)
	if err != nil {
		return nil, err
	}

	var resp bytes.Buffer
	_, err = c.do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to download live video: %w", err)
	}
	return &LiveVideoDownload{
		Data: resp.Bytes(),
	}, err
}

// liveVideoIter represents a LiveVideo iterator
type liveVideoIter struct {
	*iter
}

type LiveVideoIter interface {
	Iter
	LiveVideo() *LiveVideo
}

func (i *liveVideoIter) LiveVideo() *LiveVideo {
	return i.Current().(*LiveVideo)
}

// LiveVideoIter retrieves the list of live videos for the provided applicant.
// see https://documentation.onfido.com/#list-live-videos
func (c *client) ListLiveVideos(applicantID string) LiveVideoIter {
	return &liveVideoIter{&iter{
		c:       c,
		nextURL: "/live_videos?applicant_id=" + applicantID,
		handler: func(body []byte) ([]interface{}, error) {
			var r struct {
				LiveVideos []*LiveVideo `json:"live_videos"`
			}

			if err := json.Unmarshal(body, &r); err != nil {
				return nil, err
			}

			values := make([]interface{}, len(r.LiveVideos))
			for i, v := range r.LiveVideos {
				values[i] = v
			}
			return values, nil
		},
	}}
}
