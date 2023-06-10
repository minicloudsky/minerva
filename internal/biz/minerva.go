package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Minerva is a Minerva model.
type Minerva struct {
	Hello string
}

type SqlTypeItem struct {
	Sql  string
	Type []string
}

type MinervaRepo interface {
	ParseSqlType(ctx context.Context, sql string) (sqlType []SqlTypeItem, err error)
}

// MinervaUsecase is a Minerva usecase.
type MinervaUsecase struct {
	repo MinervaRepo
	log  *log.Helper
}

// NewMinervausecase new a Minerva usecase.
func NewMinervausecase(repo MinervaRepo, logger log.Logger) *MinervaUsecase {
	return &MinervaUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *MinervaUsecase) ParseSqlType(ctx context.Context, sql string) (sqlType []SqlTypeItem, err error) {
	sqlTypeItems, err := uc.repo.ParseSqlType(ctx, sql)
	if err != nil {
		return nil, err
	}

	return sqlTypeItems, err
}
