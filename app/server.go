package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	msg := "HTTP/1.1 404 Not Found\r\n\r\n"
	// Uncomment this block to pass the first stage
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
	req := make([]byte, 1024)
	conn.Read(req)
	fmt.Print(string(req))
	lines := strings.Split(string(req), "\r\n")
	path_str := strings.Split(lines[0], " ")[1]
	if path_str == "/" {
		msg = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path_str, "/echo/") {
		keyword := strings.Split(path_str, "/echo/")[1]
		msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(keyword), keyword)
	} else if strings.HasPrefix(path_str, "/user-agent") && len(lines) > 4 {
		items := strings.Split(lines[3], ": ")
		msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(items[1]), items[1])
	}
	conn.Write([]byte(msg))
}
