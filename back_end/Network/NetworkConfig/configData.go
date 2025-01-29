package NetworkConfig

import (
	"encoding/json"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"net/http"
	"sort"
	"strconv"
)

type Config struct {
	Index int         `json:"index"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ConfigData 存储所有的从Initialize中获得的Config
var ConfigData []Config = make([]Config, 0)

func GetConfigData() *[]Config {
	return &ConfigData
}

// InitializeConfigData 从Intialize中获取所有的Config
func InitializeConfigData() {
	if len(ConfigData) == 0 {
		i := 0
		for key, value := range InitializeStructure.MapVariableNameToVariable {
			ConfigData = append(ConfigData, Config{Index: i, Key: key, Value: value})
			i++
		}
	}
}

// HandleConfigData 处理/nameData
func HandleConfigData(w http.ResponseWriter, r *http.Request) {
	//解析前端键值对
	query := r.URL.Query()
	//设置每一页大小
	pageSize, _ := strconv.Atoi(query.Get("results"))
	if pageSize == 0 {
		pageSize = 10
	}
	//获得当前页号
	currentPage, _ := strconv.Atoi(query.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}

	//获取排序字段和排序方式
	sortField := query.Get("sortField")
	sortOrder := query.Get("sortOrder")

	//fmt.Println(" sortField:", sortField, " sortOrder:", sortOrder)

	// 过滤数据
	filteredData := ConfigData

	// 数据排序
	if sortField != "" {
		sort.Slice(filteredData, func(i, j int) bool {
			if sortField == "index" {
				if sortOrder == "ascend" {
					return filteredData[i].Index < filteredData[j].Index
				} else if sortOrder == "descend" {
					return filteredData[i].Index > filteredData[j].Index
				}
			}
			return false
		})
	}

	//分页数据
	startIndex := (currentPage - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex > len(filteredData) {
		startIndex = len(filteredData)
	}
	if endIndex > len(filteredData) {
		endIndex = len(filteredData)
	}
	paginatedData := filteredData[startIndex:endIndex]
	//fmt.Println("filterdDataLen:", len(filteredData))
	//fmt.Println("filterdData:", filteredData)

	// Prepare response
	response := struct {
		Results []Config `json:"results"`
		Total   int      `json:"total"`
	}{
		Results: paginatedData,
		Total:   len(filteredData),
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
