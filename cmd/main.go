package main

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/go-playground/validator"
)

type CliCmd struct {
	Format  string        `short:"f" enum:"json,pretty" default:"json" help:"Output format: json or pretty."`
	Timeout time.Duration `short:"t" default:"1s" help:"Timeout duration."`

	Destination string `arg:"" help:"IP or host to check the port on." validate:"ip|hostname"`
	Port        uint   `arg:"" help:"Port to check." validate:"gte=0,lte=65535"`
}

func (c *CliCmd) Run() error {
	return handleCli(c.Destination, c.Port, c.Timeout, c.Format)
}

func (c *CliCmd) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

type ServerCmd struct {
	Host string `short:"H" default:"0.0.0.0" help:"Host to serve on, can be set with env." env:"HOST"`
	Port int    `short:"P" default:"8072" help:"Port to serve on, can be set with env." env:"PORT"`
}

func (c *ServerCmd) Run() error {
	return handleServe(c.Host, c.Port)
}

var commands struct {
	Serve ServerCmd `cmd:"" help:"Start a server that will listen for ips and ports."`
	Cli   CliCmd    `cmd:"" help:"Test if the port on an ip is open."`
}

func main() {
	ctx := kong.Parse(&commands, kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
