package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, this is eyesvless"))
	})

	http.HandleFunc("/ps", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		w.Write([]byte("executed ps"))
	})

	http.HandleFunc("/vless", vlessWsHandler)

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func vlessWsHandler(w http.ResponseWriter, r *http.Request) {
	// 定义服务器地址和端口
	serverAddr := "127.0.0.1:6781"
	// 创建一个 buffer 来保存请求数据
	var buf bytes.Buffer
	// 将整个请求写入到 buffer
	r.Host = serverAddr
	if err := r.Write(&buf); err != nil {
		http.Error(w, "Error writing request to buffer", http.StatusInternalServerError)
		return
	}
	// 将 buffer 的内容转换为 []byte
	requestBytes := buf.Bytes()
	h, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "websocket: response does not implement http.Hijacker", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := h.Hijack()
	if err != nil {
		http.Error(w, "cannot get underlaying connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()
	// 连接后端的vless服务，并且开始tcp代理
	// 连接到服务器
	vlessConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer vlessConn.Close()
	vlessConn.Write(requestBytes)
	go func() {
		io.Copy(vlessConn, clientConn)
		vlessConn.Close()
	}()
	io.Copy(clientConn, vlessConn)
}
