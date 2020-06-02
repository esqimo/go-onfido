package onfido

import (
	"context"
	"net/http"
)

func (c *client) GetResource(ctx context.Context, href string, v interface{}) error {
	req, err := c.newRequest(http.MethodGet, href, nil)
	if err != nil {
		return err
	}
	_, err = c.do(ctx, req, v)
	return err
}
