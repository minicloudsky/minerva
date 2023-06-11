package sqlparser

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"minerva/internal/biz"
	"minerva/internal/pkg"
	"minerva/internal/pkg/zlog"

	"strings"
)

// DdlTypeCreateDatabase databse operation
const DdlTypeCreateDatabase = "create database"
const DdlTypeDropDatabase = "drop database"

// DdlTypeCreateTable table operation
const DdlTypeCreateTable = "create table"
const DdlTypeCreatePartitionTable = "create partition table"
const DdlTypeTruncateTable = "truncate table"
const DdlTypeDropTable = "drop table"
const DdlTypeRenameTable = "rename table"

// DdlTypeAddColumn column operation
const DdlTypeAddColumn = "add column"
const DdlTypeModifyColumn = "modify column"
const DdlTypeChangeColumn = "change column"
const DdlTypeAlterColumn = "alter column"
const DdlTypeDropColumn = "drop column"

// DdlTypeAddIndex indexes
const DdlTypeAddIndex = "add index"
const DdlTypeCreateIndex = "create index"
const DdlTypeDropIndex = "drop index"

// DdlTypePartition DdlTypePartitionTable table partition
const DdlTypePartition = "partition"

const SelectType = "select"
const DMLTypeInsert = "insert"
const DMLTypeUpdate = "update"
const DMLTypeDelete = "delete"

// DdlTypeUnknown unknown ddl
const DdlTypeUnknown = "unknown"

func ParseSqlType(sqls []string) (sqlTypes []biz.SqlTypeItem) {
	sqlTypeMap := make(map[string][]string, 0)
	for _, sql := range sqls {
		sqlTypeMap[sql] = make([]string, 0)
	}
	for _, sql := range sqls {
		p := parser.New()
		strippedSql := strings.ToLower(sql)
		strippedSql = strings.Replace(strippedSql, "\n", " ", -1)
		strippedSql = strings.Replace(strippedSql, "  ", " ", -1)
		strippedSql = strings.Replace(strippedSql, ", ", ",", -1)
		strippedSql = strings.Replace(strippedSql, " ,", ",", -1)
		strippedSql = strings.Replace(strippedSql, " , ", ",", -1)
		strippedSql = strings.Replace(strippedSql, " =", "=", -1)
		strippedSql = strings.Replace(strippedSql, "= ", "=", -1)
		strippedSql = strings.Replace(strippedSql, "algorithm=inplace,lock=none,", "", -1)
		strippedSql = strings.Replace(strippedSql, ",algorithm=inplace,lock=none", "", -1)
		strippedSql = strings.Replace(strippedSql, ",algorithm=instant,lock=default", "", -1)
		strippedSql = strings.Replace(strippedSql, "algorithm=instant,lock=default,", "", -1)
		stmts, _, err := p.Parse(strippedSql, "", "")
		if err != nil {
			zlog.StdLogHelper.Errorf("tidb sql parse err! sql: %s err: %v\n", sql, err)
			if len(sqlTypeMap[sql]) == 0 {
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeUnknown)
			}
			continue
		}
		for _, stmt := range stmts {
			switch stmt.(type) {
			case *ast.DropDatabaseStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeDropDatabase)
			case *ast.DropTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeDropTable)
			case *ast.RenameTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeRenameTable)
			case *ast.TruncateTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeTruncateTable)
			case *ast.DropIndexStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeDropIndex)
			case *ast.AlterTableStmt:
				alterStmt, ok := stmt.(*ast.AlterTableStmt)
				if !ok {
					zlog.StdLogHelper.Errorf("sql: %s not an alter table statement", sql)
				}
				for _, spec := range alterStmt.Specs {
					switch spec.Tp {
					case ast.AlterTableRenameTable:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeRenameTable)
					case ast.AlterTableDropIndex, ast.AlterTableDropForeignKey, ast.AlterTableDropPrimaryKey:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeDropIndex)
					case ast.AlterTableDropColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeDropColumn)
					case ast.AlterTableChangeColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeChangeColumn)
					case ast.AlterTableModifyColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeModifyColumn)
					case ast.AlterTableAlterColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeAlterColumn)
					case ast.AlterTableAddColumns:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeAddColumn)
					case ast.AlterTableAddConstraint:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeAddIndex)
					case ast.AlterTablePartition, ast.AlterTableExchangePartition,
						ast.AlterTableReorganizePartition, ast.AlterTableCheckPartitions,
						ast.AlterTableAddPartitions, ast.AlterTableOptimizePartition,
						ast.AlterTableCoalescePartitions, ast.AlterTableRemovePartitioning,
						ast.AlterTableTruncatePartition, ast.AlterTableRepairPartition,
						ast.AlterTableDropFirstPartition, ast.AlterTableAddLastPartition,
						ast.AlterTableReorganizeFirstPartition, ast.AlterTableReorganizeLastPartition,
						ast.AlterTableDropPartition:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypePartition)
					default:
						if len(sqlTypeMap[sql]) == 0 {
							sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeUnknown)
						}
						zlog.StdLogHelper.Errorf("unknown sql alter type, sql: %s", sql)
					}
				}
			case *ast.CreateIndexStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeCreateIndex)
			case *ast.CreateTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeCreateTable)
				createTableStmt, ok := stmt.(*ast.CreateTableStmt)
				if !ok {
					zlog.StdLogHelper.Errorf("sql: %s not an create table statement", sql)
				}
				if createTableStmt.Partition != nil && createTableStmt.Partition.PartitionMethod.Num > 0 && len(createTableStmt.Partition.Definitions) > 0 {
					sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeCreatePartitionTable)
				}
			case *ast.CreateDatabaseStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeCreateDatabase)
			case *ast.SelectStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], SelectType)
			case *ast.InsertStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DMLTypeInsert)
			case *ast.UpdateStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DMLTypeUpdate)
			case *ast.DeleteStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], DMLTypeDelete)
			default:
				if len(sqlTypeMap[sql]) == 0 {
					sqlTypeMap[sql] = append(sqlTypeMap[sql], DdlTypeUnknown)
				}
				zlog.StdLogHelper.Errorf("unknown sql type! sql: %s", sql)
			}
		}
	}
	for sql, types := range sqlTypeMap {
		uniqTypes := pkg.RemoveDuplicateStr(types)
		sqlTypes = append(sqlTypes, biz.SqlTypeItem{
			Sql:  sql,
			Type: uniqTypes,
		})
	}
	return sqlTypes
}
