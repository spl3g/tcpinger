package pinger

import (
	"context"
	"fmt"
	"io"
	"net"
)

type CheckDNSError struct {
	Err *net.DNSError
}

func (e *CheckDNSError) Error() string {
	return fmt.Sprintf("dns error: %s", e.Err.Err)
}

func (e *CheckDNSError) Is(target error) bool {
	_, ok := target.(*CheckDNSError)
	return ok
}

func (e *CheckDNSError) Unwrap() error {
	return e.Err
}

type Pinger interface {
	Check(context.Context) (PingerResponse, error)
}

type PingerResponse interface {
	Ok() bool
	FormatPretty(io.Writer)
	FormatJSON(io.Writer)
}
