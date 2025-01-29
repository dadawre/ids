package NetworkBanIP

import (
	"encoding/json"
	"fmt"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Structure"
	"net/http"
)

// AddBanIP handles adding a BanIP record.
func AddBanIP(w http.ResponseWriter, r *http.Request) {
	// 检查请求的 Content-Type 是否为 application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type, expected application/json", http.StatusBadRequest)
		return
	}

	var newBanIP Structure.BanIPDescribe

	// 读取请求体中的 JSON 数据
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newBanIP); err != nil {
		// 输出详细的错误信息
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	newBanIP.BanRuleID = -1
	fmt.Println(newBanIP)

	// 如果结构体的字段类型转换不成功（例如 banRuleID 从字符串转换为 int 失败）
	if newBanIP.BanIP == "" {
		http.Error(w, "Invalid BanRuleID format", http.StatusBadRequest)
		return
	}

	// 插入数据
	success := InitializeStructure.MainBanIPSet.Insert(newBanIP)
	if !success {
		http.Error(w, "Ban IP already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteBanIP handles deleting a BanIP record.
func DeleteBanIP(w http.ResponseWriter, r *http.Request) {
	banIP := r.URL.Query().Get("banIP")
	//fmt.Println(banIP)
	if banIP == "" {
		http.Error(w, "Ban IP is required", http.StatusBadRequest)
		return
	}

	success := InitializeStructure.MainBanIPSet.Delete(banIP)
	if !success {
		http.Error(w, "Ban IP not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
