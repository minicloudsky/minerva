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

func TestParseSelect(t *testing.T) {
	sqls := []string{"select id,name from user where id=1 and name='tony'"}
	res := ParseSqlType(sqls)
	fmt.Println(res)
}

func TestParseDML(t *testing.T) {
	sqls := []string{
		"insert into user(`name`, `city`,sex) values('apple','tokyo',0)",
		"update user set sex=0,city='tokyo' where username in ('sudo','ping','telnet')",
		"delete from user where name ='sunny' and sex=0",
	}
	res := ParseSqlType(sqls)
	fmt.Println(res)
}
