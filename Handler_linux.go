package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"unsafe"
)

//func insert(db *sqlx.DB, v interface{}) error {
//	sql, args := getInsertSql(v)
//	_, err := db.Exec(sql, args...)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func getInsertSql(v interface{}) (string, []interface{}) {
//	elem := reflect.ValueOf(v)
//	if elem.Kind() == reflect.Ptr {
//		elem = elem.Elem()
//	}
//	elemType := elem.Type()
//	tablename := normalize(elemType.Name())
//	numfields := elem.NumField()
//	insertfields := make([]string, 0, numfields)
//	insertvalues := make([]interface{}, 0, numfields)
//	for i := 0; i < numfields; i++ {
//		curfield := elem.Field(i)
//		curstruct := elemType.Field(i)
//		field, val := parseField(curfield, curstruct)
//		insertfields = append(insertfields, "`"+field+"`")
//		insertvalues = append(insertvalues, val)
//	}
//	quotes := strings.Repeat("?,", len(insertfields))
//	quotes = quotes[0 : len(quotes)-1]
//	ret := fmt.Sprintf("insert into %s(%s) values(%s)", tablename, strings.Join(insertfields, ","), quotes)
//	return ret, insertvalues
//}
//
//func parseField(v reflect.Value, s reflect.StructField) (string, interface{}) {
//	fieldname := normalize(s.Name)
//	var fieldvalue interface{}
//	switch v.Kind() {
//	case reflect.Bool:
//		if v.Bool() {
//			fieldvalue = "1"
//		} else {
//			fieldvalue = "0"
//		}
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		fieldvalue = v.Int()
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//		fieldvalue = v.Uint()
//	case reflect.Float32, reflect.Float64:
//		fieldvalue = v.Float()
//	case reflect.String:
//		fieldvalue = v.String()
//	}
//	return fieldname, fieldvalue
//}
//
//func normalize(str string) string {
//	return strings.ToLower(str)
//}

//设置服务端存储的测试用例
func SetRawConfig(c *gin.Context){
	//获取要更新的配置
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println( "new config : "+string(body))
	err := json.Unmarshal(body, &G_TestConfig)

	if err != nil {
		log.Fatal("new case config json.Unmarshal() Fatal,error:%v\n", err)
	}
	//覆盖原来配置文件
	f, err := os.OpenFile("../config/test_case_config.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		log.Println("file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt(body, n)
	}
}

func ReportHandler(c *gin.Context){

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll()fatal : error:%v\n", err)
	}
	var reportData TestResult
	_ = json.Unmarshal(body, &reportData)
	reportData.IP=c.ClientIP()
	//reportData.IP=remoteIp
	log.Printf(" reportData=%+v",reportData)
	//sql插入数据
	res, err := myDB.Exec("insert INTO result(res_ip,res_timestamp,res_lost1,res_lost2,res_lost3,res_lost4,res_lost5,res_lost_upon5,res_timeout,res_decode_timeout,res_encode_timeout,res_usb_test,res_reboot_count) values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		 c.ClientIP(), time.Now().Format("2006/1/2 15:04:05"),reportData.Loss1,
		 reportData.Loss2,reportData.Loss3,reportData.Loss4,reportData.Loss5,
		 reportData.LossAbove5,reportData.NetSlow,reportData.EncodeSlow,reportData.DecodeSlow,
		 reportData.UsbResult,reportData.RebootSuccessRate)
	//保存插入后的id,让后续追加重启测试的结果方便一点
	tmp:=rebootRecord[c.ClientIP()]
	if res!=nil{
		id, _ :=res.LastInsertId()
		log.Println("result table last insert Id : ",id)
		tmp.ID= *(*int)(unsafe.Pointer(&id))
		rebootRecord[c.ClientIP()]=tmp
	}
	if err!=nil{
		log.Println("insert result fail error : "+err.Error())
	}

}
func ErrorHandler(c *gin.Context) {

	//remoteIp:=GetRemoteIP(req.RemoteAddr)
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(c.ClientIP()," erron occurred : "+string(body))

	_,err:=myDB.Exec("insert INTO error(err_timestamp,err_ip,err_msg) values(?,?,?)",
		time.Now().Format("2006/1/2 15:04:05"),c.ClientIP(),string(body))
	if err!=nil{
		log.Println(err.Error())
	}
	delete(rebootRecord,c.ClientIP())
}

//根据重启测试的状态,返回不同的结果
func DispatchTestCase(c *gin.Context) {
	//查看是否在重启map中
	if val, exist := rebootRecord[c.ClientIP()];exist{
			val.AlreadyReboot+=1
			if (val.TotalReboot-val.AlreadyReboot)!=0{
				log.Println(c.ClientIP()+" reboot success ")

				//数据库中更新重启计数
				_,err:=myDB.Exec("update result set res_reboot_count=? where res_id=?",
					strconv.Itoa(val.AlreadyReboot)+"/"+strconv.Itoa(val.TotalReboot), val.ID)
				if err!=nil{
					log.Println("update id : ",val.ID," fail ,error : ",err.Error())
				}
				rebootRecord[c.ClientIP()]=val
				//发送重启回复
				c.String(REBOOT,"reboot")
			}else {
				_,err:=myDB.Exec("update result set res_reboot_count=? where res_id=?",
					strconv.Itoa(val.AlreadyReboot)+"/"+strconv.Itoa(val.TotalReboot), val.ID)
				if err!=nil{
					log.Println("update id : ",val.ID," fail ,error : ",err.Error())
				}
				//完成测试后去除这个ip
				log.Println("test ended,remove client : ",c.ClientIP())
				delete(rebootRecord,c.ClientIP())
				//完成重启测试后,关闭客户机
				c.String(SHUTDOWN,"shutdown")
			}
	}else {
		if G_TestConfig.RebootTest.Active {
		//若不存在就存入,并且发送测试用例
		log.Println(c.ClientIP()+" first connect")
		rebootRecord[c.ClientIP()]=G_TestConfig.RebootTest
		//w.WriteHeader(401)
		//log.Println("testRecorder=%+v", ConfigMap[req.Host])
		responseJson,_:=json.Marshal(G_TestConfig)
		//log.Printf("Load testConfig=%+v\n", G_TestConfig)

		log.Printf("send test config %+v to %+v",G_TestConfig,c.ClientIP())
		c.String(http.StatusOK, string(responseJson))
		}
	}
}