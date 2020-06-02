package onfido

import (
	"context"
	"net/http"
)

func (c *client) getResourceFromHref(ctx context.Context, href string, v interface{}) error {
	req, err := c.newRequest(http.MethodGet, href, nil)
	if err != nil {
		return err
	}
	_, err = c.do(ctx, req, v)
	return err
}

func (c *client) GetCheckFromHref(ctx context.Context, href string) (*Check, error) {
	var resp Check
	if err := c.getResourceFromHref(ctx, href, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
