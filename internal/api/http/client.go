package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Post(url string, body []byte) error {
	res, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unexpected status: %s", res.Status))
	}

	return nil
}
