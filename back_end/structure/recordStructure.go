package Structure

import (
	"fmt"
	"graduationDesign/back_end/Initialize/Parameter"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ActionDescribeTriggerLog 数据包匹配规则结构体
type ActionDescribeTriggerLog struct {
	ID              int    `json:"ID" gorm:"primaryKey"`
	SrcIPAddr       string `json:"SrcIPAddr"`
	DstIPAddr       string `json:"DstIPAddr"`
	SrcPortAddr     string `json:"SrcPortAddr"`
	DstPortAddr     string `json:"DstPortAddr"`
	TopProtocol     string `json:"TopProtocol"`
	RuleID          int    `json:"RuleID"`
	RuleName        string `json:"RuleName"`
	DescribeMessage string `json:"DescribeMessage"`
	TriggerTime     string `json:"TriggerTime"`
	ActionType      int8   `json:"ActionType"`
	TcpFlag         string `json:"TcpFlag"`
	PacketLength    int    `json:"PacketLength"`
}

// MySQLManager 管理 MySQL 连接和操作
type MySQLManager struct {
	DB *gorm.DB
}

// ConnectAddrStr 连接到 MySQL 数据库
func ConnectAddrStr(user, password, host string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, host)
}

// NewMySQLManager 连接 MySQL 并返回管理对象
func NewMySQLManager(dsn string) (*MySQLManager, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}
	return &MySQLManager{DB: db}, nil
}

// EnsureDatabase 存在某架构，不存在则创建
func (m *MySQLManager) EnsureDatabase(SchemaName string) error {
	// 保护数据库名，避免 SQL 注入
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", SchemaName)
	return m.DB.Exec(sql).Error
}

// TableExists 检测表是否存在
func (m *MySQLManager) TableExists(schemaName string, tableName string) (bool, error) {
	var count int64
	err := m.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", schemaName, tableName).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateTable 创建存储 ActionDescribeTriggerLog 的表（如果不存在）
func (m *MySQLManager) CreateTable(tableName string) error {
	return m.DB.Table(tableName).AutoMigrate(&ActionDescribeTriggerLog{})
}

// InsertRecord 插入数据
func (m *MySQLManager) InsertRecord(table string, log *ActionDescribeTriggerLog) error {
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
func (m *MySQLManager) DeleteRecord(tableName string, condition map[string]interface{}) error {
	return m.DB.Table(tableName).Where(condition).Delete(&ActionDescribeTriggerLog{}).Error
}

// UpdateRecord 修改数据
func (m *MySQLManager) UpdateRecord(tableName string, condition map[string]interface{}, updates map[string]interface{}) error {
	return m.DB.Table(tableName).Where(condition).Updates(updates).Error
}

// QueryRecords 查询数据，支持分页和计数
func (m *MySQLManager) QueryRecords(tableName string, condition map[string]interface{}, limit int, offset int) ([]ActionDescribeTriggerLog, int64, error) {
	var results []ActionDescribeTriggerLog
	var count int64

	// 先计算符合条件的总数
	err := m.DB.Table(tableName).Where(condition).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询符合条件的数据
	err = m.DB.Table(tableName).Where(condition).Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func main() {
	// MySQL 连接信息 (请替换成你的 MySQL 配置)
	dsn := ConnectAddrStr(Parameter.SqlUser, Parameter.SqlPassword, Parameter.SqlHost)

	// 创建 MySQL 管理实例
	manager, err := NewMySQLManager(dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 确保数据库存在
	if err := manager.EnsureDatabase(Parameter.SqlSchemaName); err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	// 切换到 `IDSdata` 数据库
	manager.DB = manager.DB.Exec("USE IDSdata")
	fullPathTableName := Parameter.SqlSchemaName + "." + Parameter.SqlRecordTableName

	// 检查表是否存在
	exists, err := manager.TableExists(Parameter.SqlSchemaName, Parameter.SqlRecordTableName)
	if err != nil {
		log.Fatalf("检查表失败: %v", err)
	}
	if !exists {
		log.Println("表不存在，正在创建...")
		if err := manager.CreateTable(fullPathTableName); err != nil {
			log.Fatalf("创建表失败: %v", err)
		}
	}

	// 示例: 插入数据
	record := &ActionDescribeTriggerLog{
		ID:              1,
		SrcIPAddr:       "192.168.1.1",
		DstIPAddr:       "10.0.0.1",
		SrcPortAddr:     "12345",
		DstPortAddr:     "80",
		TopProtocol:     "TCP",
		RuleID:          1001,
		RuleName:        "SQL Injection",
		DescribeMessage: "Detected SQL injection attack",
		TriggerTime:     "2025-01-31 12:34:56",
		ActionType:      1,
		TcpFlag:         "SYN",
		PacketLength:    512,
	}

	if err := manager.InsertRecord(fullPathTableName, record); err != nil {
		log.Fatalf("插入数据失败: %v", err)
	}

	// 示例: 查询数据
	condition := map[string]interface{}{"SrcIPAddr": "192.168.1.1"}
	records, total, err := manager.QueryRecords(fullPathTableName, condition, 10, 0)
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	log.Printf("查询到 %d 条符合条件的数据", total)
	for _, r := range records {
		log.Printf("ID: %s, SrcIP: %s, DstIP: %s, Rule: %s", r.ID, r.SrcIPAddr, r.DstIPAddr, r.RuleName)
	}

	// 示例: 更新数据
	updateData := map[string]interface{}{"RuleName": "Updated Rule"}
	if err := manager.UpdateRecord(fullPathTableName, condition, updateData); err != nil {
		log.Fatalf("更新数据失败: %v", err)
	}

	// 示例: 删除数据
	if err := manager.DeleteRecord(fullPathTableName, condition); err != nil {
		log.Fatalf("删除数据失败: %v", err)
	}
}
