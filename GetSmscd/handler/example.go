package handler

import (
	"context"

	"github.com/micro/go-log"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/GetSmscd/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"math/rand"
	"reflect"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
//获取短信验证码
func (e *Example) GetSmscd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" GET smscd  api/v1.0/smscode/:id ")
	//初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*验证uuid的缓存*/

	//验证手机号
	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user)
	if err == nil {
		beego.Info("用户已经存在")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info(err)

	//连接redis数据库
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)
	beego.Info(string(redis_config))

	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config))
	//bm,err:=cache.NewCache("redis",`{"key":"ihome","conn":"127.0.0.1:6379","dbNum":"0"} `)//创建1个缓存句柄
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	beego.Info(req.Id, reflect.TypeOf(req.Id))
	//查询相关数据

	value := bm.Get(req.Id)
	if value == nil {
		beego.Info("获取到缓存数据查询失败", value)

		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil
	}
	beego.Info(value, reflect.TypeOf(value))
	value_str, _ := redis.String(value, nil)

	beego.Info(value_str, reflect.TypeOf(value_str))
	//数据对比
	if req.Text != value_str {
		beego.Info("图片验证码 错误 ")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(8999) + 1000 //1000-9999
	beego.Info(size)

	/*	//短信map
		messageconfig := make(map[string]string)
		//id
		messageconfig["appid"] = "29672"
		//key
		messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
		//编码格式
		messageconfig["signtype"] = "md5"


		//短信操作对象
		messagexsend := submail.CreateMessageXSend()
		//短信发送到那个手机号
		submail.MessageXSendAddTo(messagexsend, req.Mobile  )
		//短信发送的模板
		submail.MessageXSendSetProject(messagexsend, "NQ1J94")
		//发送的验证码
		submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(size))
		//发送
		fmt.Println("MessageXSend ", submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig))*/

	/*通过手机号将验证短信进行缓存*/

	err = bm.Put(req.Mobile, size, time.Second*300)
	if err != nil {
		beego.Info("缓存出现问题")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

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
