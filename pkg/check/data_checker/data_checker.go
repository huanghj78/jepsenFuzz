package datachecker

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ngaut/log"

	_ "gitee.com/opengauss/openGauss-connector-go-pq" // 确保导入OpenGauss的驱动
	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

type Checker struct {
	Addresses []string
	Tables    []string
}

func (c Checker) Check(_ core.Model, _ []core.Operation) (bool, error) {
	var clients []*sql.DB
	for _, address := range c.Addresses {
		dsn := fmt.Sprintf("postgres://testuser:test@123@%s/test?sslmode=disable&target_session_attrs=any", address)
		db, err := sql.Open("opengauss", dsn)
		if err != nil {
			log.Infof("Failed to open database connection: %v", err)
			continue
		}
		clients = append(clients, db)
	}

	for _, table := range c.Tables {
		var firstResultSetStr string
		for i, client := range clients {
			rows, err := client.Query("SELECT  id, balance, balance2 FROM " + table)
			if err != nil {
				log.Infof("Failed to query table %s: %v", table, err)
				return false, err
			}
			defer rows.Close()

			var resultSetStr strings.Builder
			columns, err := rows.Columns()
			if err != nil {
				log.Infof("Failed to get columns from result set: %v", err)
				return false, err
			}

			for rows.Next() {
				values := make([]sql.RawBytes, len(columns))
				scanArgs := make([]interface{}, len(values))
				for i := range values {
					scanArgs[i] = &values[i]
				}

				err = rows.Scan(scanArgs...)
				if err != nil {
					log.Infof("Failed to scan row: %v", err)
					return false, err
				}

				for j, col := range values {
					if col != nil {
						resultSetStr.Write(col)
						// Add a comma after each column except the last one
						if j < len(values)-1 {
							resultSetStr.WriteString(",")
						}
					}
				}
				// Add a newline after each row
				resultSetStr.WriteString("\n")
			}

			if i == 0 {
				// Use the first result set as the baseline
				firstResultSetStr = resultSetStr.String()
				log.Info(firstResultSetStr)
			} else {
				// Compare the current result set with the first one
				log.Info(resultSetStr.String())
				if resultSetStr.String() != firstResultSetStr {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func (c Checker) Name() string {
	return "DataChecker"
}

// compareResultSets 比较两个结果集是否相同
// func compareResultSets(rs1, rs2 *sql.Rows) bool {
// 	columns1, err := rs1.Columns()
// 	if err != nil {
// 		log.Infof("Failed to get columns from result set 1: %v", err)
// 		return false
// 	}
// 	columns2, err := rs2.Columns()
// 	if err != nil {
// 		log.Infof("Failed to get columns from result set 2: %v", err)
// 		return false
// 	}

// 	if len(columns1) != len(columns2) {
// 		return false
// 	}

// 	for i, c1 := range columns1 {
// 		if c1 != columns2[i] {
// 			return false
// 		}
// 	}

// 	// 读取并比较行
// 	for rs1.Next() {
// 		var row1 []interface{}
// 		var row2 []interface{}
// 		dest := make([]interface{}, len(columns1))
// 		for i := range dest {
// 			dest[i] = new(interface{})
// 		}

// 		if !rs2.Next() {
// 			// 结果集2的行数少于结果集1
// 			return false
// 		}

// 		err = rs1.Scan(dest...)
// 		if err != nil {
// 			log.Infof("Failed to scan row from result set 1: %v", err)
// 			return false
// 		}
// 		row1 = make([]interface{}, len(columns1))
// 		for i, v := range dest {
// 			row1[i] = *v.(*interface{})
// 		}

// 		err = rs2.Scan(dest...)
// 		if err != nil {
// 			log.Infof("Failed to scan row from result set 2: %v", err)
// 			return false
// 		}
// 		row2 = make([]interface{}, len(columns1))
// 		for i, v := range dest {
// 			row2[i] = *v.(*interface{})
// 		}

// 		if !compareRows(row1, row2) {
// 			return false
// 		}
// 	}

// 	// 检查是否有额外的行
// 	if rs2.Next() {
// 		// 结果集2的行数多于结果集1
// 		return false
// 	}

// 	return true
// }

// // compareRows 比较两行数据是否相同
// func compareRows(row1, row2 []interface{}) bool {
// 	if len(row1) != len(row2) {
// 		return false
// 	}

// 	for i := range row1 {
// 		if row1[i] != row2[i] {
// 			return false
// 		}
// 	}

// 	return true
// }
