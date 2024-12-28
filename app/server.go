package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var cmdLine = flag.String("directory", "", "Specify the folder in which files are to be handled")

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.Parse()

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
		go handleConnection(conn, cmdLine)
	}

}

func handleConnection(conn net.Conn, directory *string) {
	defer conn.Close()

	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		sendResponse(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
		return
	}
	req = req[:n]
	lines := strings.Split(string(req), "\r\n")

	path, userAgent, encodingFormats := parseRequest(lines)
	if path == "" {
		sendResponse(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return
	}

	switch {
	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path, encodingFormats)
	case strings.HasPrefix(path, "/user-agent"):
		handleUserAgent(conn, userAgent)
	case strings.HasPrefix(path, "/files/"):
		handleFiles(conn, directory, path, lines)
	default:
		handleDefault(conn, path)
	}
}

func handleEcho(conn net.Conn, path string, encodingFormats []string) {
	keyword := strings.TrimPrefix(path, "/echo/")
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(keyword), keyword)
	for _, encodingFormat := range encodingFormats {
		if encodingFormat == "gzip" {
			encodedKeyword, err := gzipAndEncode(keyword)
			if err != nil {
				fmt.Printf("Failed to gzip and encode: %v\n", err)
				sendResponse(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
				return
			}
			msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%v", len(encodedKeyword), encodedKeyword)
			break
		}
	}
	sendResponse(conn, msg)
}

func handleUserAgent(conn net.Conn, userAgent string) {
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(userAgent), userAgent)
	sendResponse(conn, msg)
}

func handleFiles(conn net.Conn, directory *string, path string, lines []string) {
	fileName := strings.TrimPrefix(path, "/files/")
	if *directory == "" || fileName == "" {
		sendResponse(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
		return
	}
	absolutePath := filepath.Join(*directory, fileName)

	if strings.HasPrefix(lines[0], "GET") {
		fileData, err := os.ReadFile(absolutePath)
		if err != nil {
			sendResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
			return
		}
		msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%v", len(fileData), string(fileData))
		sendResponse(conn, msg)
	} else if strings.HasPrefix(lines[0], "POST") {
		data := []byte(lines[len(lines)-1])
		err := os.WriteFile(absolutePath, data, 0644)
		if err != nil {
			sendResponse(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
			return
		}
		sendResponse(conn, "HTTP/1.1 201 Created\r\n\r\n")
	}
}

func handleDefault(conn net.Conn, path string) {
	if path == "/" {
		sendResponse(conn, "HTTP/1.1 200 OK\r\n\r\n")
	} else {
		sendResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}
