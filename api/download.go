package api

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ****************************************************  单文件  *****************************************************
// 下载一个文件
func SingleDownloader(w http.ResponseWriter, r *http.Request) {
	fileArrs := strings.Split(r.URL.Path, "/")
	filename := fileArrs[len(fileArrs)-1]

	// 打印文件名
	fmt.Println("Download: ", filename)

	// 打开文件
	file, err := os.Open(filepath + "/" + filename)
	if err != nil {
		http.Error(w, "Error opening the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// 将文件内容写入响应体
	io.Copy(w, file)
}

// 下载一个结果文件
func SingleResultDownloader(w http.ResponseWriter, r *http.Request) {
	fileArrs := strings.Split(r.URL.Path, "/")
	filename := fileArrs[len(fileArrs)-1]

	// 打印文件名
	fmt.Println("Download: ", filename)

	// 打开文件
	file, err := os.Open(resultpath + "/" + filename)
	if err != nil {
		http.Error(w, "Error opening the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// 将文件内容写入响应体
	io.Copy(w, file)
}

// ****************************************************  多文件  *****************************************************
// 下载多个文件，打包成zip文件
func MultiDownloader(w http.ResponseWriter, r *http.Request) {
	Cors(w)
	method := r.Method

	// 如果是Post请求，下载多个文件
	if method == http.MethodPost {
		var requestData struct {
			Files []string `json:"fileNames"`
		}

		// 解析请求体
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 打印要下载的文件列表
		fmt.Println("Download files: ", requestData.Files)

		// 创建一个zip文件
		zipPath := filepath + "/download.zip"
		zipFile, err := os.Create(zipPath)
		if err != nil {
			http.Error(w, "Error creating the zip file", http.StatusInternalServerError)
			return
		}
		defer zipFile.Close()

		// 创建一个zip.Writer
		zipWriter := zip.NewWriter(zipFile)

		// 对于每个文件，打开文件，创建zip文件，将文件内容写入zip文件
		for _, file := range requestData.Files {
			// 打开文件
			srcFile, err := os.Open(filepath + "/" + file)
			if err != nil {
				http.Error(w, "Error opening the file: "+file, http.StatusInternalServerError)
				return
			}

			// 创建zip文件条目
			zipEntry, err := zipWriter.Create(file)
			if err != nil {
				http.Error(w, "Error creating the zip entry for file: "+file, http.StatusInternalServerError)
				return
			}

			// 将文件内容写入zip条目
			_, err = io.Copy(zipEntry, srcFile)
			srcFile.Close() // 在每次循环结束时关闭文件
			if err != nil {
				http.Error(w, "Error writing the file: "+file, http.StatusInternalServerError)
				return
			}
		}

		// 关闭zip.Writer以完成写入
		err = zipWriter.Close()
		if err != nil {
			http.Error(w, "Error closing the zip writer", http.StatusInternalServerError)
			return
		}

		// 打开已完成的zip文件进行读取
		zipFile, err = os.Open(zipPath)
		if err != nil {
			http.Error(w, "Error opening the zip file for reading", http.StatusInternalServerError)
			return
		}
		defer zipFile.Close()
		defer os.Remove(zipPath)

		// 设置响应头
		w.Header().Set("Content-Disposition", "attachment; filename=download.zip")
		w.Header().Set("Content-Type", "application/zip")

		// 将zip文件内容写入响应体
		_, err = io.Copy(w, zipFile)
		if err != nil {
			http.Error(w, "Error sending the zip file", http.StatusInternalServerError)
		}
	}
}
