package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ****************************************************  单文件  *****************************************************
// 删除单个文件
func SingleDeleter(w http.ResponseWriter, r *http.Request) {
	fileArrs := strings.Split(r.URL.Path, "/")
	filename := fileArrs[len(fileArrs)-1]

	// 删除文件
	err := os.Remove(filepath + "/" + filename)
	if err != nil {
		http.Error(w, "Error deleting the file", http.StatusInternalServerError)
		return
	}
}

// 删除单个结果文件
func SingleResultDeleter(w http.ResponseWriter, r *http.Request) {
	fileArrs := strings.Split(r.URL.Path, "/")
	filename := fileArrs[len(fileArrs)-1]

	// 删除文件
	err := os.Remove(resultpath + "/" + filename)
	if err != nil {
		http.Error(w, "Error deleting the file", http.StatusInternalServerError)
		return
	}
}

// ****************************************************  多文件  *****************************************************
// 批量删除文件
func MultiDeleter(w http.ResponseWriter, r *http.Request) {
	// 从body中获取要删除的文件列表
	// 解析请求体
	var requestData struct {
		Files struct {
			FileNames []string `json:"fileNames"`
		} `json:"files"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 打印要删除的文件列表
	fmt.Println("Delete files: ", requestData.Files.FileNames)

	// 删除文件
	for _, file := range requestData.Files.FileNames {
		err := os.Remove(filepath + "/" + file)
		if err != nil {
			http.Error(w, "Error deleting the file", http.StatusInternalServerError)
			return
		}
	}

	// 返回成功信息
	json.NewEncoder(w).Encode("Files deleted successfully")
}

// 批量删除结果文件
func MultiResultDeleter(w http.ResponseWriter, r *http.Request) {
	// 从body中获取要删除的结果文件
	// 解析请求体
	var requestData struct {
		Files struct {
			FileNames []string `json:"fileNames"`
		} `json:"files"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 打印要删除的结果文件
	fmt.Println("Delete results: ", requestData.Files.FileNames)

	// 删除文件
	for _, file := range requestData.Files.FileNames {
		err := os.Remove(resultpath + "/" + file)
		if err != nil {
			http.Error(w, "Error deleting the file", http.StatusInternalServerError)
			return
		}
	}

	// 返回成功信息
	json.NewEncoder(w).Encode("Results deleted successfully")
}
