package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"strings"
)

func sendResponse(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Failed to send response: %v\n", err)
	}
}

func parseRequest(lines []string) (string, string, []string) {
	var path, userAgent string
	encodingFormats := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "GET") || strings.HasPrefix(line, "POST") {
			path = strings.Fields(line)[1]
		} else if strings.HasPrefix(line, "User-Agent") {
			userAgent = strings.SplitN(line, ": ", 2)[1]
		} else if strings.HasPrefix(line, "Accept-Encoding") {

			encodingFormat := strings.SplitN(line, ": ", 2)[1]
			encodingFormats = strings.Split(encodingFormat, ", ")
		}
	}
	return path, userAgent, encodingFormats
}

func gzipAndEncode(data string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write([]byte(data))
	if err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return b.String(), nil
}
