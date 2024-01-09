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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
  defer conn.Close()
	req := make([]byte, 1024)
	conn.Read(req)
	reqData := strings.Split(string(req), "\r\n")
	path := strings.Split(reqData[0], " ")[1]
	var body string
	if path == "/" {
		ok(conn, "")
	} else if path == "/user-agent" {
		for _, line := range reqData {
			if strings.HasPrefix(line, "User-Agent") {
				body = strings.TrimPrefix(line, "User-Agent: ")
				break
			}
		}
		ok(conn, body)
	} else if strings.HasPrefix(path, "/echo/") {
		body = path[6:]
		ok(conn, body)
	} else {
		notfound(conn)
	}
}

func ok(conn net.Conn, body string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body)
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}

func notfound(conn net.Conn) {
	notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n"
	_, err := conn.Write([]byte(notFoundResponse))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}

}
