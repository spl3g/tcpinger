package main

import (
	"time"

	"github.com/spl3g/tcpinger/cli"
	"github.com/spl3g/tcpinger/pinger"
)

func handleCli(destination string, port uint, timeout time.Duration, formatStr string) error {
	pinger := pinger.NewTCPPinger(destination, port)
	format, err := cli.OutputFormatFromString(formatStr)
	if err != nil {
		return err
	}

	handler := cli.NewCliHandler(&pinger, timeout, format)
	return handler.Run()
}
