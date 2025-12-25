package pinger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"text/tabwriter"
)

type TCPPinger struct {
	Host string
	Port uint
}

func NewTCPPinger(ip string, port uint) TCPPinger {
	return TCPPinger{ip, port}
}

func (p *TCPPinger) Check(ctx context.Context) (PingerResponse, error) {
	resp := TCPPingerResponse{p.Host, p.Port, false}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(p.Host, fmt.Sprint(p.Port)))

	var dnsError *net.DNSError
	if errors.Is(err, syscall.ECONNREFUSED) {
		return resp, nil
	} else if errors.As(err, &dnsError) {
		return resp, &CheckDNSError{dnsError}
	} else if err != nil {
		return resp, err
	}
	conn.Close()

	resp.ok = true
	return resp, nil
}

type TCPPingerResponse struct {
	Host string
	Port uint
	ok   bool
}

func (p TCPPingerResponse) Ok() bool {
	return p.ok
}

func (p TCPPingerResponse) FormatPretty(w io.Writer) {
	var status string
	if p.ok {
		status = "open"
	} else {
		status = "closed"
	}

	fmt.Fprintf(w, "Host %s:\n", p.Host)
	tw := tabwriter.NewWriter(w, 1, 1, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tSTATUS")
	fmt.Fprintf(tw, "%d/tcp\t%s\n", p.Port, status)
	tw.Flush()
}

type tcpJsonOutput struct {
	IP     string `json:"ip"`
	Port   uint   `json:"port"`
	Status string `json:"status"`
}

func (p TCPPingerResponse) FormatJSON(w io.Writer) {
	var status string
	if p.ok {
		status = "open"
	} else {
		status = "closed"
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(tcpJsonOutput{p.Host, p.Port, status})
	if err != nil {
		panic("Failed to marshal tcp output somehow")
	}
}
