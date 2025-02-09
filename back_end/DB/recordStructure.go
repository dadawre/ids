package DB

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"graduationDesign/back_end/Initialize/Parameter"
	"net/http"
	"strconv"
)

// ActionDescribeTriggerLog 数据包匹配规则结构体
type ActionDescribeTriggerLog struct {
	ID              int    `json:"ID" gorm:"primaryKey" gorm:"column:id"`
	SrcIPAddr       string `json:"SrcIPAddr" gorm:"column:src_ip_addr"`
	DstIPAddr       string `json:"DstIPAddr" gorm:"column:dst_ip_addr"`
	SrcPortAddr     string `json:"SrcPortAddr" gorm:"column:src_port_addr"`
	DstPortAddr     string `json:"DstPortAddr" gorm:"column:dst_port_addr"`
	TopProtocol     string `json:"TopProtocol" gorm:"column:top_protocol"`
	RuleID          int    `json:"RuleID" gorm:"column:rule_id"`
	RuleName        string `json:"RuleName" gorm:"column:rule_name"`
	DescribeMessage string `json:"DescribeMessage" gorm:"column:describe_message"`
	TriggerTime     string `json:"TriggerTime" gorm:"column:trigger_time"`
	ActionType      int8   `json:"ActionType" gorm:"column:action_type"`
	TcpFlag         string `json:"TcpFlag" gorm:"column:tcp_flag"`
	PacketLength    int    `json:"PacketLength" gorm:"column:packet_length"`
}

// QueryParams 查询参数结构体
type QueryParams struct {
	Page     int
	PageSize int
	Filters  map[string]string
}

// RecordSQLManager 管理 MySQL 连接和操作
type RecordSQLManager struct {
	DB *gorm.DB
}

// ConnectAddrStr 连接到 MySQL 数据库
func ConnectAddrStr(user, password, host string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, host)
}

// NewMySQLManager 连接 MySQL 并返回管理对象
func NewMySQLManager(dsn string) (*RecordSQLManager, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}
	return &RecordSQLManager{DB: db}, nil
}

// EnsureDatabase 存在某架构，不存在则创建
func (m *RecordSQLManager) EnsureDatabase(SchemaName string) error {
	// 保护数据库名，避免 SQL 注入
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", SchemaName)
	return m.DB.Exec(sql).Error
}

// TableExists 检测表是否存在
func (m *RecordSQLManager) TableExists(schemaName string, tableName string) (bool, error) {
	var count int64
	err := m.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", schemaName, tableName).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateTable 创建存储 ActionDescribeTriggerLog 的表（如果不存在）
func (m *RecordSQLManager) CreateTable(tableName string) error {
	return m.DB.Table(tableName).AutoMigrate(&ActionDescribeTriggerLog{})
}

// InsertRecord 插入数据
func (m *RecordSQLManager) InsertRecord(table string, log *ActionDescribeTriggerLog) error {
	// 每个插入操作使用独立事务
	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(table).Create(log).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// DeleteRecord 删除符合条件的数据
func (m *RecordSQLManager) DeleteRecord(tableName string, condition map[string]interface{}) error {
	return m.DB.Table(tableName).Where(condition).Delete(&ActionDescribeTriggerLog{}).Error
}

// UpdateRecord 修改数据
func (m *RecordSQLManager) UpdateRecord(tableName string, condition map[string]interface{}, updates map[string]interface{}) error {
	return m.DB.Table(tableName).Where(condition).Updates(updates).Error
}

// GetLogs 整合后的HTTP处理函数
func (m *RecordSQLManager) GetLogs(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	params := parseQueryParams(r)
	fmt.Println(params)
	// 执行查询
	logs, total, err := m.QueryWithPagination(Parameter.SqlRecordTableName, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"items": logs,
		"total": total,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "响应编码失败", http.StatusInternalServerError)
	}
}

// QueryWithPagination 支持模糊查询和分页的查询方法
func (m *RecordSQLManager) QueryWithPagination(tableName string, params QueryParams) ([]ActionDescribeTriggerLog, int64, error) {
	var results []ActionDescribeTriggerLog
	var count int64

	// **每次查询都创建新的会话，避免污染**
	query := m.DB.Session(&gorm.Session{NewDB: true}).Table(tableName)

	// 应用过滤条件
	query = applyFilters(query, params.Filters)

	// 计算总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %v", err)
	}

	// 应用分页
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// 执行查询
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("查询数据失败: %v", err)
	}

	return results, count, nil
}

// 前端接收键转SQL键
var sqlColumnsMap = map[string]string{
	"id":               "ID",
	"src_ip_addr":      "SrcIPAddr",
	"dst_ip_addr":      "DstIPAddr",
	"src_port_addr":    "SrcPortAddr",
	"dst_port_addr":    "DstPortAddr",
	"top_protocol":     "TopProtocol",
	"rule_id":          "RuleID",
	"rule_name":        "RuleName",
	"describe_message": "DescribeMessage",
	"trigger_time":     "TriggerTime",
	"action_type":      "ActionType",
	"tcp_flag":         "TcpFlag",
	"packet_length":    "PacketLength",
}

// 解析查询参数
func parseQueryParams(r *http.Request) QueryParams {
	query := r.URL.Query()

	// 分页参数
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))

	// 过滤参数
	filters := make(map[string]string)
	for key, field := range sqlColumnsMap {
		filters[key] = query.Get(field)
	}

	return QueryParams{
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
}

// 应用过滤条件
func applyFilters(query *gorm.DB, filters map[string]string) *gorm.DB {
	for field, value := range filters {
		if value == "" {
			continue
		}

		// 特殊处理数字字段
		switch field {
		case "ID", "RuleID", "PacketLength", "ActionType":
			if num, err := strconv.Atoi(value); err == nil {
				query = query.Where(fmt.Sprintf("%s = ?", field), num)
			}
		default:
			// 字符串字段使用LIKE模糊查询
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
		}
	}
	return query
}

//// QueryRecords 查询数据，支持分页和计数
//func (m *RecordSQLManager) QueryRecords(tableName string, condition map[string]interface{}, limit int, offset int) ([]ActionDescribeTriggerLog, int64, error) {
//	var results []ActionDescribeTriggerLog
//	var count int64
//
//	// 先计算符合条件的总数
//	err := m.DB.Table(tableName).Where(condition).Count(&count).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	// 查询符合条件的数据
//	err = m.DB.Table(tableName).Where(condition).Limit(limit).Offset(offset).Find(&results).Error
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return results, count, nil
//}

//func main() {
//	// MySQL 连接信息 (请替换成你的 MySQL 配置)
//	dsn := ConnectAddrStr(Parameter.SqlUser, Parameter.SqlPassword, Parameter.SqlHost)
//
//	// 创建 MySQL 管理实例
//	manager, err := NewMySQLManager(dsn)
//	if err != nil {
//		log.Fatalf("数据库连接失败: %v", err)
//	}
//
//	// 确保数据库存在
//	if err := manager.EnsureDatabase(Parameter.SqlSchemaName); err != nil {
//		log.Fatalf("创建数据库失败: %v", err)
//	}
//
//	// 切换到 `IDSdata` 数据库
//	manager.DB = manager.DB.Exec("USE IDSdata")
//	fullPathTableName := Parameter.SqlSchemaName + "." + Parameter.SqlRecordTableName
//
//	// 检查表是否存在
//	exists, err := manager.TableExists(Parameter.SqlSchemaName, Parameter.SqlRecordTableName)
//	if err != nil {
//		log.Fatalf("检查表失败: %v", err)
//	}
//	if !exists {
//		log.Println("表不存在，正在创建...")
//		if err := manager.CreateTable(fullPathTableName); err != nil {
//			log.Fatalf("创建表失败: %v", err)
//		}
//	}
//
//	// 示例: 插入数据
//	record := &ActionDescribeTriggerLog{
//		ID:              1,
//		SrcIPAddr:       "192.168.1.1",
//		DstIPAddr:       "10.0.0.1",
//		SrcPortAddr:     "12345",
//		DstPortAddr:     "80",
//		TopProtocol:     "TCP",
//		RuleID:          1001,
//		RuleName:        "SQL Injection",
//		DescribeMessage: "Detected SQL injection attack",
//		TriggerTime:     "2025-01-31 12:34:56",
//		ActionType:      1,
//		TcpFlag:         "SYN",
//		PacketLength:    512,
//	}
//
//	if err := manager.InsertRecord(fullPathTableName, record); err != nil {
//		log.Fatalf("插入数据失败: %v", err)
//	}
//
//	// 示例: 查询数据
//	condition := map[string]interface{}{"SrcIPAddr": "192.168.1.1"}
//	records, total, err := manager.QueryRecords(fullPathTableName, condition, 10, 0)
//	if err != nil {
//		log.Fatalf("查询数据失败: %v", err)
//	}
//	log.Printf("查询到 %d 条符合条件的数据", total)
//	for _, r := range records {
//		log.Printf("ID: %d, SrcIP: %s, DstIP: %s, Rule: %s", r.ID, r.SrcIPAddr, r.DstIPAddr, r.RuleName)
//	}
//
//	// 示例: 更新数据
//	updateData := map[string]interface{}{"RuleName": "Updated Rule"}
//	if err := manager.UpdateRecord(fullPathTableName, condition, updateData); err != nil {
//		log.Fatalf("更新数据失败: %v", err)
//	}
//
//	// 示例: 删除数据
//	if err := manager.DeleteRecord(fullPathTableName, condition); err != nil {
//		log.Fatalf("删除数据失败: %v", err)
//	}
//}
