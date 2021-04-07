package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logFile *os.File
	myDB *sqlx.DB
	rebootRecord map[string]RebootTest
)
//初始化
func init(){
	//输出日志到控制台和文件
	gin.DisableConsoleColor()

	rebootRecord=make(map[string]RebootTest)
	logFile, _= os.Create("./log/"+time.Now().String()+".log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.Println("init")
	//初始化数据库
	db, err:= sqlx.Open("mysql", "root:123456@tcp(10.1.108.250:3300)/PreTest?charset=utf8")
	myDB=db
	if err!=nil{
		log.Fatal("open mysql fail , error = ",err)
	}

	//读取配置文件
	jsonfile,err:=ioutil.ReadFile("../config/test_case_config.json")
	if err!=nil{
		log.Fatal("read case config fail, error = ",err)
	}
	err = json.Unmarshal(jsonfile, &G_TestConfig)
	if err != nil {
		log.Fatal("json.Unmarshal() Fatal,error:%v\n", err)
	}
	log.Printf("Load testConfig=%+v\n", G_TestConfig)

}

func main() {

	router := gin.Default()
	defer myDB.Close()
	defer logFile.Close()
	//读取网页
	router.LoadHTMLGlob("../web/*")

	// API 定义开始
	//分发测试用例
	router.GET("/case",DispatchTestCase)

	//网页
	web:=router.Group("/web")
	{
		web.GET("/config", func(context *gin.Context) {
			// 返回HTML文件
			context.HTML(http.StatusOK, "SetTestConfig.html", nil)
		})

		web.GET("/result", func(context *gin.Context) {
			// 返回HTML文件
			context.HTML(http.StatusOK, "ShowTestResult.html", nil)
		})

		web.GET("/error", func(context *gin.Context) {
			// 返回HTML文件
			context.HTML(http.StatusOK, "ShowError.html", nil)
		})
	}

	//处理上报的数据
	reportHandler := router.Group("/report")
	{
		reportHandler.POST("/streamResult", ReportHandler)
		reportHandler.POST("/error", ErrorHandler)
	}

	//网页内嵌js查询接口
	query := router.Group("/query")
	{
		query.GET("/rawConfig", func(context *gin.Context) {
			context.JSON(http.StatusOK,G_TestConfig)
		})

		query.GET("/result", func(context *gin.Context) {
			//数据库查询
			var res []TestResult
			//todo 需要优化,这里一次性读取数据库会有内存问题
			err := myDB.Select(&res, "select res_ip,res_timestamp,res_lost1,res_lost2,res_lost3,res_lost4,res_lost5,res_lost_upon5,res_timeout,res_decode_timeout,res_encode_timeout,res_usb_test,res_reboot_count from result")
			if err != nil {
				log.Fatal("exec failed, ", err)
				return
			}
			var jsons string
			for _,item := range res{
				temp,_:=json.Marshal(item)
				jsons+=string(temp)+"\n"
			}
			//发送结果
			context.String(http.StatusOK,jsons)

		})

		query.GET("/error", func(context *gin.Context) {
			var e []Error
			//todo 需要优化,这里一次性读取数据库会有内存问题
			err := myDB.Select(&e, "select * from error")
			if err != nil {
				log.Fatal("exec failed, ", err)
				return
			}
			var jsons string
			for _,item := range e{
				temp,_:=json.Marshal(item)
				jsons+=string(temp)+"\n"
			}
			//发送结果
			context.String(http.StatusOK,jsons)
		})
	}
	//设置测试用例
	set := router.Group("/set")
	{
		set.POST("/rawConfig", SetRawConfig)
	}

	err := router.Run(":80")
	if err!=nil{
		log.Fatal("gin run fail")
	}

}


//func main() {
//	CreateDir("./log")
//	logfile:=setLogOutput("./log/serverlog.txt")
//	defer logfile.Close()
//
//	//保存测试进度
//	ConfigMap=make(map[string]TestConfig)
//	ResultMap=make(map[string]TestResult)
//	//读取配置文件
//	jsonfile,_:=ioutil.ReadFile("../config/test_case_config.json")
//
//	err := json.Unmarshal(jsonfile, &G_TestConfig)
//	if err != nil {
//		log.Fatal("json.Unmarshal() Fatal,error:%v\n", err)
//	}
//	log.Printf("Load testConfig=%+v\n", G_TestConfig)
//	//运行http服务
//
//	http.HandleFunc(dispatchTestCase,DispatchTestCase)
//
//	http.HandleFunc(reportHandler,ReportHandler)
//
//	http.HandleFunc(setRawConfig,SetRawConfig)
//
//	http.HandleFunc(errorHandler,ErrorHandler)
//
//	http.HandleFunc(queryResult,func(w http.ResponseWriter, req *http.Request){
//		log.Println("html query test result")
//
//		for _, v := range ResultMap {
//			js,err :=json.Marshal(v)
//			if err!=nil{
//				log.Fatal(err.Error())
//			}
//			_, _ = fmt.Fprintf(w, string(js)+"\n")
//			//log.Println("write json ",string(js))
//		}
//	})
//
//	http.HandleFunc(queryRawConfig,func(w http.ResponseWriter, req *http.Request){
//		log.Println("html query raw config")
//		_, _ = fmt.Fprintf(w, string(jsonfile))
//		//log.Println("write json ",string(js))
//	})
//
//	http.HandleFunc(showTestResult,func(w http.ResponseWriter, req *http.Request){
//
//		log.Println("send ShowTestResult.html")
//		s, e := ioutil.ReadFile("../web/ShowTestResult.html") //读取html
//		if e != nil {
//			panic(e)
//		}
//		_, _ = fmt.Fprintf(w, string(s)) //这个就是response
//	})
//
//	http.HandleFunc(setTestConfig,func(w http.ResponseWriter, req *http.Request){
//
//		log.Println("send SetTestConfig.html")
//		s, e := ioutil.ReadFile("../web/SetTestConfig.html") //读取html
//		if e != nil {
//			panic(e)
//		}
//		_, _ = fmt.Fprintf(w, string(s)) //这个就是response
//	})
//
//
//
//	//http.HandleFunc("/activeReboot",activeReboot)
//	err = http.ListenAndServe(":80", nil)
//	if err != nil {
//		panic(err)
//	}
//
//}

//package main
//
//import (
//"encoding/json"
//"github.com/gin-gonic/gin"
//"io"
//"io/ioutil"
//"log"d
//"net/http"
//
//)
////保存每个盒子的测试进度
//var ClientMap map[string]TestConfig
//

//
//}
//
////根据服务端存储的测试进度进行测试用例分发
//func dispatchTestCase(w http.ResponseWriter, req *http.Request) {
//
//	//获取这台主机的测试进度
//	remoteIP:=getRemoteIP(req.RemoteAddr)
//
//	//若存在则检查重启标志位
//	if val, exist := ClientMap[remoteIP];exist{
//		if val.ServerTest.RebootCase.Active {
//			if val.ServerTest.RebootCase.TestCount!=0{
//				log.Println(remoteIP+" reboot success ")
//				val.ServerTest.RebootCase.FirstConnect=false
//				val.ServerTest.RebootCase.TestCount-=1
//				val.ServerTest.RebootCase.SuccessCount+=1
//				ClientMap[remoteIP]=val
//				//发送更新后的测试用例
//				log.Printf("updated testRecorder=%+v\n", ClientMap[remoteIP])
//				responseJson,_:=json.Marshal(val)
//				_, _ = io.WriteString(w, string(responseJson))
//			}else {
//				//已经完成重启测试就计算成功率
//				temp:=ClientMap[remoteIP]
//				temp.TestResult.RebootSuccessRate=float64( ClientMap[remoteIP].ServerTest.RebootCase.SuccessCount/
//					testConfig.ServerTest.RebootCase.TestCount)
//				ClientMap[remoteIP]=temp
//				log.Printf(remoteIP+" TestResult=%+v\n", ClientMap[remoteIP].TestResult)
//				//完成测试后去除这个ip
//				delete(ClientMap,remoteIP)
//			}
//		}
//
//
//	}else {
//		//若不存在就存入map,并且发送初始测试用例
//		log.Println(remoteIP+" first connect")
//		ClientMap[remoteIP]=testConfig
//		w.Header().Set("Content-Type","application/json")
//		//w.WriteHeader(401)
//		//log.Println("testRecorder=%+v", ClientMap[req.Host])
//		responseJson,_:=json.Marshal(testConfig)
//		_, _ = io.WriteString(w, string(responseJson))
//	}
//}
//
//func reportHandler(w http.ResponseWriter, req *http.Request){
//
//	remoteIp:=getRemoteIP(req.RemoteAddr)
//
//	body, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		log.Fatal("ioutil.ReadAll()fatal : error:%v\n", err)
//	}
//	var reportData TestResult
//	_ = json.Unmarshal(body, &reportData)
//
//	log.Printf(remoteIp+ " reportData=%+v",reportData)
//
//	temp:=ClientMap[remoteIp]
//	temp.TestResult=reportData
//	ClientMap[remoteIp]=temp
//
//}
