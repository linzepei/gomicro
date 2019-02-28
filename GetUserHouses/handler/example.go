package handler

import (
	"context"

	"github.com/micro/go-log"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	example "go-1/GetUserHouses/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"reflect"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//打印被调用的函数
	beego.Info("获取当前用户所发布的房源 GetUserHouses /api/v1.0/user/houses")
	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*通过session 获取我们当前登陆用户的user_id*/
	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
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
	//拼接key
	sessioniduserid := req.Sessionid + "user_id"

	value_id := bm.Get(sessioniduserid)
	beego.Info(value_id, reflect.TypeOf(value_id))
	id := int(value_id.([]uint8)[0])
	beego.Info(id, reflect.TypeOf(id))

	/*通过user_id 获取到当前的用户所发布的房源信息*/
	house_list := []models.House{}

	//创建数据库句柄
	o := orm.NewOrm()
	qs := o.QueryTable("house")

	num, err := qs.Filter("user__id", id).All(&house_list)
	if err != nil {

		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	if num == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	/*成功返回数据给前端*/

	house, err := json.Marshal(house_list)

	rsp.Mix = house

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
