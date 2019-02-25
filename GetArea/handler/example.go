package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"gomicro/IhomeWeb/models"
	"gomicro/IhomeWeb/utils"
	"time"

	"github.com/micro/go-log"

	//	"github.com/astaxie/beego/cache"
	_ "github.com/gomodule/redigo/redis"
	example "gomicro/GetArea/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("请求地区信息 GetArea api/areas")
	//初始化 错误码
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	/**
	1.从缓存中获取数据 如果有发送给前端 如果没有就走2
	2.没有数据就从mysql中查找数据
	3.将查找到的数据存到缓存中
	4.将查找到的数据发送给前端
	*/
	//准备连接redis信息
	redis_conf := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_password,
	}
	beego.Info(redis_conf)
	//将map进行转换成为json
	redis_conf_json, _ := json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("redis连接失败", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//获取数据 定义一个key 就是用来作area查询的 area_info
	area_value := bm.Get("area_info")
	if area_value != nil {
		//如果有数据发送给前端
		beego.Info("获取到地域信息缓存")
		area_map := []map[string]interface{}{}
		//将获取到的数据进行json解码
		json.Unmarshal(area_value.([]byte), &area_map)
		//beego.Info("得到从缓存中提取到的area数据", area_map)
		for _, value := range area_map {
			//beego.Info(key, value)
			tmp := example.Response_Areas{Aid: int32(value["aid"].(float64)), Aname: value["aname"].(string)}
			rsp.Data = append(rsp.Data, &tmp)
		}
		//已经将数据发送给了前端，后面就不需要执行了
		return nil
	}

	//beego操作数据库的orm方法
	//创建orm句柄
	o := orm.NewOrm()
	//查询什么
	qs := o.QueryTable("area")
	//用什么接收
	var area []models.Area
	num, err := qs.All(&area)
	if err != nil {
		beego.Info("数据库查询失败", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	if num == 0 {
		beego.Info("数据库没有数据")
		rsp.Error = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//将查询到的数据写入缓存
	//需要将获取到的数据转换为json
	area_json, _ := json.Marshal(area)
	//操作redis将数据带入
	err = bm.Put("area_info", area_json, time.Second*3600)
	if err != nil {
		beego.Info("数据缓存失败", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//将查询到的数据按照proto的格式发送给web服务
	for _, value := range area {
		//beego.Info(key, value)
		tmp := example.Response_Areas{Aid: int32(value.Id), Aname: value.Name}
		rsp.Data = append(rsp.Data, &tmp)
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
