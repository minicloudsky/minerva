package sqlparser

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/samber/lo"
	"minerva/internal/biz"
	"minerva/internal/pkg/mlog"

	"strings"
)

func ParseSqlType(sqls []string) (sqlTypes []biz.SqlTypeCheckResult) {
	sqlTypeMap := make(map[string][]biz.SqlType, 0)
	for _, sql := range sqls {
		sqlTypeMap[sql] = make([]biz.SqlType, 0)
	}
	for _, sql := range sqls {
		p := parser.New()
		cleanSql := strings.ToLower(sql)
		cleanSql = strings.Replace(cleanSql, "\n", " ", -1)
		cleanSql = strings.Replace(cleanSql, "  ", " ", -1)
		cleanSql = strings.Replace(cleanSql, ", ", ",", -1)
		cleanSql = strings.Replace(cleanSql, " ,", ",", -1)
		cleanSql = strings.Replace(cleanSql, " , ", ",", -1)
		cleanSql = strings.Replace(cleanSql, " =", "=", -1)
		cleanSql = strings.Replace(cleanSql, "= ", "=", -1)
		stmts, _, err := p.Parse(cleanSql, "", "")
		if err != nil {
			mlog.StdLogHelper.Errorf("tidb sql parse err! sql: %s err: %v\n", sql, err)
			if len(sqlTypeMap[sql]) == 0 {
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeUnknown)
			}
			continue
		}
		for _, stmt := range stmts {
			switch stmt.(type) {
			case *ast.DropDatabaseStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeDropDatabase)
			case *ast.DropTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeDropTable)
			case *ast.RenameTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeRenameTable)
			case *ast.TruncateTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeTruncateTable)
			case *ast.DropIndexStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeDropIndex)
			case *ast.AlterTableStmt:
				alterStmt, ok := stmt.(*ast.AlterTableStmt)
				if !ok {
					mlog.StdLogHelper.Errorf("sql: %s not an alter table statement", sql)
				}
				for _, spec := range alterStmt.Specs {
					switch spec.Tp {
					case ast.AlterTableRenameTable:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeRenameTable)
					case ast.AlterTableDropIndex, ast.AlterTableDropForeignKey, ast.AlterTableDropPrimaryKey:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeDropIndex)
					case ast.AlterTableDropColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeDropColumn)
					case ast.AlterTableChangeColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeChangeColumn)
					case ast.AlterTableModifyColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeModifyColumn)
					case ast.AlterTableAlterColumn:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeAlterColumn)
					case ast.AlterTableAddColumns:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeAddColumn)
					case ast.AlterTableAddConstraint:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeAddIndex)
					case ast.AlterTablePartition, ast.AlterTableExchangePartition,
						ast.AlterTableReorganizePartition, ast.AlterTableCheckPartitions,
						ast.AlterTableAddPartitions, ast.AlterTableOptimizePartition,
						ast.AlterTableCoalescePartitions, ast.AlterTableRemovePartitioning,
						ast.AlterTableTruncatePartition, ast.AlterTableRepairPartition,
						ast.AlterTableDropFirstPartition, ast.AlterTableAddLastPartition,
						ast.AlterTableReorganizeFirstPartition, ast.AlterTableReorganizeLastPartition,
						ast.AlterTableDropPartition:
						sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypePartition)
					default:
						if len(sqlTypeMap[sql]) == 0 {
							sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeUnknown)
						}
						mlog.StdLogHelper.Errorf("unknown sql alter type, sql: %s", sql)
					}
				}
			case *ast.CreateIndexStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeCreateIndex)
			case *ast.CreateTableStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeCreateTable)
				createTableStmt, ok := stmt.(*ast.CreateTableStmt)
				if !ok {
					mlog.StdLogHelper.Errorf("sql: %s not an create table statement", sql)
				}
				if createTableStmt.Partition != nil &&
					createTableStmt.Partition.PartitionMethod.Num > 0 &&
					len(createTableStmt.Partition.Definitions) > 0 {
					sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeCreatePartitionTable)
				}
			case *ast.CreateDatabaseStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeCreateDatabase)
			case *ast.SelectStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.SelectType)
			case *ast.InsertStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DMLTypeInsert)
			case *ast.UpdateStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DMLTypeUpdate)
			case *ast.DeleteStmt:
				sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DMLTypeDelete)
			default:
				if len(sqlTypeMap[sql]) == 0 {
					sqlTypeMap[sql] = append(sqlTypeMap[sql], biz.DdlTypeUnknown)
				}
				mlog.StdLogHelper.Errorf("unknown sql type! sql: %s", sql)
			}
		}
	}
	for sql, types := range sqlTypeMap {
		uniqTypes := lo.Uniq[biz.SqlType](types)
		var risks []biz.SqlRiskLevel
		var higestRisk biz.SqlRiskLevel
		for _, t := range uniqTypes {
			risks = append(risks, t.SqlTypeRisk())
		}
		if lo.Contains[biz.SqlRiskLevel](risks, biz.SqlRiskHigh) {
			higestRisk = biz.SqlRiskHigh
		} else if lo.Contains[biz.SqlRiskLevel](risks, biz.SqlRiskMedium) {
			higestRisk = biz.SqlRiskMedium
		} else {
			higestRisk = biz.SqlRiskLow
		}
		sqlTypes = append(sqlTypes, biz.SqlTypeCheckResult{
			Sql:  sql,
			Type: uniqTypes,
			Risk: higestRisk,
		})
	}
	return sqlTypes
}
