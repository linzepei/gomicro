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
	example "go-1/GetIndex/proto/example"
	"go-1/homeweb/models"
	"go-1/homeweb/utils"
	"time"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetIndex(ctx context.Context, req *example.Request, rsp *example.Response) error {

	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	data := []interface{}{}
	//1 从缓存服务器中请求 "home_page_data" 字段,如果有值就直接返回
	//先从缓存中获取房屋数据,将缓存数据返回前端即可
	redis_config_map := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": "sher",
	}

	redis_config, _ := json.Marshal(redis_config_map)
	beego.Info(string(redis_config))

	cache_conn, err := cache.NewCache("redis", string(redis_config))

	//cache_conn, err := cache.NewCache("redis", `{"key":"ilhome","conn":"127.0.0.1:6379","dbNum":"0"} `)
	if err != nil {
		beego.Debug("connect cache error", err)
	}

	house_page_key := "home_page_data"
	house_page_value := cache_conn.Get(house_page_key)
	if house_page_value != nil {
		beego.Debug("======= get house page info  from CACHE!!! ========")
		//直接将二进制发送给客户端
		rsp.Max = house_page_value.([]byte)
		return nil
	}

	houses := []models.House{}

	//2 如果缓存没有,需要从数据库中查询到房屋列表
	o := orm.NewOrm()

	if _, err := o.QueryTable("house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses); err == nil {
		for _, house := range houses {
			o.LoadRelated(&house, "Area")
			o.LoadRelated(&house, "User")
			o.LoadRelated(&house, "Images")
			o.LoadRelated(&house, "Facilities")
			data = append(data, house.To_house_info())
		}

	}
	beego.Info(data, houses)
	//将data存入缓存数据
	house_page_value, _ = json.Marshal(data)
	cache_conn.Put(house_page_key, house_page_value, 3600*time.Second)

	rsp.Max = house_page_value.([]byte)
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
