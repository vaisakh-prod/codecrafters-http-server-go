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

func handleConnection(conn net.Conn, Directory *string) {
	msg := "HTTP/1.1 404 Not Found\r\n\r\n"
	req := make([]byte, 1024)

	conn.Read(req)

	lines := strings.Split(string(req), "\r\n")

	defer conn.Close()

	var path string
	items := []string{}

	fmt.Println(lines)
	for _, line := range lines {
		if strings.Contains(line, "GET") || strings.Contains(line, "POST") {
			path = strings.Split(line, " ")[1]
		} else if strings.Contains(line, "User-Agent") {
			items = strings.Split(line, ": ")
		}
	}
	fmt.Println(lines[0])
	if strings.Contains(lines[0], "GET") {
		fmt.Println(path)
		if path == "/" {
			msg = "HTTP/1.1 200 OK\r\n\r\n"
		} else if strings.HasPrefix(path, "/echo/") {
			keyword := strings.Split(path, "/echo/")[1]
			msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(keyword), keyword)
		} else if strings.HasPrefix(path, "/user-agent") && len(items) > 0 {
			msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%v", len(items[1]), items[1])
		} else if strings.HasPrefix(path, "/files/") {
			fileName := strings.Split(path, "/files/")[1]
			if *Directory != "" && fileName != "" {
				absolutePath := filepath.Join(*Directory, fileName)

				fileData, err := os.ReadFile(absolutePath)
				//fmt.Println(absolutePath, fileData)
				if err != nil {
					fmt.Print(err)
				}
				if err == nil {
					file := string(fileData)
					msg = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%v", len(file), file)
				}
			} else {
				panic("Directory or fileName is None")
			}
		}
	}
	conn.Write([]byte(msg))
}
