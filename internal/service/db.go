package service

import (
	"context"
	"errors"
	"fmt"
	"redis-client-go/internal/define"
	"redis-client-go/internal/helper"
	"strconv"
	"strings"
)

func DbList(identity string) ([]*define.DbItem, error) {
	if identity == "" {
		return nil, errors.New("连接唯一标识不能为空")

	}
	rdb, err := helper.GetRedisClient(identity, 0)
	if err != nil {
		return nil, err
	}
	keySpace, err := rdb.Info(context.Background(), "keyspace").Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	v := strings.Split(keySpace, "\n")

	for i := 1; i < len(v)-1; i++ {
		databases := strings.Split(v[i], ":")
		if len(databases) < 2 {
			continue
		}
		vv := strings.Split(databases[1], ",")
		if len(vv) < 1 {
			continue
		}
		keyNumber := strings.Split(vv[0], "=")
		if len(keyNumber) < 2 {
			continue
		}
		num, err := strconv.Atoi(keyNumber[1])
		if err != nil {
			continue
		}
		m[databases[0]] = num
	}

	databasesRes, err := rdb.ConfigGet(context.Background(), "databases").Result()
	if err != nil {
		return nil, err
	}
	if len(databasesRes) < 1 {
		return nil, errors.New("连接数据异常")
	}
	var dbNumStr string
	for _, value := range databasesRes {
		dbNumStr += fmt.Sprintf("%s", value)
	}
	dbNum, err := strconv.ParseInt(dbNumStr, 10, 64)
	if err != nil {
		return nil, err
	}
	data := make([]*define.DbItem, 0)
	for i := 0; i < int(dbNum); i++ {
		item := &define.DbItem{
			Key: "db" + strconv.Itoa(i),
		}
		if n, ok := m["db"+strconv.Itoa(i)]; ok {
			item.Number = n
		}
		data = append(data, item)
	}
	return data, nil
}

func DbInfo(identity string) ([]*define.KeyValue, error) {
	if identity == "" {
		return nil, errors.New("连接唯一标识不能为空")
	}
	rdb, err := helper.GetRedisClient(identity, 0)
	if err != nil {
		return nil, err
	}

	// info 获取数据库键的个数
	keySpace, err := rdb.Info(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	// info 数据格式
	// key:value
	data := make([]*define.KeyValue, 0)
	infos := strings.Split(keySpace, "\n")
	for _, info := range infos {
		v := strings.Split(info, ":")
		if len(v) == 2 {
			data = append(data, &define.KeyValue{
				Key:   v[0],
				Value: v[1],
			})
		}
	}
	return data, nil
}
