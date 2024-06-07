package main

import (
	"UPC-GO/api"
	"UPC-GO/register"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	success := register.RegisterService()
	if success {
		fmt.Println("Service registered successfully")
	} else {
		fmt.Println("Service registration failed")
	}

	// 循环发送心跳
	// 创建一个定时器，每隔HeartbeatInterval发送一次心跳信号
	HeartbeatInterval := 60 * time.Second
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				success := register.SendHeartbeat()
				if success {
					fmt.Println("Heartbeat sent successfully")
				} else {
					fmt.Println("Failed to send heartbeat")
				}
			}
		}
	}()

	log.Println("Starting server on :4000")
	if err := StartServer(":4000"); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

// HTTP服务器
type Server struct {
}

// NewServer 创建一个新的服务器实例
func NewServer() *Server {
	return &Server{}
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	// 静态文件服务器 /static/ -> ./public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/ws", api.HandleWebSocket)

	// 路由
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/api", ConnectHandler)
	http.HandleFunc("/api/files", api.FilesHandler)             // get /api/files 获取所有文件的列表
	http.HandleFunc("/api/files/", api.FileProcessor)           // get /api/files/:filename 对一个文件进行操作
	http.HandleFunc("/api/files/download", api.MultiDownloader) // get /api/files/download 下载多个文件
	http.HandleFunc("/api/results", api.ResultsHandler)         // get /api/results 获取所有结果的列表
	http.HandleFunc("/api/results/", api.ResultProcessor)       // get /api/results/:resultName 下载或删除一个结果
	http.HandleFunc("/api/upload", api.UploadHandler)           // post /api/upload 上传文件
	http.HandleFunc("/api/images", api.ImagesHandler)           // get /api/images 获取所有docker images 的列表
	http.HandleFunc("/api/images/", api.ImageProcessor)         // get /api/images/:imageName 对一个docker image 进行操作
	http.HandleFunc("/api/pull/", api.ImagePuller)              // post /api/pull/:imageName 拉取一个docker image

	// 创建一个 http.Server 实例
	server := &http.Server{Addr: addr}

	// 启动服务器的 Goroutine，这样我们可以在主线程中等待服务器关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// 等待服务器关闭
	waitForShutdown(server)
	return nil
}

// 启动服务器
func StartServer(addr string) error {
	server := NewServer()
	return server.Start(addr)
}

// 等待服务器关闭
func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 关闭服务器，注销服务
	register.UnregisterService()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

// 返回index.html
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	api.Cors(w)
	http.ServeFile(w, r, "./views/index.html")
}

// get /api 是测试这个服务器是否正常工作
func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	api.Cors(w)
	w.Write([]byte("Connect success!"))
}
