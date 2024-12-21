package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := "HTTP/1.1 200 OK\r\n\r\n"
		if r.URL.Path != "/" {
			msg = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
		w.Write([]byte(msg))
	})

	err := http.ListenAndServe("0.0.0.0:4221", nil)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
}
