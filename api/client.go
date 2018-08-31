package api

import (
	"context"
	"time"

	"github.com/ferux/phraseGen/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Client is used to connect to the server and making all interactions between it and himself.
type Client struct {
	connString string
	client     api.APIClient
	conn       *grpc.ClientConn
	done       chan struct{}
	l          *logrus.Entry
	isRunning  bool
}

// NewClient dials to the connection string and returns niew client with already active connection.
func NewClient(cs string) *Client {
	return &Client{
		connString: cs,
		done:       make(chan struct{}),
		l:          logrus.New().WithField("pkg", "api"),
	}
}

// Run runs the client
func (c *Client) Run() error {
	conn, err := grpc.Dial(c.connString, grpc.WithInsecure())

	if err != nil {
		return err
	}

	c.conn = conn
	c.client = api.NewAPIClient(conn)
	c.isRunning = true
	return nil
}

// Close closes the connection to the server
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetMessage from the server
func (c *Client) GetMessage() (*api.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	msg, err := c.client.GetMessage(ctx, &api.Query{})
	defer cancel()
	return msg, err
}

// AskStatus asks about server's status.
func (c *Client) AskStatus() (*api.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	msg, err := c.client.AskStatus(ctx, &api.Query{})
	defer cancel()
	return msg, err
}
