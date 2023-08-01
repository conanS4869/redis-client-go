package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"redis-client-go/internal/define"
	"redis-client-go/internal/helper"
)

func ZSetValueDelete(req *define.ZSetValueRequest) error {
	rdb, err := helper.GetRedisClient(req.ConnIdentity, req.Db)
	err = rdb.ZRem(context.Background(), req.Key, req.Member).Err()
	return err
}

func ZSetValueCreate(req *define.ZSetValueRequest) error {
	rdb, err := helper.GetRedisClient(req.ConnIdentity, req.Db)
	err = rdb.ZAdd(context.Background(), req.Key, redis.Z{
		Score:  req.Score,
		Member: req.Member,
	}).Err()
	return err
}
