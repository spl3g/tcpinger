package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spl3g/tcpinger/pinger"
)

type OutputFormat int

const (
	FormatJSON OutputFormat = iota
	FormatPretty
)

func OutputFormatFromString(str string) (OutputFormat, error) {
	switch str {
	case "json":
		return FormatJSON, nil
	case "pretty":
		return FormatPretty, nil
	default:
		return -1, fmt.Errorf("unknown format: %s", str)
	}
}

type CliHandler struct {
	pinger  pinger.Pinger
	timeout time.Duration
	format  OutputFormat
}

func NewCliHandler(pinger pinger.Pinger, timeout time.Duration, format OutputFormat) CliHandler {
	return CliHandler{pinger, timeout, format}
}

func (c *CliHandler) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.pinger.Check(ctx)
	if err != nil {
		return err
	}

	switch c.format {
	case FormatJSON:
		resp.FormatJSON(os.Stdout)
	case FormatPretty:
		resp.FormatPretty(os.Stdout)
	}

	return nil
}
