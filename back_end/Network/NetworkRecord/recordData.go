package NetworkRecord

import (
	"encoding/json"
	"graduationDesign/back_end/Structure"
	"net/http"
	"strconv"
)

// 处理前端传来的筛选条件
func filterLogs(filter map[string]string) []Structure.ActionDescribeTriggerLog {

	var filteredLogs []Structure.ActionDescribeTriggerLog
	//for _, log := range InitializeStructure.ActionDescribeTriggerLogArray {
	//	match := true
	//	// 根据前端传来的字段进行过滤
	//	if filter["DescribeMessage"] != "" && !strings.Contains(log.DescribeMessage, filter["DescribeMessage"]) {
	//		match = false
	//	}
	//	if filter["SrcIPAddr"] != "" && !strings.Contains(log.SrcIPAddr, filter["SrcIPAddr"]) {
	//		match = false
	//	}
	//	if filter["DstIPAddr"] != "" && !strings.Contains(log.DstIPAddr, filter["DstIPAddr"]) {
	//		match = false
	//	}
	//	if filter["ActionType"] != "" && fmt.Sprintf("%d", log.ActionType) != filter["ActionType"] {
	//		match = false
	//	}
	//	if match {
	//		filteredLogs = append(filteredLogs, log)
	//	}
	//}
	return filteredLogs
}

// GetLogs 获取日志数据，支持前端发送的筛选条件
func GetLogs(w http.ResponseWriter, r *http.Request) {
	// 解析前端传来的查询参数（筛选条件）
	filter := make(map[string]string)
	query := r.URL.Query()

	// 从查询参数中读取筛选条件
	filter["DescribeMessage"] = query.Get("DescribeMessage")
	filter["SrcIPAddr"] = query.Get("SrcIPAddr")
	filter["DstIPAddr"] = query.Get("DstIPAddr")
	filter["ActionType"] = query.Get("ActionType")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	// 过滤数据
	filteredLogs := filterLogs(filter)

	// 计算分页
	start := (page - 1) * pageSize
	end := page * pageSize
	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}
	paginatedLogs := filteredLogs[start:end]

	// 返回数据
	response := map[string]interface{}{
		"items": paginatedLogs,
		"total": len(filteredLogs),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
