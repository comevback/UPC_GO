package main

import (
	"UPC-GO/api"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting server on :4000")
	if err := StartServer(":4000"); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	// 静态文件服务器 /static/ -> ./public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 路由
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/api", ConnectHandler)
	http.HandleFunc("/api/files", api.FilesHandler)             // get /api/files 获取所有文件的列表
	http.HandleFunc("/api/files/", api.FileProcessor)           // get /api/files/:filename 下载或删除一个文件
	http.HandleFunc("/api/files/download", api.MultiDownloader) // get /api/files/download 下载多个文件
	http.HandleFunc("/api/results", api.ResultsHandler)         // get /api/results 获取所有结果的列表
	http.HandleFunc("/api/results/", api.ResultProcessor)       // get /api/results/:resultName 下载或删除一个结果
	http.HandleFunc("/api/upload", api.UploadHandler)           // post /api/upload 上传文件
	http.HandleFunc("/api/images", api.ImagesHandler)           // get /api/images 获取所有docker images 的列表
	http.HandleFunc("/api/images/", api.ImageProcessor)         // get /api/images/:imageName 获取一个docker image 的详细信息
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
