package user

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type Info struct {
	UID        string `json:"uid"`
	LoginTime  int64  `json:"loginTime"`
	ExpireTime int64  `json:"expireTime"`
}

const (
	TokenPrefix = "Bearer "
	TokenHeader = "Authorization"
)

var authKey = []byte("seckill2021")

func padding(src []byte, blkSize int) []byte {
	l := len(src)
	for i := 0; i < blkSize-l%blkSize; i++ {
		src = append(src, byte(0))
	}
	return src
}

func Auth(token string) *Info {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	cipher, err := aes.NewCipher(padding(authKey, 16))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	src, err1 := base64.StdEncoding.DecodeString(token)
	if err1 != nil || len(src) == 0 {
		return nil
	}
	src = padding(src, cipher.BlockSize())
	output := make([]byte, len(src))
	cipher.Decrypt(output, src)
	var info *Info
	err = json.Unmarshal(output, &info)
	if err != nil || info.ExpireTime < time.Now().Unix() {
		return nil
	}
	return info
}

func Login(uid string, passwd string) (*Info, string) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	info := &Info{
		UID:        uid,
		LoginTime:  time.Now().Unix(),
		ExpireTime: time.Now().Unix() + 24*3600,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return nil, ""
	}
	cipher, err1 := aes.NewCipher(padding(authKey, 16))
	if err1 != nil {
		logrus.Error(err1)
		return nil, ""
	}
	data1 := padding(data, cipher.BlockSize())
	dst := make([]byte, len(data1))
	cipher.Encrypt(dst, data1)
	logrus.Info(len(data), len(data1))
	return info, base64.StdEncoding.EncodeToString(dst[:len(data)])
}
