package main

func main() {
	ipAddress := "192.168.1.1"
	port := 1025
	timeOutSecond := 5

	// create plc connection
	conn := NewPlcConn(ipAddress, port, timeOutSecond)
	// defer close connection
	defer conn.Close()

	// send []byte data to plc
	data := []byte{0x01, 0x02, 0x03, 0x04}
	conn.Send(data)
}
