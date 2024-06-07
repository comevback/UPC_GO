package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// 上传单个或多个文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	Cors(w)

	// 如果是OPTIONS预检请求，返回200
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 必须是POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 上传的目标路径
	targetPath := filepath
	// 如果文件夹不存在，创建一个
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		os.Mkdir(targetPath, os.ModePerm)
	}

	// 解析请求
	err := r.ParseMultipartForm(100 << 20) // 100MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}
	// 清理暂存文件
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["file"]
	var uploadedFiles []string

	// 对于接收到的每个文件，提取文件名，创造目标文件路径，通过io.Copy()函数将文件内容写入目标文件
	for _, fileHeader := range files {
		// 打印文件信息，包括文件名和文件大小，把文件大小转化为人类可读的格式
		size := fileHeader.Size
		sizeStr := getSize(size)
		fmt.Printf("Uploaded: %s --- Size: %s\n", fileHeader.Filename, sizeStr)

		// 打开这个文件
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 创建目标文件
		dstPath := targetPath + "/" + fileHeader.Filename
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Error creating the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// 将文件内容写入目标文件
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}

		// 记录上传的文件名到uploadedFiles数组
		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}
	//fmt.Println("Uploaded files: ", uploadedFiles)
	json.NewEncoder(w).Encode(uploadedFiles)
}
