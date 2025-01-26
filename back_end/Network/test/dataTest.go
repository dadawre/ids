package test

//
//import (
//	"encoding/json"
//	"fmt"
//	"graduationDesign/back_end/Network"
//	"log"
//	"math/rand"
//	"net/http"
//	"sort"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/gorilla/mux"
//)
//
//// User represents a single user entry
//type User struct {
//	Name struct {
//		First string `json:"first"`
//		Last  string `json:"last"`
//	} `json:"name"`
//	Gender string `json:"gender"`
//	Email  string `json:"email"`
//	Login  struct {
//		UUID string `json:"uuid"`
//	} `json:"login"`
//}
//
//// GenerateMockData generates a list of mock users
//func GenerateMockData(count int) []User {
//	firstNames := []string{"John", "Jane", "Michael", "Emily", "Chris", "Sarah"}
//	lastNames := []string{"Doe", "Smith", "Brown", "Johnson", "Davis", "Martinez"}
//	genders := []string{"male", "female"}
//
//	users := make([]User, count)
//	for i := 0; i < count; i++ {
//		users[i].Name.First = firstNames[rand.Intn(len(firstNames))]
//		users[i].Name.Last = lastNames[rand.Intn(len(lastNames))]
//		users[i].Gender = genders[rand.Intn(len(genders))]
//		users[i].Email = strings.ToLower(users[i].Name.First + "." + users[i].Name.Last + "@example.com")
//		users[i].Login.UUID = strconv.Itoa(rand.Intn(1000000))
//	}
//	return users
//}
//
//// HandleNameData handles the /nameData endpoint
//func HandleNameData(w http.ResponseWriter, r *http.Request) {
//	// Parse query parameters
//	query := r.URL.Query()
//	//设置每一页大小
//	pageSize, _ := strconv.Atoi(query.Get("results"))
//	if pageSize == 0 {
//		pageSize = 10
//	}
//	//获得当前页号
//	currentPage, _ := strconv.Atoi(query.Get("page"))
//	if currentPage == 0 {
//		currentPage = 1
//	}
//	//filterStrs := query.Get("filters[gender][0]")
//	//获取过滤关键字
//	var filterStrs []string = make([]string, 0)
//	for i, filterStr := 1, query.Get(fmt.Sprintf("filters[gender][%d]", 0)); filterStr != ""; i, filterStr = i+1, query.Get(fmt.Sprintf("filters[gender][%d]", i)) {
//		filterStrs = append(filterStrs, filterStr)
//	}
//
//	//获取排序字段和排序方式
//	sortField := query.Get("sortField")
//	sortOrder := query.Get("sortOrder")
//
//	fmt.Println("filterStrs: ", filterStrs, " sortField:", sortField, " sortOrder:", sortOrder)
//	fmt.Println(len(filterStrs), "  ", filterStrs)
//
//	//生成模拟数据
//	allData := GenerateMockData(200)
//
//	// 过滤数据
//	var filteredData []User
//	if len(filterStrs) > 0 {
//		for _, user := range allData {
//			for _, filterStr := range filterStrs {
//				if strings.EqualFold(user.Gender, filterStr) {
//					filteredData = append(filteredData, user)
//				}
//			}
//
//		}
//	} else {
//		filteredData = allData
//	}
//
//	// 数据排序
//	if sortField != "" {
//		sort.Slice(filteredData, func(i, j int) bool {
//			if sortField == "Name.First" {
//				if sortOrder == "ascend" {
//					return strings.ToLower(filteredData[i].Name.First) < strings.ToLower(filteredData[j].Name.First)
//				} else if sortOrder == "descend" {
//					return strings.ToLower(filteredData[i].Name.First) > strings.ToLower(filteredData[j].Name.First)
//				}
//			}
//			return false
//		})
//	}
//
//	//分页数据
//	startIndex := (currentPage - 1) * pageSize
//	endIndex := startIndex + pageSize
//	if startIndex > len(filteredData) {
//		startIndex = len(filteredData)
//	}
//	if endIndex > len(filteredData) {
//		endIndex = len(filteredData)
//	}
//	paginatedData := filteredData[startIndex:endIndex]
//	fmt.Println("filterdDataLen:", len(filteredData))
//	// Prepare response
//	response := struct {
//		Results []User `json:"results"`
//		Total   int    `json:"total"`
//	}{
//		Results: paginatedData,
//		Total:   len(filteredData),
//	}
//
//	// Return JSON response
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(response)
//}
//
//
//func main() {
//	rand.Seed(time.Now().UnixNano())
//
//	// Create a new router
//	router := mux.NewRouter()
//
//	// Define the /nameData endpoint
//	router.HandleFunc("/nameData", HandleNameData).Methods("GET")
//	router.HandleFunc("/login", network.loginHandler)
//	router.HandleFunc("/protected", network.authMiddleware(network.protectedHandler)) // 使用中间件保护路由
//	// Apply the CORS middleware
//	log.Println("Server running on http://localhost:8080")
//	log.Fatal(http.ListenAndServe(":8080", CorsMiddleware(router)))
//
//}
