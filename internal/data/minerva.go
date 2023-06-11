package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"minerva/internal/biz"
	"minerva/internal/pkg/sqlparser"
	"strings"
)

type minervaRepo struct {
	data *Data
	log  *log.Helper
}

func (r *minervaRepo) ParseSqlType(ctx context.Context, sql string) (sqlType []biz.SqlTypeItem, err error) {
	sql = strings.TrimSpace(sql)
	sqls := strings.Split(sql, ";")
	sqlTypes := sqlparser.ParseSqlType(sqls)
	return sqlTypes, nil
}

// NewMinervaRepo .
func NewMinervaRepo(data *Data, logger log.Logger) biz.MinervaRepo {
	return &minervaRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
