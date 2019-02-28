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
	example "go-1/PutOrders/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"reflect"
	"strconv"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PutOrders(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("==============api/v1.0/orders  Postorders post succ!!=============")
	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//1通过session得到当前的user_id
	//构建连接缓存的数据
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
	sessioniduserid := req.Sessionid + "user_id"

	value_id := bm.Get(sessioniduserid)
	beego.Info(value_id, reflect.TypeOf(value_id))
	user_id := int(value_id.([]uint8)[0])
	beego.Info(user_id, reflect.TypeOf(user_id))

	//2通过url参数得到当前订单id

	order_id, _ := strconv.Atoi(req.Orderid)

	//3解析客户端请求的json数据得到action参数
	beego.Info(req.Action)
	//得到请求指令
	action := req.Action

	//5查找订单表,找到该订单并确定当前订单状态是wait_accept
	o := orm.NewOrm()

	order := models.OrderHouse{}
	err = o.QueryTable("order_house").Filter("id", order_id).Filter("status", models.ORDER_STATUS_WAIT_ACCEPT).One(&order)
	if err != nil {

		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	if _, err := o.LoadRelated(&order, "House"); err != nil {

		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	house := order.House
	//6校验该订单的user_id是否是当前用户的user_id
	//返回错误json
	if house.User.Id != user_id {

		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = "订单用户不匹配,操作无效"
		return nil
	}
	//7action为accept
	if action == "accept" {
		//如果是接受订单,将订单状态变成待评价状态
		order.Status = models.ORDER_STATUS_WAIT_COMMENT
		beego.Debug("action = accpet!")
	} else if action == "reject" {
		//如果是拒绝接单, 尝试获得拒绝原因,并把拒单原因保存起来
		order.Status = models.ORDER_STATUS_REJECTED
		//8更换订单状态为status为reject
		reason := req.Action
		//添加评论
		order.Comment = reason
		beego.Debug("action = reject!, reason is ", reason)
	}

	//更新该数据到数据库中
	if _, err := o.Update(&order); err != nil {

		rsp.Errno = utils.RECODE_DATAERR
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
