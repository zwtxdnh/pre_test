package main

import (
	hook "github.com/robotn/gohook"
	"log"
	"math/rand"
	"time"
)
var (
	usbAccuracyRate  float64
	screenX int
	screenY int

)

func sendUsbEvent2Box(usbEvent USBProtocol){
	data:=USBStructToBytes(&usbEvent)
	log.Printf("send usbEvent=%+v\n to box", usbEvent)
	log.Println("len = ",len(data),"data = ",data)
	_, err := box.Write(data)
	checkError(err)
}

func GenerateRangeNum(min, max int) int {
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	//fmt.Printf("rand is %v\n", randNum)
	return randNum
}

//func setUsbEvent(usbEvent USBProtocol,triggerCh chan bool){
//	//若要发送的是键盘事件
//	if usbEvent.deviceFlag==keybord{
//		log.Printf("set keyevent = "+key_map[usbEvent.keyEvent])
//
//		registerEvent(usbEvent,triggerCh)
//
//	}else {//若要发送的是鼠标事件
//		switch usbEvent.mouseEvent.mouseCode {
//		case MBTN:{
//			log.Println("pass MBTN event")
//
//			break
//		}
//		case LBTN:{
//			log.Printf("set mouse_event LBTN X=%+v Y=%+v",usbEvent.mouseEvent.X,usbEvent.mouseEvent.Y)
//			robotgo.AddMouse("left")
//			x,y:=robotgo.GetMousePos()
//			log.Printf("GetMousePos x = %+v y = %+v",x,y)
//			if x == int(usbEvent.mouseEvent.X) && y == int(usbEvent.mouseEvent.Y) {
//				triggerCh<-true
//			}
//			break
//
//		}
//		case RBTN:{
//			log.Printf("set mouse_event RBTN X=%+v Y=%+v", usbEvent.mouseEvent.X,usbEvent.mouseEvent.Y)
//			robotgo.AddMouse("right")
//			x,y:=robotgo.GetMousePos()
//			log.Printf("GetMousePos x = %+v y = %+v",x,y)
//			if x == int(usbEvent.mouseEvent.X) && y == int(usbEvent.mouseEvent.Y) {
//				triggerCh<-true
//			}
//			break
//		}
//		case 0:{
//			log.Printf("set mouse_event NOBTN X=%+v Y=%+v", usbEvent.mouseEvent.X,usbEvent.mouseEvent.Y)
//			time.Sleep(time.Second)
//			x,y:=robotgo.GetMousePos()
//			log.Printf("GetMousePos x = %+v y = %+v",x,y)
//			if x == int(usbEvent.mouseEvent.X) && y == int(usbEvent.mouseEvent.Y) {
//				triggerCh<-true
//			}
//			break
//		}
//		}
//	}
//}


func getRandomUsbEvent()USBProtocol {

	//return  USBProtocol{deviceFlag:keybord,keyEvent:key_list[rand.Intn(len(key_list))]}

	flag:=rand.Float32()<0.5
	if flag==keybord{
		return  USBProtocol{deviceFlag:flag,keyEvent:key_list[rand.Intn(len(key_list))]}
		//keycode=key_list[rand.Intn(len(key_list))]
	}else {
		//mouse_event.mouseCode=mouse_list[rand.Intn(len(mouse_list))]
		//x,y:=robotgo.GetMousePos()
		//mouse_event.X=int16(rand.Intn(screenSizeX))
		//mouse_event.Y=int16(rand.Intn(screenSizeY))
		return  USBProtocol{deviceFlag:flag,
			mouseEvent:MouseEvent{move2X:int16(rand.Intn(screenX)),
				move2Y:int16(rand.Intn(screenY)),
				mouseCode:mouse_list[rand.Intn(len(mouse_list))],
				screenSizeX:int16(screenX),
				screenSizeY:int16(screenY)}}
	}

}

func registerEvent(usbEvent USBProtocol, isTriggered chan bool) {

	isEnd:=false
	if usbEvent.deviceFlag==keybord{
		log.Printf("Register key event,keycode : %+v",key_map[usbEvent.keyEvent])
		hook.Register(hook.KeyHold, []string{key_map[usbEvent.keyEvent]}, func(e hook.Event) {
			isTriggered<-true
			isEnd=true
		})

	}else if usbEvent.deviceFlag==monse{
		//newX:=usbEvent.mouseEvent.currentX + int16(usbEvent.mouseEvent.offsetX)
		//newY:=usbEvent.mouseEvent.currentY + int16(usbEvent.mouseEvent.offsetY)
		me:=hook.Event{
			Kind:      mouse_kind[usbEvent.mouseEvent.mouseCode],
			When:      time.Time{},
			Mask:      0,
			Reserved:  0,
			Keycode:   0,
			Rawcode:   0,
			Keychar:   0,
			Button:   uint16(usbEvent.mouseEvent.mouseCode),
			Clicks:    0,
			X:         usbEvent.mouseEvent.move2X,
			Y:         usbEvent.mouseEvent.move2Y,
			Amount:    0,
			Rotation:  0,
			Direction: 0,
		}

		hook.RegisterMouse(me,func(e hook.Event){
			isTriggered<-true
			isEnd=true
		})

		log.Printf("Register Mouse Event : %+v",usbEvent.mouseEvent)

		//ct := false
		//timeup:=false
		//go func() {
		//	time.Sleep(time.Second)
		//	//log.Printf("event : %+v is no Triggered",usbEvent)
		//	if !isEnd{
		//		isTriggered<-false
		//		hook.End()
		//		timeup=true
		//	}
		//}()
		//
		//for !timeup{
		//	e := <-s
		//	if usbEvent.mouseEvent.mouseCode==MBTN {
		//		if e.Kind == hook.WheelDown && e.Button == ukey {
		//			isTriggered<-true
		//			isEnd=true
		//			hook.End()
		//
		//			break
		//		}
		//
		//	} else  {
		//		if e.Kind == hook.MouseMove && e.X == newX&& e.Y == newY {
		//			ct = true
		//		}
		//		if ct && e.Kind == hook.MouseDown && e.Button == ukey {
		//			isTriggered<-true
		//			isEnd=true
		//			hook.End()
		//			break
		//		}
		//	}
		//}
	}
	//定时
	//log.Println("unlock")
	time.Sleep(time.Second*2)
	if !isEnd{
		isTriggered<-false
	}
}

//发送随机键鼠事件并且监听检查
func SendRandomUsbEvent(exitch chan bool){
	var(
		usbTestCount int
		usbSuccessCount int
	)

	//初始化随机数种子
	rand.Seed(time.Now().Unix())

	log.Println("screenSize X : ",screenX," screenSize Y : ",screenY)
	hookEvent:=hook.Start()
	go func() {
		<-hook.Process(hookEvent)
	}()

	for{
		//初始化一个usb事件通知通道
		var usbCh = make(chan bool)
		usbEvent:=getRandomUsbEvent()
		//等待usb事件触发

		go registerEvent(usbEvent,usbCh)

		time.Sleep(time.Millisecond*100)
		//发送给树莓派,一个键鼠事件
		sendUsbEvent2Box(usbEvent)

		select {
		//若被通知退出,那么直接return
		case <- exitch:
			log.Println("SendRandomUsbEvent exit.")
			goto end
			//若还没退出就去查看usb事件是否触发
		default:{
			if <-usbCh {
				log.Printf("usbEvent=%+v is triggered\n\n", usbEvent)
				usbSuccessCount=usbSuccessCount+1
			} else{
				log.Printf("usbEvent=%+v is no triggered\n\n", usbEvent)

			}
			hook.UnRegister()
		}
		}
		usbTestCount=usbTestCount+1
	}
	end:
	hook.UnRegister()
	hook.End()
	usbAccuracyRate=float64(usbSuccessCount)/float64(usbTestCount)
}


