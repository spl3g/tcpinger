# TCPinger
Tests if the tcp port is open.  
There are two ways to interface with the program: as an `http server` or as a `cli`. They can be chosen with the subcommands.

CLI supports multiple output styles, chosen with `-f` flag: `json` and `pretty`.

Building:
```sh
$ go build -o tcpinger ./cmd
```

Running:
```sh
$ ./tcpinger -h
Usage: tcpinger <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  serve [flags]
    Start a server that will listen for ips and ports.

  cli <destination> <port> [flags]
    Test if the port on an ip is open.

Run "tcpinger <command> --help" for more information on a command.
```
```sh
$ ./tcpinger cli 77.88.55.242 443
{"ip":"77.88.55.242","port":443,"status":"open"}
```
```sh
$ ./tcpinger cli 77.88.55.242 443 -f pretty
Host 77.88.55.242:
PORT     STATUS
443/tcp  open
```
```sh
$ ./tcpinger serve
2025/12/25 20:57:13 Started HTTP server on 0.0.0.0:8072
...
```

The HTTP server API consists of:
- POST /api/v1/check  
  Request:
  ```json
  {
    "ip": "77.88.55.242",
    "port": 443,
    "timeout": "1s"
  }
  ```
