package NetworkBanIP

import (
	"encoding/json"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Structure"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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

	allData := InitializeStructure.MainBanIPSet.GetAllTransformArray()
	var filteredData []Structure.BanIPDescribe
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
