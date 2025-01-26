package NetworkRule

import (
	"encoding/json"
	"graduationDesign/back_end/Initialize"
	"graduationDesign/back_end/structure/ruleStructure"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

var mu sync.Mutex

func RuleFetch(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(Initialize.ArrayRuleStrings)
	// 解析前端键值对
	query := r.URL.Query()
	pageSize, _ := strconv.Atoi(query.Get("results"))
	if pageSize == 0 {
		pageSize = 10 // 默认每页大小
	}

	currentPage, _ := strconv.Atoi(query.Get("page"))
	if currentPage == 0 {
		currentPage = 1 // 默认当前页
	}

	sortField := query.Get("sortField")
	sortOrder := query.Get("sortOrder")

	// 数据过滤
	mu.Lock()
	filteredData := Initialize.ArrayRuleStrings
	mu.Unlock()

	// 数据排序
	if sortField != "" {
		sort.Slice(filteredData, func(i, j int) bool {
			switch sortField {
			case "ID":
				if sortOrder == "ascend" {
					return filteredData[i].ID < filteredData[j].ID
				} else if sortOrder == "descend" {
					return filteredData[i].ID > filteredData[j].ID
				}
			case "Priority":
				if sortOrder == "ascend" {
					return filteredData[i].Priority < filteredData[j].Priority
				} else if sortOrder == "descend" {
					return filteredData[i].Priority > filteredData[j].Priority
				}
			case "ProtocolsString":
				if sortOrder == "ascend" {
					return filteredData[i].ProtocolsString < filteredData[j].ProtocolsString
				} else if sortOrder == "descend" {
					return filteredData[i].ProtocolsString > filteredData[j].ProtocolsString
				}
			}
			return false
		})
	}

	// 分页
	startIndex := (currentPage - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex > len(filteredData) {
		startIndex = len(filteredData)
	}
	if endIndex > len(filteredData) {
		endIndex = len(filteredData)
	}
	paginatedData := filteredData[startIndex:endIndex]

	// 准备响应数据
	response := struct {
		Results []ruleStructure.RuleNodeStringValue `json:"results"`
		Total   int                                 `json:"total"`
	}{
		Results: paginatedData,
		Total:   len(filteredData),
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
