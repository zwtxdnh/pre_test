package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
)

func testcode() {

	//s := hook.Start()
	////定时
	//go func() {
	//	time.Sleep(5*time.Second)
	//	hook.End()
	//	//isTriggered<-false
	//}()
	//
	//hook.Register(hook.KeyDown, []string{"k"}, func(e hook.Event) {
	//	//isTriggered<-true
	//	log.Println("k")
	//	hook.End()
	//})
	//hook.UnRegister()
	//<-hook.Process(s)
	//
	//time.Sleep(10*time.Second)
	//ch:=make(chan bool)
	//robotgo.AddMouse()
	//hookEvent:=hook.Start()
	//
	////初始化一个usb事件通知通道
	//var usbCh = make(chan bool)
	//usbEvent:=getRandomUsbEvent()
	//usbEvent.deviceFlag=keybord
	//usbEvent.keyEvent=KEY_DOWN
	////等待usb事件触发
	//go registerEvent(,usbEvent,usbCh)
	////time.Sleep(time.Millisecond*500)
	////通知盒子发送usb事件
	////sendUsbEvent2Box(usbEvent)
	//go func() {
	//	<-hook.Process(hookEvent)
	//}()
	//if <-usbCh {
	//	log.Printf("usbEvent=%+v is triggered\n\n", usbEvent)
	//} else {
	//	log.Printf("usbEvent=%+v is no triggered\n\n", usbEvent)
	//}
	//
	//hook.UnRegister()
	//hook.End()

	//fi, err := os.Open("log/test1.log")
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err)
	//	return
	//}
	//defer fi.Close()
	//buf := bufio.NewScanner(fi)
	//for {
	//	if !buf.Scan() {
	//		break
	//	}
	//	line := buf.Text()
	//
	//	encodeReg:=regexp.MustCompile(`encode\[avg:(\d+)us, max:(\d+)us, slow:\d+, total_slow:\d+\], decode`)
	//	encodeData:=encodeReg.FindAllStringSubmatch(line,-1)
	//	if encodeData != nil{
	//		fmt.Printf("%+v\n",encodeData)
	//		//continue
	//	}
	//	decodeReg:=regexp.MustCompile(`decode\[avg:(\d+)us, max:(\d+)us, slow:\d+, total_slow:\d+\], net`)
	//	decodeData:=decodeReg.FindAllStringSubmatch(line,-1)
	//	if decodeData != nil{
	//		fmt.Printf("%+v\n",decodeData)
	//		//continue
	//	}
	//	delayReg:=regexp.MustCompile(`delay\[avg:(\d+)us, max:(\d+)us, slow:\d+, total_slow:\d+\]`)
	//	delayData:=delayReg.FindAllStringSubmatch(line,-1)
	//	if delayData != nil{
	//		fmt.Printf("%+v\n",delayData)
	//		//continue
	//	}
	//	fpsReg:=regexp.MustCompile(`\[capture_frame:(\d+), play_stream:(\d+), render:(\d+)`)
	//	fpsData:=fpsReg.FindAllStringSubmatch(line,-1)
	//	if fpsData != nil{
	//		fmt.Printf("%+v\n",fpsData)
	//		//continue
	//	}
	//
	//
	//}
//loadConfig()
log.Println(check_video())
log.Println(check_audio())

os.Exit(0)
}


var hostConfig HostConfig

func checkError(err error) {
	if err != nil {
		log.Println(err.Error())
		e := fmt.Sprintf("%s", err.Error())
		_,_ = http.Post(hostConfig.URL+"error",
			"application/x-www-form-urlencoded",
			bytes.NewBuffer([]byte(e)))
		os.Exit(1)
	}
}

func main() {
	//testcode()
	CreateDir("./log")
	logfile:=setLogOutput("./log/clientlog.txt")
	defer logfile.Close()

	content, err := ioutil.ReadFile("../config/host_config.json")
	checkError(err)

	// 获取配置
	err = json.Unmarshal(content, &hostConfig)
	checkError(err)
	log.Printf("load host config %+v",hostConfig)

	//这里主要根据串流程序的tcp链接去找到盒子端的ip,
	//然后根据盒子端的ip连接到盒子的键鼠模拟器程序
	Connect2Box()

	//向服务器请求测试用例
	GetTestCase(hostConfig.URL+"/case")

	//运行所有测试用例,usb,视频,游戏片段等...
	RunAllCase()
	//time.Sleep(time.Minute)

}



