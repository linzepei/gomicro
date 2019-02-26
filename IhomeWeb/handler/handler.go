package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	example "github.com/micro/examples/template/srv/proto/example"
	"github.com/micro/go-micro/client"
	//调用area的proto
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	DELETESESSION "gomicro/DeleteSession/proto/example"
	GETAREA "gomicro/GetArea/proto/example"
	GETIMAGECD "gomicro/GetImageCd/proto/example"
	GETSESSION "gomicro/GetSession/proto/example"
	GETSMSCD "gomicro/GetSmscd/proto/example"
	"gomicro/IhomeWeb/models"
	"gomicro/IhomeWeb/utils"
	POSTLOGIN "gomicro/PostLogin/proto/example"
	POSTRET "gomicro/PostRet/proto/example"
	"image"
	"image/png"
)

func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取地区信息
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求地区信息 GetArea api/v1.0/areas")
	//创建服务获取句柄
	server := grpc.NewService()
	//服务初始化
	server.Init()

	//调用服务返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())

	//调用服务返回数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//接收数据
	//准备接收切片
	area_list := []models.Area{}
	//循环接收数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname}
		area_list = append(area_list, tmp)
	}

	// 返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
		"data":   area_list,
	}

	//会传数据的时候三直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取验证码图片
func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	beego.Info("获取验证码图片 GetImageCd /api/v1.0/imagecode/:uuid")

	//创建服务
	server := grpc.NewService()
	server.Init()

	// 调用服务
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", server.Client())

	//获取uuid
	uuid := ps.ByName("uuid")

	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//接收图片信息的 图片格式
	var img image.RGBA

	img.Stride = int(rsp.Stride)
	img.Pix = []uint8(rsp.Pix)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)

	var image captcha.Image

	image.RGBA = &img

	//将图片发送给浏览器
	png.Encode(w, image)

}

//获取短信验证码GetSmscd
func GetSmscd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取短信验证码 GetSmscd api/v1.0/smscode/:mobile ")
	//通过传入参数URL下Query获取前端的在url里的带参
	//beego.Info(r.URL.Query())
	//map[text:[9346] id:[474494b0-18eb-4eb7-9e68-a5ecf3c8317b]]
	//获取参数
	test := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]
	mobile := ps.ByName("mobile")

	//通过正则进行手机号的判断
	//创建正则条件
	/*	mobile_reg:=regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
		//通过条件判断字符串是否匹配规则 返回正确或失败
		bl :=mobile_reg.MatchString(mobile)
		//如果手机号不匹配那就直接返回错误不调用服务
		if bl == false{
			// 创建返回数据的map
			response := map[string]interface{}{
				"error": utils.RECODE_MOBILEERR ,
				"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
			}

			//设置返回数据的格式
			w.Header().Set("Content-Type","application/json")

			// 发送数据
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	*/

	//创建并初始化服务
	server := grpc.NewService()
	server.Init()

	// 调用服务
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmscd", server.Client())
	rsp, err := exampleClient.GetSmscd(context.TODO(), &GETSMSCD.Request{
		Mobile:   mobile,
		Imagestr: test,
		Uuid:     id,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 创建返回数据的map
	response := map[string]interface{}{

		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
	}

	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")

	// 发送数据
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//PostRet
func PostRet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("PostRet  注册 /api/v1.0/users")

	//服务创建
	server := grpc.NewService()
	server.Init()

	//接收post发送过来的数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["mobile"].(string) == "" || request["password"].(string) == "" || request["sms_code"].(string) == "" {
		//准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//发送给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return

	}

	// 调用请求
	exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", server.Client())
	rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//读取cookie   统一cookie   userlogin
	//func (r *Request) Cookie(name string) (*Cookie, error)

	cookie, err := r.Cookie("userlogin")
	if err != nil || "" == cookie.Value {
		//创建1个cookie对象
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionId, Path: "/", MaxAge: 3600}
		//对浏览器的cookie进行设置
		http.SetCookie(w, &cookie)
	}

	//准备回传数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取session信息
func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取session信息 GetSession /api/v1.0/session")

	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		// 直接返回说名用户未登陆
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		// 将数据回发给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", server.Client())
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]string)
	data["name"] = rsp.UserName

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//登陆  /api/v1.0/sessions
func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("登陆  PostLogin /api/v1.0/sessions")

	// 接收前端发送过来的json数据进行解码
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["mobile"].(string) == "" || request["password"].(string) == "" {
		//准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//发送给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		return
	}

	//创建服务
	server := grpc.NewService()
	server.Init()

	// 调用服务
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", server.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//设置cookie
	//Cookie读取
	cookie, err := r.Cookie("userlogin")

	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.Sessionid, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//退出登陆
func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	beego.Info("DeleteSession  退出登陆 /api/v1.0/session")

	//创建服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", server.Client())

	//获取cookie
	cookie, err := r.Cookie("userlogin")

	if err != nil || cookie.Value == "" {
		//准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//发送给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//删除sessionid
	cookie, err = r.Cookie("userlogin")
	if cookie.Value != "" || err == nil {
		cookie := http.Cookie{Name: "userlogin", Path: "/", MaxAge: -1, Value: ""}
		http.SetCookie(w, &cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取首页轮播图
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取首页轮播图 GetIndex api/v1.0/houses/index")

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	//会传数据的时候三直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
