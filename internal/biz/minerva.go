package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// SqlType represents DDL operation types
type SqlType string

// SqlType operation constants
const (
	DdlTypeCreateDatabase       SqlType = "CreateDatabase"
	DdlTypeDropDatabase         SqlType = "DropDatabase"
	DdlTypeCreateTable          SqlType = "CreateTable"
	DdlTypeCreatePartitionTable SqlType = "CreatePartitionTable"
	DdlTypeTruncateTable        SqlType = "TruncateTable"
	DdlTypeDropTable            SqlType = "DropTable"
	DdlTypeRenameTable          SqlType = "RenameTable"
	DdlTypeAddColumn            SqlType = "AddColumn"
	DdlTypeModifyColumn         SqlType = "ModifyColumn"
	DdlTypeChangeColumn         SqlType = "ChangeColumn"
	DdlTypeAlterColumn          SqlType = "AlterColumn"
	DdlTypeDropColumn           SqlType = "DropColumn"
	DdlTypeAddIndex             SqlType = "AddIndex"
	DdlTypeCreateIndex          SqlType = "CreateIndex"
	DdlTypeDropIndex            SqlType = "DropIndex"
	DdlTypePartition            SqlType = "Partition"
	SelectType                  SqlType = "Select"
	DMLTypeInsert               SqlType = "Insert"
	DMLTypeUpdate               SqlType = "Update"
	DMLTypeDelete               SqlType = "Delete"
	DdlTypeUnknown              SqlType = "Unknown"
)

type SqlRiskLevel string

const (
	SqlRiskHigh   SqlRiskLevel = "High"
	SqlRiskMedium SqlRiskLevel = "Medium"
	SqlRiskLow    SqlRiskLevel = "Low"
)

func (s SqlType) SqlTypeRisk() SqlRiskLevel {
	switch s {
	case DdlTypeDropDatabase, DdlTypeCreatePartitionTable, DdlTypeTruncateTable, DdlTypeDropTable,
		DdlTypeRenameTable, DdlTypeDropColumn, DdlTypeDropIndex, DdlTypePartition, DdlTypeUnknown:
		return SqlRiskHigh
	case DdlTypeAddColumn, DdlTypeModifyColumn, DdlTypeChangeColumn, DdlTypeAlterColumn,
		DdlTypeAddIndex, DdlTypeCreateIndex:
		return SqlRiskMedium
	case DdlTypeCreateDatabase, DdlTypeCreateTable,
		DMLTypeInsert, DMLTypeUpdate, DMLTypeDelete, SelectType:
		return SqlRiskLow
	default:
		return SqlRiskHigh
	}
}

type SqlTypeCheckResult struct {
	Sql  string
	Type []SqlType
	Risk SqlRiskLevel
}

type MinervaRepo interface {
	ParseSqlType(ctx context.Context, sql string) (sqlType []SqlTypeCheckResult, err error)
}

// MinervaUsecase is a Minerva usecase.
type MinervaUsecase struct {
	repo MinervaRepo
	log  *log.Helper
}

// NewMinervaUsecase new a Minerva usecase.
func NewMinervaUsecase(repo MinervaRepo, logger log.Logger) *MinervaUsecase {
	return &MinervaUsecase{repo: repo, log: log.NewHelper(logger)}
}
