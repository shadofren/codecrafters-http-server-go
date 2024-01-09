package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	handle(conn)
}

func handle(conn net.Conn) {
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	reqData := strings.Split(string(req), "\r\n")
	path := strings.Split(reqData[0], " ")[1]
	okResponse := "HTTP/1.1 200 OK\r\n\r\n"
	notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n"
	if path == "/" {
		_, err = conn.Write([]byte(okResponse))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/echo/") {
		body := path[6:]
		size := len(body)
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			size, body)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}

	} else {
		_, err = conn.Write([]byte(notFoundResponse))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	}

}
