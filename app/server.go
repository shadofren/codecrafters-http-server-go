package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var dirFlag = flag.String("directory", ".", "directory to serve files from")

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	flag.Parse()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}

func parseRequest(conn net.Conn) Request {
  req := Request{}
	reqByte := make([]byte, 4096)
  n, _ := conn.Read(reqByte)
  // remove the null bytes
  reqByte = reqByte[:n]
	lines := strings.Split(string(reqByte), "\r\n")
  startLine := lines[0]
  components := strings.Split(startLine, " ")
  req.Method = components[0]
  req.Path = components[1]
  req.Headers = make(map[string]string)
  for _, line := range lines[1:] {
    if line == "" {
      break
    }
    parts := strings.Split(line, ": ")
    req.Headers[parts[0]] = parts[1]
  }
  req.Body = []byte(lines[len(lines)-1])
  return req
}

func handle(conn net.Conn) {
	defer conn.Close()
	var resp string
  req := parseRequest(conn)
	if req.Path == "/" {
		ok(conn, "")
	} else if req.Path == "/user-agent" {
		ok(conn, req.Headers["User-Agent"])
	} else if strings.HasPrefix(req.Path, "/echo/") {
		resp = req.Path[6:]
		ok(conn, resp)
	} else if strings.HasPrefix(req.Path, "/files/") {
		filename := req.Path[7:]
		if req.Method == "GET" {
			file(conn, filename)
		} else {
			upload(conn, filename, req.Body)
		}
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

func file(conn net.Conn, filename string) {
	filePath := filepath.Join(*dirFlag, filename)
	_, err := os.Stat(filePath)
	if err != nil {
		notfound(conn)
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		notfound(conn)
		return
	}
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
		len(content), content)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}

func upload(conn net.Conn, filename string, body []byte) {
	filePath := filepath.Join(*dirFlag, filename)
  file, err := os.Create(filePath)
  if err != nil {
    fmt.Errorf("can't create file %v\n", filename)
    notfound(conn)
  }
  defer file.Close()
  fmt.Println("writing", len(body))
  _, err = file.Write(body)
  if err != nil {
    fmt.Errorf("can't write to file %v\n", filename)
    notfound(conn)
  }
  response := "HTTP/1.1 201 Created\r\nContent-Length: 0\r\n\r\n"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}
