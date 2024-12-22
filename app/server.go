package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	msg := "HTTP/1.1 404 Not Found\r\n\r\n"
	req := make([]byte, 1024)

	conn.Read(req)

	lines := strings.Split(string(req), "\r\n")

	defer conn.Close()

	var path string
	items := []string{}

	for _, line := range lines {
		if strings.Contains(line, "GET") {
			path = strings.Split(line, " ")[1]
		} else if strings.Contains(line, "User-Agent") {
			items = strings.Split(line, ": ")
		}
	}
	if path == "" {
		msg = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo/") {
		keyword := strings.Split(path, "/echo/")[1]
		msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(keyword), keyword)
	} else if strings.HasPrefix(path, "/user-agent") && len(items) > 0 {
		msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(items[1]), items[1])
	}

	conn.Write([]byte(msg))
}
