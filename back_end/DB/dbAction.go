package DB

//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"log"
//)
//
//// ConnectDB 连接到 MySQL 数据库
//func ConnectDB(user, password, host, dbname string) (*sql.DB, error) {
//	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbname)
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		return nil, err
//	}
//	return db, db.Ping()
//}
//
//// EnsureDatabase 检测并创建 IDSdata 数据库
//func EnsureDatabase(db *sql.DB) error {
//	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS IDSdata")
//	return err
//}
//
//// TableExists 检测表是否存在
//func TableExists(db *sql.DB, tableName string) (bool, error) {
//	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'IDSdata' AND table_name = ?"
//	var count int
//	err := db.QueryRow(query, tableName).Scan(&count)
//	if err != nil {
//		return false, err
//	}
//	return count > 0, nil
//}
//
//// InsertData 向表中插入数据
//func InsertData(db *sql.DB, schema string, table string, columns []string, values []interface{}) error {
//	colStr := "(" + joinStrings(columns, ", ") + ")"
//	placeholders := "(" + joinStrings(makePlaceholders(len(values)), ", ") + ")"
//	query := fmt.Sprintf("INSERT INTO %s.%s %s VALUES %s", schema, table, colStr, placeholders)
//
//	_, err := db.Exec(query, values...)
//	return err
//}
//
//// DeleteData 删除表中的数据
//func DeleteData(db *sql.DB, schema, table, condition string, args ...interface{}) error {
//	query := fmt.Sprintf("DELETE FROM %s.%s WHERE %s", schema, table, condition)
//	_, err := db.Exec(query, args...)
//	return err
//}
//
//// UpdateData 更新表中的数据
//func UpdateData(db *sql.DB, schema, table, updates, condition string, args ...interface{}) error {
//	query := fmt.Sprintf("UPDATE %s.%s SET %s WHERE %s", schema, table, updates, condition)
//	_, err := db.Exec(query, args...)
//	return err
//}
//
//// QueryData 查询表中的数据
//func QueryData(db *sql.DB, schema, table, columns, condition string, args ...interface{}) (*sql.Rows, error) {
//	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s", schema, columns, table, condition)
//	return db.Query(query, args...)
//}
//
//// 辅助函数：拼接字符串
//func joinStrings(items []string, sep string) string {
//	result := ""
//	for i, item := range items {
//		if i > 0 {
//			result += sep
//		}
//		result += item
//	}
//	return result
//}
//
//// 辅助函数：生成占位符
//func makePlaceholders(count int) []string {
//	placeholders := make([]string, count)
//	for i := range placeholders {
//		placeholders[i] = "?"
//	}
//	return placeholders
//}
//
//func main() {
//	// 数据库连接参数
//	user := "root"
//	password := "password"
//	host := "127.0.0.1:3306"
//	dbname := ""
//
//	// 连接数据库
//	db, err := ConnectDB(user, password, host, dbname)
//	if err != nil {
//		log.Fatal("连接数据库失败:", err)
//	}
//	defer db.Close()
//
//	// 确保 IDSdata 数据库存在
//	if err := EnsureDatabase(db); err != nil {
//		log.Fatal("创建 IDSdata 失败:", err)
//	}
//	fmt.Println("IDSdata 数据库已准备好")
//
//	// 连接 IDSdata 数据库
//	db, err = ConnectDB(user, password, host, "IDSdata")
//	if err != nil {
//		log.Fatal("连接 IDSdata 数据库失败:", err)
//	}
//
//	// 检测表是否存在
//	table := "example_table"
//	exists, err := TableExists(db, table)
//	if err != nil {
//		log.Fatal("检查表失败:", err)
//	}
//	if exists {
//		fmt.Println("表存在:", table)
//	} else {
//		fmt.Println("表不存在:", table)
//	}
//	//
//	//// 插入数据示例
//	//err = InsertData(db, "example_table", []string{"name", "age"}, []interface{}{"Alice", 25})
//	//if err != nil {
//	//	log.Println("插入数据失败:", err)
//	//} else {
//	//	fmt.Println("数据插入成功")
//	//}
//	//
//	//// 查询数据示例
//	//rows, err := QueryData(db, "example_table", "name, age", "age > ?", 20)
//	//if err != nil {
//	//	log.Fatal("查询数据失败:", err)
//	//}
//	//defer rows.Close()
//	//
//	//for rows.Next() {
//	//	var name string
//	//	var age int
//	//	if err := rows.Scan(&name, &age); err != nil {
//	//		log.Fatal("扫描行失败:", err)
//	//	}
//	//	fmt.Printf("查询到数据: %s, %d\n", name, age)
//	//}
//	//
//	//// 其他操作示例
//	//err = DeleteData(db, "example_table", "age < ?", 18)
//	//err = UpdateData(db, "example_table", "age = age + 1", "name = ?", "Alice")
//}
