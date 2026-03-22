package ufw

import (
	"context"
	"fmt"
	"strings"
)

type Action string

const (
	ActionAllow  Action = "allow"
	ActionDeny   Action = "deny"
	ActionReject Action = "reject"
	ActionLimit  Action = "limit"
)

type Config struct {
	// Reserved for future options (sudo, timeouts, etc).
}

type Runner interface {
	Run(ctx context.Context, name string, args ...string) (stdout string, stderr string, exitCode int, err error)
}

type Client struct {
	runner Runner
	cfg    Config
}

func NewClient(runner Runner, cfg Config) *Client {
	return &Client{runner: runner, cfg: cfg}
}

func (c *Client) Status(ctx context.Context) (Status, error) {
	out, err := c.run(ctx, "status")
	if err != nil {
		return Status{}, err
	}
	return ParseStatus(out)
}

func (c *Client) StatusNumbered(ctx context.Context) (NumberedStatus, error) {
	out, err := c.run(ctx, "status", "numbered")
	if err != nil {
		return NumberedStatus{}, err
	}
	return ParseNumberedStatus(out)
}

func (c *Client) StatusVerbose(ctx context.Context) (string, error) {
	return c.run(ctx, "status", "verbose")
}

func (c *Client) AppList(ctx context.Context) ([]string, error) {
	out, err := c.run(ctx, "app", "list")
	if err != nil {
		return nil, err
	}
	return ParseAppList(out)
}

func (c *Client) AppInfo(ctx context.Context, profile string) (string, error) {
	if strings.TrimSpace(profile) == "" {
		return "", fmt.Errorf("empty profile")
	}
	return c.run(ctx, "app", "info", profile)
}

func (c *Client) ApplyApp(ctx context.Context, action Action, profile string) error {
	switch action {
	case ActionAllow, ActionDeny, ActionReject, ActionLimit:
	default:
		return fmt.Errorf("unsupported action: %q", action)
	}
	if strings.TrimSpace(profile) == "" {
		return fmt.Errorf("empty profile")
	}
	_, err := c.run(ctx, string(action), profile)
	return err
}

func (c *Client) Enable(ctx context.Context) error {
	_, err := c.run(ctx, "--force", "enable")
	return err
}

func (c *Client) Disable(ctx context.Context) error {
	_, err := c.run(ctx, "disable")
	return err
}

func (c *Client) Reload(ctx context.Context) error {
	_, err := c.run(ctx, "reload")
	return err
}

func (c *Client) DeleteNumber(ctx context.Context, number int) error {
	if number <= 0 {
		return fmt.Errorf("invalid rule number: %d", number)
	}
	_, err := c.run(ctx, "--force", "delete", fmt.Sprintf("%d", number))
	return err
}

func (c *Client) Apply(ctx context.Context, action Action, args []string) error {
	switch action {
	case ActionAllow, ActionDeny, ActionReject, ActionLimit:
	default:
		return fmt.Errorf("unsupported action: %q", action)
	}
	if len(args) == 0 {
		return fmt.Errorf("empty args")
	}
	argv := append([]string{string(action)}, args...)
	_, err := c.run(ctx, argv...)
	return err
}

func (c *Client) run(ctx context.Context, args ...string) (string, error) {
	stdout, stderr, code, err := c.runner.Run(ctx, "ufw", args...)
	if err != nil {
		return "", err
	}
	if code != 0 {
		if stderr == "" {
			stderr = stdout
		}
		return "", fmt.Errorf("ufw %v failed (exit=%d): %s", args, code, trimOneLine(stderr))
	}
	out := stdout
	if strings.TrimSpace(out) == "" && strings.TrimSpace(stderr) != "" {
		out = stderr
	}
	return out, nil
}

func trimOneLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return s
}
