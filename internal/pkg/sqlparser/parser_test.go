package sqlparser

import (
	"fmt"
	"testing"
)

func TestParseSqlType(t *testing.T) {
	sqls := []string{"alter table t_user modify username varchar(64) default '' not null comment 'username';"}
	res := ParseSqlType(sqls)
	fmt.Println(res)
}
