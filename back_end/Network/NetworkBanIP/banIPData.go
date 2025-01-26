package NetworkBanIP

import (
	"encoding/json"
	"fmt"
	"graduationDesign/back_end/Initialize"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// AddBanIP handles adding a BanIP record.
func AddBanIP(w http.ResponseWriter, r *http.Request) {
	// 检查请求的 Content-Type 是否为 application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type, expected application/json", http.StatusBadRequest)
		return
	}

	var newBanIP Initialize.BanIPDescribe

	// 读取请求体中的 JSON 数据
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newBanIP); err != nil {
		// 输出详细的错误信息
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 如果结构体的字段类型转换不成功（例如 banRuleID 从字符串转换为 int 失败）
	if newBanIP.BanRuleID == 0 && newBanIP.BanIP != "" {
		http.Error(w, "Invalid BanRuleID format", http.StatusBadRequest)
		return
	}

	// 插入数据
	success := Initialize.BanBanIPSet.Insert(newBanIP)
	if !success {
		http.Error(w, "Ban IP already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteBanIP handles deleting a BanIP record.
func DeleteBanIP(w http.ResponseWriter, r *http.Request) {
	banIP := r.URL.Query().Get("banIP")
	fmt.Println(banIP)
	if banIP == "" {
		http.Error(w, "Ban IP is required", http.StatusBadRequest)
		return
	}
	success := Initialize.BanBanIPSet.Delete(banIP)
	if !success {
		http.Error(w, "Ban IP not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FetchBanIpData handles fetching and pagination of Ban IP records.
func FetchBanIpData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	allData := Initialize.BanBanIPSet.GetAll()
	var filteredData []Initialize.BanIPDescribe
	for _, desc := range allData {
		if query == "" || strings.Contains(desc.BanIP, query) || strings.Contains(desc.Describe, query) {
			filteredData = append(filteredData, desc)
		}
	}

	// 排序：根据 BanIP 进行排序（你可以根据需要排序其他字段）
	sort.Slice(filteredData, func(i, j int) bool {
		return filteredData[i].BanIP < filteredData[j].BanIP // 按 BanIP 字段升序排序
	})

	total := len(filteredData)
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pagedData := filteredData[start:end]
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]interface{}{"data": pagedData, "total": total})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
