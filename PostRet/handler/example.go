package handler

import (
	"context"

	"github.com/micro/go-log"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/PostRet/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"reflect"
	"strconv"
	"time"
)

type Example struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostRet(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" POST userreg    /api/v1.0/users !!!")
	//初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	redis_config, _ := json.Marshal(redis_config_map)

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}

	//查询相关数据
	value := bm.Get(req.Mobile)
	if value == nil {
		beego.Info("获取到缓存数据查询失败", value)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}
	beego.Info(value, reflect.TypeOf(value))
	//进行解码
	var info interface{}
	json.Unmarshal(value.([]byte), &info)
	beego.Info(info, reflect.TypeOf(info))

	//类型转换
	s := int(info.(float64))
	beego.Info(s, reflect.TypeOf(s))
	s1, err := strconv.Atoi(req.SmsCode)

	if s1 != s {
		beego.Info("短信验证码错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	user := models.User{}
	user.Name = req.Mobile //就用手机号登陆
	//密码正常情况下 md5 sha256 sm9  存入数据库的是你加密后的编码不是明文存入
	//user.Password_hash = GetMd5String(req.Password)
	user.Password_hash = req.Password
	user.Mobile = req.Mobile
	//创建数据库剧本
	o := orm.NewOrm()
	//插入数据库
	id, err := o.Insert(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info("id", id)

	//生成sessionID 保证唯一性
	h := GetMd5String(req.Mobile + req.Password)
	//返回给客户端session
	rsp.SessionID = h

	//拼接key sessionid + name
	bm.Put(h+"name", string(user.Mobile), time.Second*3600)
	//拼接key sessionid + user_id
	bm.Put(h+"user_id", string(user.Id), time.Second*3600)
	//拼接key sessionid + mobile
	bm.Put(h+"mobile", string(user.Mobile), time.Second*3600)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
