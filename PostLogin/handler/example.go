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
	example "go-1/PostLogin/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"time"
)

type Example struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostLogin(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("登陆 api/v1.0/sessions")

	//返回给前端的map结构体
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//查询数据库
	var user models.User
	o := orm.NewOrm()

	//select * from user
	//创建查询句柄
	qs := o.QueryTable("user")
	//qs.Filter("profile__age", 18)
	//查询符合的数据
	err := qs.Filter("mobile", req.Mobile).One(&user)
	if err != nil {

		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}

	beego.Info("密码", req.Password)
	beego.Info("数据库密码", user.Password_hash)

	//判断密码是否正确
	if utils.Md5String(req.Password) != user.Password_hash {
		rsp.Errno = utils.RECODE_PWDERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}

	//编写redis缓存数据库信息
	redis_config_map := map[string]string{
		"key": utils.G_server_name,
		//"conn":"127.0.0.1:6379",
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config))

	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}

	//生成sessionID
	h := GetMd5String(req.Mobile + req.Password)
	rsp.SessionID = h

	beego.Info(h)

	//拼接key sessionid + name
	bm.Put(h+"name", string(user.Name), time.Second*3600)
	//拼接key sessionid + user_id
	bm.Put(h+"user_id", string(user.Id), time.Second*3600)
	//拼接key sessionid + mobile
	bm.Put(h+"mobile", string(user.Mobile), time.Second*3600)

	//成功返回数据
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
