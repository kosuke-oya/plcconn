# plcconn
This is sinmple []byte sending interface by golang

## feature
mutex, tcp connection, timeout, error handling

# how to use
go get at first
```bash
go get github.com/kosuke-oya/plcconn
```
    
sample
```go
package main

import (
    "github.com/kosuke-oya/plcconn"
)

func main() {
	ipAddress := "192.168.1.1"
	port := 1025
	timeOutSecond := 5

	// create plc connection
	conn := plcconn.NewPlcConn(ipAddress, port, timeOutSecond)
	// defer close connection
	defer conn.Close()

	// send []byte data to plc
	data := []byte{0x01, 0x02, 0x03, 0x04}
	r,err := conn.Send(data)
        if err != nil {
           fmt.Println(err)
        }

        // show result
        fmt.Printf("%d",r)
}

```
