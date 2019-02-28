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
	example "go-1/GetArea/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" GetArea    api/v1.0/areas !!!")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//连接redis创建句柄
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}
	//确定连接信息
	beego.Info(redis_config_map)
	//将map转化为json
	redis_config, _ := json.Marshal(redis_config_map)
	//连接redis
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/*1获取缓存数据*/
	areas_info_value := bm.Get("areas_info")
	//如果不为空则说明成功
	if areas_info_value != nil {
		beego.Info("获取到缓存发送给前端")

		//用来存放解码的json
		ares_info := []map[string]interface{}{}
		//解码
		err = json.Unmarshal(areas_info_value.([]byte), &ares_info)

		//进行循环赋值
		for key, value := range ares_info {

			beego.Info(key, value)
			//创建对于数据类型并进行赋值
			area := example.Response_Address{Aid: int32(value["aid"].(float64)), Aname: value["aname"].(string)}

			//递增到切片
			rsp.Data = append(rsp.Data, &area)
		}

		return nil
	}
	beego.Info("没有拿到缓存")

	/*2如果没有缓存我们就从mysql 里进行查询*/

	//orm的操作创建orm句柄
	o := orm.NewOrm()

	//接受地区信息的切片
	var areas []models.Area
	//创建查询条件
	qs := o.QueryTable("area")
	//查询全部地区
	num, err := qs.All(&areas)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	beego.Info("写入缓存")

	/*3获取数据写入缓存*/

	//将查询到的数据编码成json格式
	ares_info_str, _ := json.Marshal(areas)

	//Put(key string, val interface{}, timeout time.Duration) error
	//存入缓存中
	err = bm.Put("areas_info", ares_info_str, time.Second*3600)
	if err != nil {
		beego.Info("数据库中查出数据信息存入缓存中失误", err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//返回地区信息
	for key, value := range areas {
		beego.Info(key, value)

		area := example.Response_Address{Aid: int32(value.Id), Aname: string(value.Name)}
		//var area  example.ResponseArea_Address
		//area.Aid = value["aid"].(int32)
		//area.Aname = value["aname"].(string)
		rsp.Data = append(rsp.Data, &area)
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
