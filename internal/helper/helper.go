package helper

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"os/user"
	"redis-client-go/internal/define"
	"time"
)

func GetConfPath() string {
	current, _ := user.Current()
	return current.HomeDir
}

func GetRedisClient(connectionIdentity string, db int) (*redis.Client, error) {
	conn, err := GetConnection(connectionIdentity)
	if err != nil {
		return nil, err
	}
	redisOption := &redis.Options{
		Addr:         net.JoinHostPort(conn.Addr, conn.Port),
		Username:     conn.Username,
		Password:     conn.Password,
		DB:           db,
		ReadTimeout:  -1,
		WriteTimeout: -1,
	}
	if conn.Type == "ssh" {
		redisOption.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return getRedisConn(conn.SSHUsername, conn.SSHPassword, conn.SSHAddr+":"+conn.SSHPort, addr)
		}
	}
	rdb := redis.NewClient(redisOption)
	return rdb, err
}

func GetConnection(identity string) (*define.Connection, error) {
	conf := new(define.Config)
	nowPath, _ := os.Getwd()
	data, err := os.ReadFile(nowPath + string(os.PathSeparator) + define.ConfigName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	for _, v := range conf.Connections {
		if v.Identity == identity {
			return v, nil
		}
	}
	return nil, errors.New("连接数据不存在")
}

func getRedisConn(username, password, addr, redisAddr string) (net.Conn, error) {
	client, err := getSSHClient(username, password, addr)
	if err != nil {
		return nil, err
	}
	return client.Dial("tcp", redisAddr)
}
func getSSHClient(username, password, addr string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}
	return ssh.Dial("tcp", addr, config)
}
