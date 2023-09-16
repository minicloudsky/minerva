# Minerva

Minerva is a SQL type detection tool based on TiDB parser, which automatically identifies SQL types such as creating tables, modifying table fields, adding indexes, deleting indexes, and modifying tables. It can be used in scenarios such as SQL risk assessment.
> Minerva, also known as Athena in Roman mythology, in Greek mythology, was the goddess of wisdom, strategy, and war. She was the daughter of Zeus and the Roman version of the goddess of war, Pallas Athena. Minerva is often depicted as a female warrior wearing a golden helmet and carrying a spear and shield. In addition to planning strategies for war, Minerva was considered the guardian of culture, she presided over medicine, poetry, weaving, and music. During the Roman Empire, Minerva was widely worshipped and her temples and statues were built in many cities and regions.


## Usage
### Docker Quickstart
```shell
docker run -p 8088:8088 -p 9000:9000 --name minerva -d littlecloudsky/minerva
```

### Check a sql's type and risk level
```shell
curl --location --request POST 'http://127.0.0.1:8088/minerva/parse-sql-type' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sql": "insert into user(`name`, `city`,sex) values('apple','tokyo',0);update user set sex=0,city='tokyo' where username in ('sudo','ping','telnet');delete from user where name ='sunny' and sex=0;alter table t_user modify username varchar(64) default '' not null comment 'username';"
}'

{
    "code": 0,
    "message": "OK",
    "data": {
        "sqlTypes": [
            {
                "sql": "delete from user where name ='sunny' and sex=0",
                "sqlType": [
                    "Delete"
                ],
                "risk": "Low"
            },
            {
                "sql": "alter table t_user modify username varchar(64) default '' not null comment 'username'",
                "sqlType": [
                    "ModifyColumn"
                ],
                "risk": "Medium"
            },
            {
                "sql": "insert into user(`name`, `city`,sex) values('apple','tokyo',0)",
                "sqlType": [
                    "Insert"
                ],
                "risk": "Low"
            },
            {
                "sql": "update user set sex=0,city='tokyo' where username in ('sudo','ping','telnet')",
                "sqlType": [
                    "Update"
                ],
                "risk": "Low"
            }
        ]
    }
}
```

## License

The Minerva is open-sourced software licensed under the [MIT license](./LICENSE).
