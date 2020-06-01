package onfido

import (
	"context"
	"net/http"
)

func (c *Client) getHref(ctx context.Context, href string, v interface{}) error {
	req, err := c.newRequest(http.MethodGet, href, nil)
	if err != nil {
		return err
	}
	_, err = c.do(ctx, req, v)
	return err
}

func (c *Client) GetCheckFromHref(ctx context.Context, href string) (*Check, error) {
	var resp Check
	if err := c.getHref(ctx, href, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetReportFromHref(ctx context.Context, href string) (*Report, error) {
	var resp Report
	if err := c.getHref(ctx, href, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetDocumentFromHref(ctx context.Context, href string) (*Document, error) {
	var resp Document
	if err := c.getHref(ctx, href, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
