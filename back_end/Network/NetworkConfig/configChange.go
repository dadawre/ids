package NetworkConfig

import (
	"encoding/json"
	"fmt"
	"graduationDesign/back_end/Initialize"
	"io"
	"net/http"
)

type ConfigRequest struct {
	Index int         `json:"index"`
	Value interface{} `json:"value"`
}

type ConfigResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ChangeConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 读取请求体数据
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取请求体失败:", err)
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 打印请求体内容，帮助调试
	fmt.Println("请求体内容:", string(body))

	// 解析 JSON 请求体
	var req ConfigRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println("解析失败:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("可以修改:", req)
	ConfigData[req.Index].Value = req.Value
	Initialize.MapVariableNameToVariable[ConfigData[req.Index].Key] = req.Value
	fmt.Println(Initialize.MapVariableNameToVariable[ConfigData[req.Index].Key])
	// 模拟配置修改成功
	res := ConfigResponse{Success: true, Message: "配置修改成功"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
