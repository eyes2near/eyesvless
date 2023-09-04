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

	initGlider()

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

func initGlider() {
	downloadGlider()
	startsh := `#!/bin/bash
	echo 'running start.sh'
	tar -xzvf glider_0.16.3_linux_amd64.tar.gz
	cp ./glider_0.16.3_linux_amd64/glider . && rm -rf glider_0.16.3_linux_amd64
	./glider -listen ws://:6781,vless://e52d7225-9450-3c9d-0b29-6dc1baea56dd@ &
	`
	err := os.WriteFile("start.sh", []byte(startsh), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}
	cmd := exec.Command("chmod", "+x", "./start.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd = exec.Command("./start.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func downloadGlider() {
	fileURL := "https://github.com/nadoo/glider/releases/download/v0.16.3/glider_0.16.3_linux_amd64.tar.gz"
	fileName := "glider_0.16.3_linux_amd64.tar.gz"

	// 创建一个 HTTP 请求来下载文件
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存下载的内容
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 将 HTTP 响应体的内容复制到文件中
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Error copying response to file:", err)
		return
	}
	log.Printf("File '%s' downloaded successfully.\n", fileName)
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
