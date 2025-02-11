# go-ipc

<!-- 
go test -coverpkg ./... -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html
-->

[![Go Coverage](https://github.com/Kilemonn/go-ipc/wiki/coverage.svg)](https://raw.githack.com/wiki/Kilemonn/go-ipc/coverage.html)


An inter-process communication library written in Golang. Using "unix" sockets to support interprocess communication. **This also works on Windows!**

## Get started

You can add this project as a dependency using:

> go get github.com/Kilemonn/go-ipc

Once added you can setup connections as per below (error handling removed for brevity):
```go
package main

import (
    "time"

	"github.com/Kilemonn/go-ipc/client"
	"github.com/Kilemonn/go-ipc/server"
)

func main() {
	ipcChannelName := "readme-example"
	svr, _ := server.NewIPCServer(ipcChannelName, nil)
	defer svr.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	defer client.Close()

	accepted, err := svr.Accept(time.Millisecond * 1000)
	defer accepted.Close()

	content := "some-data"
	n, err := client.Write([]byte(content))

	b := make([]byte, len(content))
	n, err = accepted.Read(b)

	if content == string(b) {
		fmt.Println("Sent and read is the same!")
	}
}
```

## Notes

- The `IPCClient` implements the `io.ReadWriterCloser` interface and provides a property to set read timeouts instead of remaining blocking.
- This is aimed to be a very primative and simple implementation for now, other settings and customisability can be built upon what is currently implemented. However feel free to request any features that would be useful.
