package main

//
//import (
//	"graduationDesign/back_end/DB"
//	"graduationDesign/back_end/Initialize/Parameter"
//	"log"
//	"testing"
//)
//
//func TestDB(t *testing.T) {
//	// MySQL 连接信息 (请替换成你的 MySQL 配置)
//	dsn := DB.ConnectAddrStr(Parameter.SqlUser, Parameter.SqlPassword, Parameter.SqlHost)
//
//	// 创建 MySQL 管理实例
//	manager, err := DB.NewMySQLManager(dsn)
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
//
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
//	record := &DB.ActionDescribeTriggerLog{
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
//	if err := manager.InsertRecord(Parameter.SqlRecordTableName, record); err != nil {
//		log.Fatalf("插入数据失败: %v", err)
//	}
//
//	// 示例: 查询数据
//	condition := map[string]interface{}{"src_ip_addr": "192.168.1.1"}
//	records, total, err := manager.QueryRecords(Parameter.SqlRecordTableName, condition, 10, 0)
//	if err != nil {
//		log.Fatalf("查询数据失败: %v", err)
//	}
//	log.Printf("查询到 %d 条符合条件的数据", total)
//	for _, r := range records {
//		log.Printf("ID: %d, SrcIP: %s, DstIP: %s, Rule: %s", r.ID, r.SrcIPAddr, r.DstIPAddr, r.RuleName)
//	}
//
//	// 示例: 更新数据
//	updateData := map[string]interface{}{"packet_length": 66}
//	if err := manager.UpdateRecord(Parameter.SqlRecordTableName, condition, updateData); err != nil {
//		log.Fatalf("更新数据失败: %v", err)
//	}
//
//	// 示例: 删除数据
//	if err := manager.DeleteRecord(Parameter.SqlRecordTableName, condition); err != nil {
//		log.Fatalf("删除数据失败: %v", err)
//	}
//}
