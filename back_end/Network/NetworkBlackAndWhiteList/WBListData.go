package NetworkBlackAndWhiteList

import (
	"encoding/json"
	"fmt"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/List"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 定义互斥锁，用于保证并发安全
var (
	mu sync.Mutex
)

// AddHandler 用于处理添加数据的请求（POST请求），接收 JSON 格式的 addr、isWhite、valueType 参数
func AddHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Addr      string          `json:"addr"`
		IsWhite   bool            `json:"isWhite"`
		ValueType List.WBAddrType `json:"valueType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "无效的输入数据", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	InitializeStructure.BlackWhiteList.IDindex++
	newEntry := List.WBValueStructure{
		ID:        InitializeStructure.BlackWhiteList.IDindex,
		Addr:      input.Addr,
		IsWhite:   input.IsWhite,
		ValueType: input.ValueType,
	}
	InitializeStructure.BlackWhiteList.Addrs[input.Addr] = newEntry
	List.InsertIPset(input.Addr, input.ValueType)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEntry)
}

// DeleteHandler 用于处理删除数据的请求，通过 URL 中的 "id" 参数来指定删除的记录
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	fmt.Println(idStr)
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	var foundKey *List.WBValueStructure
	for _, entry := range InitializeStructure.BlackWhiteList.Addrs {
		if entry.ID == id {
			foundKey = &entry
			break
		}
	}

	if foundKey == nil {
		http.Error(w, "未找到对应记录", http.StatusNotFound)
		return
	}

	delete(InitializeStructure.BlackWhiteList.Addrs, foundKey.Addr)
	List.DeleteIPset(foundKey.Addr, foundKey.ValueType)
	w.WriteHeader(http.StatusOK)
}

// QueryHandler 用于处理查询请求（GET请求），支持分页、过滤（包括地址关键字过滤）及排序功能
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// 分页参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 过滤参数
	target := r.URL.Query().Get("target")         // 地址关键字过滤
	filterIsWhite := r.URL.Query().Get("isWhite") // 是否白名单过滤
	filterValueType := r.URL.Query().Get("valueType")

	// 排序参数
	sortField := r.URL.Query().Get("sortField")
	sortOrder := r.URL.Query().Get("sortOrder") // "ascend" 或 "descend"

	var filtered []List.WBValueStructure
	mu.Lock()
	for _, entry := range InitializeStructure.BlackWhiteList.Addrs {
		// 地址过滤（使用 target 参数）
		if target != "" && !strings.Contains(entry.Addr, target) {
			continue
		}
		// 是否白名单过滤
		if filterIsWhite != "" {
			isWhite, err := strconv.ParseBool(filterIsWhite)
			if err == nil && entry.IsWhite != isWhite {
				continue
			}
		}
		// 地址类型过滤
		if filterValueType != "" {
			vt, err := strconv.Atoi(filterValueType)
			if err == nil && int(entry.ValueType) != vt {
				continue
			}
		}
		filtered = append(filtered, entry)
	}
	mu.Unlock()

	// 根据排序参数对 filtered 切片进行排序
	if sortField != "" {
		sort.Slice(filtered, func(i, j int) bool {
			var less bool
			switch sortField {
			case "id":
				less = filtered[i].ID < filtered[j].ID
			case "addr":
				less = strings.Compare(filtered[i].Addr, filtered[j].Addr) < 0
			case "isWhite":
				// false 排在 true 前面
				if filtered[i].IsWhite == filtered[j].IsWhite {
					less = false
				} else {
					less = !filtered[i].IsWhite && filtered[j].IsWhite
				}
			case "valueType":
				less = int(filtered[i].ValueType) < int(filtered[j].ValueType)
			default:
				less = true
			}
			if sortOrder == "descend" {
				return !less
			}
			return less
		})
	}

	// 分页处理
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	paginated := filtered[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginated)
}
