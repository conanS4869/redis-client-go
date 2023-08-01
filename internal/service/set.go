package service

import (
	"context"
	"redis-client-go/internal/define"
	"redis-client-go/internal/helper"
)

func SetValueDelete(req *define.SetValueRequest) error {
	rdb, err := helper.GetRedisClient(req.ConnIdentity, req.Db)
	err = rdb.SRem(context.Background(), req.Key, req.Value).Err()
	return err
}

func SetValueCreate(req *define.SetValueRequest) error {
	rdb, err := helper.GetRedisClient(req.ConnIdentity, req.Db)
	err = rdb.SAdd(context.Background(), req.Key, req.Value).Err()
	return err
}
