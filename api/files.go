package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// get /api/files 是获取所有文件的列表
// post /api/upload 是上传文件
// get /api/files/:filename 是下载一个文件
// delete /api/files/:filename 是删除一个文件
// get /api/results 是获取所有结果的列表
// get /api/results/:resultName 是下载一个结果
// delete /api/results/:resultName 是删除一个结果
// delete /api/files 是批量删除文件
// delete /api/results 是批量删除结果
// get /api/files/download 是下载多个文件

// cors 跨域请求
func Cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 定义上传文件和结果文件路径
const filepath string = "./uploads"
const resultpath string = "./results"

// 转换文件大小为人类可读的格式
func getSize(size int64) string {
	var sizeStr string
	if size < 1024 {
		sizeStr = fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		sizeStr = fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	} else {
		sizeStr = fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
	}
	return sizeStr
}

// 获取所有文件的列表，或者批量删除文件
func FilesHandler(w http.ResponseWriter, r *http.Request) {
	// 跨域请求
	Cors(w)
	method := r.Method

	// 如果是GET请求，获取所有文件列表
	if method == http.MethodGet {
		// 如果files文件夹不存在，创建一个
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			os.Mkdir(filepath, os.ModePerm)
		}
		// 如果files文件夹存在，获取文件列表，以数组形式返回
		files, _ := os.ReadDir(filepath)
		filesArray := make([]string, 0)
		for _, file := range files {
			// 忽略.gitkeep文件 __MACOSX文件夹 .DS_Store文件
			if file.Name() != ".gitkeep" && file.Name() != "__MACOSX" && file.Name() != ".DS_Store" {
				filesArray = append(filesArray, file.Name())
			}
		}
		// 返回文件列表
		json.NewEncoder(w).Encode(filesArray)
	}

	// 如果是DELETE请求，删除请求体中的文件列表
	if method == http.MethodDelete {
		MultiDeleter(w, r)
	}
}

// 获取所有结果的列表
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	// 跨域请求
	Cors(w)
	method := r.Method

	// 如果是GET请求，获取所有结果文件列表
	if method == http.MethodGet {
		// 如果results文件夹不存在，创建一个
		if _, err := os.Stat(resultpath); os.IsNotExist(err) {
			os.Mkdir(resultpath, os.ModePerm)
		}
		// 如果results文件夹存在，获取文件列表，以数组形式返回
		files, _ := os.ReadDir(resultpath)
		filesArray := make([]string, 0)
		for _, file := range files {
			// 忽略.gitkeep文件 __MACOSX文件夹 .DS_Store文件
			if file.Name() != ".gitkeep" && file.Name() != "__MACOSX" && file.Name() != ".DS_Store" {
				filesArray = append(filesArray, file.Name())
			}
		}
		// 返回文件列表
		json.NewEncoder(w).Encode(filesArray)
	}

	// 如果是DELETE请求，删除请求体中的结果文件
	if method == http.MethodDelete {
		MultiResultDeleter(w, r)
	}
}

// 下载一个文件或删除一个文件
func FileProcessor(w http.ResponseWriter, r *http.Request) {
	Cors(w)
	method := r.Method

	// 如果是GET请求，下载文件
	if method == http.MethodGet {
		SingleDownloader(w, r)
	}

	// 如果是DELETE请求，删除文件
	if method == http.MethodDelete {
		SingleDeleter(w, r)
	}
}

// 下载一个结果或删除一个结果
func ResultProcessor(w http.ResponseWriter, r *http.Request) {
	Cors(w)
	method := r.Method

	// 如果是GET请求，下载这个结果文件
	if method == http.MethodGet {
		SingleResultDownloader(w, r)
	}

	// 如果是DELETE请求，删除这个结果文件
	if method == http.MethodDelete {
		SingleResultDeleter(w, r)
	}
}
