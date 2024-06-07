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
	thisFile := filepath + "/" + filename

	// 如果是一个文件夹，删除文件夹
	fileInfo, err := os.Stat(thisFile)
	if os.IsNotExist(err) {
		http.Error(w, "File or folder not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file info: %v", err), http.StatusInternalServerError)
		return
	}

	// 如果是文件夹, 删除文件夹
	if fileInfo.IsDir() {
		err := os.RemoveAll(thisFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting the folder: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("Folder deleted successfully")
		return
	}

	// 删除文件
	err = os.Remove(thisFile)
	if err != nil {
		http.Error(w, "Error deleting the file", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	fmt.Println("Deleted: ", filename)
	json.NewEncoder(w).Encode("File deleted successfully")
}

// 删除单个结果文件
func SingleResultDeleter(w http.ResponseWriter, r *http.Request) {
	fileArrs := strings.Split(r.URL.Path, "/")
	filename := fileArrs[len(fileArrs)-1]
	thisFile := resultpath + "/" + filename

	fileInfo, err := os.Stat(thisFile)
	if os.IsNotExist(err) {
		http.Error(w, "File or folder not found", http.StatusNotFound)
		return
	}

	// 如果是文件夹, 删除文件夹
	if fileInfo.IsDir() {
		err := os.RemoveAll(thisFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting the folder: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("Folder deleted successfully")
		return
	}

	// 删除文件
	err = os.Remove(thisFile)
	if err != nil {
		http.Error(w, "Error deleting the result", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	fmt.Println("Deleted: ", filename)
	json.NewEncoder(w).Encode("Result deleted successfully")
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

	// 删除文件
	for _, file := range requestData.Files.FileNames {
		// 检查文件是否存在
		fileInfo, err := os.Stat(filepath + "/" + file)
		if os.IsNotExist(err) {
			http.Error(w, "File or folder not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Error getting file info: %v", err), http.StatusInternalServerError)
			return
		}

		// 如果是一个文件夹，删除文件夹
		if fileInfo.IsDir() {
			err := os.RemoveAll(filepath + "/" + file)
			if err != nil {
				http.Error(w, "Error deleting the folder", http.StatusInternalServerError)
				return
			}
			continue
		} else {
			err := os.Remove(filepath + "/" + file)
			if err != nil {
				http.Error(w, "Error deleting the file", http.StatusInternalServerError)
				return
			}
		}
	}

	// 返回成功信息
	fmt.Println("Deleted files: ", requestData.Files.FileNames)
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

	// 删除文件
	for _, file := range requestData.Files.FileNames {
		// 检查文件是否存在
		fileInfo, err := os.Stat(resultpath + "/" + file)
		if os.IsNotExist(err) {
			http.Error(w, "File or folder not found", http.StatusNotFound)
			return
		}

		// 如果是一个文件夹，删除文件夹
		if fileInfo.IsDir() {
			err := os.RemoveAll(resultpath + "/" + file)
			if err != nil {
				http.Error(w, "Error deleting the folder", http.StatusInternalServerError)
				return
			}
			continue
		} else {
			err := os.Remove(resultpath + "/" + file)
			if err != nil {
				http.Error(w, "Error deleting the result", http.StatusInternalServerError)
				return
			}
		}
	}

	// 返回成功信息
	fmt.Println("Deleted results: ", requestData.Files.FileNames)
	json.NewEncoder(w).Encode("Results deleted successfully")
}
