package Network

import (
	"github.com/gorilla/mux"
	"graduationDesign/back_end/Initialize/InitializeStructure"
	"graduationDesign/back_end/Network/NetworkBanIP"
	"graduationDesign/back_end/Network/NetworkBlackAndWhiteList"
	"graduationDesign/back_end/Network/NetworkConfig"
	"graduationDesign/back_end/Network/NetworkRule"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func InitializeNetwork() {
	NetworkConfig.InitializeConfigData()
	NewStartWebRoute()
}

// CorsMiddleware CORS 跨域response设置
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (change "*" to specific origin if needed)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewStartWebRoute() *mux.Router {
	rand.Seed(time.Now().UnixNano())

	// Create a new router
	router := mux.NewRouter()

	// Define the /nameData endpoint
	router.HandleFunc("/configData", NetworkConfig.HandleConfigData)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/changeConfig", NetworkConfig.ChangeConfigHandler)
	router.HandleFunc("/RuleAdd", NetworkRule.RuleAdd)
	router.HandleFunc("/RuleDelete", NetworkRule.RuleDelete)
	router.HandleFunc("/RuleEdit", NetworkRule.RuleEdit)
	router.HandleFunc("/RuleFetch", NetworkRule.RuleFetch)
	router.HandleFunc("/logData", InitializeStructure.LogSQLManager.GetLogs)
	router.HandleFunc("/banIPData", NetworkBanIP.FetchBanIpData).Methods("GET")
	router.HandleFunc("/DeleteBanIP", NetworkBanIP.DeleteBanIP)
	router.HandleFunc("/AddBanIP", NetworkBanIP.AddBanIP)

	router.HandleFunc("/AddWBData", NetworkBlackAndWhiteList.AddHandler)
	router.HandleFunc("/DeleteWBData", NetworkBlackAndWhiteList.DeleteHandler)
	router.HandleFunc("/QueryWBData", NetworkBlackAndWhiteList.QueryHandler)

	router.HandleFunc("/protected", authMiddleware(protectedHandler)) // 使用中间件保护路由
	// Apply the CORS middleware
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", CorsMiddleware(router)))
	return router
}
