package main

import (
	"encoding/json"
	."github.com/CodyGuo/win"
	"github.com/cakturk/go-netstat/netstat"
	"github.com/go-vgo/robotgo"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

/*
#include <windows.h>
void reboot()
{
	ExitWindowsEx(EWX_REBOOT, 0);
}

*/
import "C"

var box *net.TCPConn

//bool check_video()
//{
//bool flag = true;
//HANDLE   client_video_ack_handle_ = ::OpenEvent(EVENT_ALL_ACCESS, FALSE, _T("Global\\swyun_client_video_ack_event"));
//HANDLE   service_video_ok_handle_ = ::OpenEvent(EVENT_ALL_ACCESS, FALSE, _T("Global\\service_video_ok_event"));
//HANDLE handle_array[2] = { client_video_ack_handle_, service_video_ok_handle_ };
//if(client_video_ack_handle_==NULL)
//{
//log_print(kError, _T("未获得视频ACK句柄"));
//}
//if(service_video_ok_handle_==NULL)
//{
//log_print(kError, _T("未获得视频发送句柄"));
//}
//DWORD ret = WaitForMultipleObjects(2, handle_array, TRUE, WAIT_FOR_EVENT);//第三变量表示是否等待全部事件,第四变量表示等待时间
//if (ret == WAIT_OBJECT_0)
//{
//flag = true;
//}
//else
//{
//DWORD res = GetLastError();
//log_print(kError, _T("串流视频模块出错了,返回值是%lu,错误码是%lu"),ret,res);
//flag = false;
//}
//CloseHandle(client_video_ack_handle_);
//CloseHandle(service_video_ok_handle_);
//return flag;
//}

func RunAllCase(){
	//todo 这里的代码要把程序改成服务运行才可以生效
	//go func() {
	//	for true {
	//		if !check_video(){
	//			break
	//		}
	//		time.Sleep(time.Second)
	//	}
	//
	//	syscall.Exit(1)
	//}()
	//执行测试用例
	taskmap:=map[int]func(){}
	taskmap[G_TestConfig.CloudHostTest.GameCase.Order]=gameCase
	taskmap[G_TestConfig.CloudHostTest.VideoCase.Order]=videoCase

	for i:=0;i<G_TestConfig.CloudHostTest.Loop;i++{
		for j:=0;j<len(taskmap);j++{
			log.Println("run test case ",j)
			taskmap[j]()
		}
	}
		//上报主机端测试结果
		uploadResult()

		//若重启测试激活,那么报告完后重启主机
		if G_TestConfig.RebootTest.Active{
			reboot()
		}else {
			//没激活测试完成直接关机
			shutdown()
		}

}

func Connect2Box() {

	for{
		tabs, _ := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
			return s.State == netstat.Established && s.LocalAddr.Port == 8800
		})
		if tabs != nil{
			//获取到盒子的ip
			addr, err := net.ResolveTCPAddr("tcp4", /**/tabs[0].RemoteAddr.IP.String()+":9111")
			checkError(err)
			log.Println("Try connect to USB ",addr)
			//连接盒子
			conn, err := net.DialTCP("tcp4", nil, addr)
			if err==nil{
				//连接上后才能退出
				log.Println("USB emulator is connected ")
				box=conn
				break
			}

		}else{
			log.Println("Box not connected")
			time.Sleep(time.Second)
		}
	}
	time.Sleep(time.Second*15)
}

//todo 实现游戏测试一定时间
//使用命令调起furemark，同时随机发送键鼠事件
func gameCase(){
	if !G_TestConfig.CloudHostTest.GameCase.Active{
		log.Println("pass gameCase")
		return
	}

	screenX, screenY =robotgo.GetScreenSize()
	log.Println("gameCase runing")
	//L" /nogui /fullscreen /width=" + to_wstring(GetSystemMetrics(SM_CXSCREEN)) +
	//	//	L" /height=" + to_wstring(GetSystemMetrics(SM_CYSCREEN)) + L" /noscore \
	//	///log_score /log_gpu_data_polling_factor=10 /run_mode=1 /max_time="
	//	//+ to_wstring(testConfig.GPUTestTIme * 60 * 1000)
	//	///nogui /width=1024 /height=728 /noscore /log_score /log_gpu_data_polling_factor=10 /run_mode=1 /max_time=6000
	//	//cmdline:=[]string{"-d=skydiver.3dmdef","-l=10000"}
	cmdline:=[]string{"/nogui",
		"/width="+strconv.Itoa(screenX),
		"/height="+strconv.Itoa(screenY),
		"/run_mode=1",
		"/noscore",
		"/fullscreen=1",
		"/max_time="+strconv.Itoa(G_TestConfig.CloudHostTest.GameCase.Time*60*1000)}
	log.Println()
	cmd:=exec.Command(hostConfig.MarkPath,cmdline...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Start()
	checkError(err)

	killfunc:=func(){log.Println("kill pid = "+strconv.Itoa(cmd.Process.Pid))
		log.Println("top windows : ",robotgo.GetTitle())
		robotgo.KeyTap("esc")
		_ = cmd.Process.Kill()}

	WaitForKill(SendRandomUsbEvent, G_TestConfig.CloudHostTest.GameCase.Time,killfunc)

	//_ = cmd.Wait()
	//log.Println(cmd.Output())
	log.Println("gameCase over")
	log.Println("usb accuracy rate : ",usbAccuracyRate)
	//time.Sleep(time.Duration(2)*time.Second)
	//cmd_in.WriteString("\"C:/Program Files/Futuremark/3DMark/3DMarkCmd.exe\" -d=skydiver.3dmdef -l=10000")

}

//todo 实现播放视频一定时间
func videoCase(){
	if !G_TestConfig.CloudHostTest.VideoCase.Active{
		log.Println("pass videoCase")
		return
	}
	log.Println("videoCase runing")
	cmdline:=[]string{"-loop","0","pause","../res/example.MP4","-fs"}
	cmd:=exec.Command("../res/mplayer/mplayer.exe",cmdline...)

	err := cmd.Start()
	checkError(err)
	killfunc:=func(){log.Println("kill pid = "+strconv.Itoa(cmd.Process.Pid))
		log.Println("top windows : ",robotgo.GetTitle())
		_ = cmd.Process.Kill()}
	WaitForKill(nil, G_TestConfig.CloudHostTest.VideoCase.Time, killfunc)

	//_ = cmd.Wait()
	//输出命令行运行结果
	//log.Println(cmd.Output())
	log.Println("videoCase over")
}


func upPrivileges(){
	var hToken HANDLE
	var tkp TOKEN_PRIVILEGES

	OpenProcessToken(GetCurrentProcess(), TOKEN_ADJUST_PRIVILEGES|TOKEN_QUERY, &hToken)
	LookupPrivilegeValueA(nil, StringToBytePtr(SE_SHUTDOWN_NAME), &tkp.Privileges[0].Luid)
	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED
	AdjustTokenPrivileges(hToken, false, &tkp, 0, nil, nil)
}

//重启
func reboot(){
	log.Println("could host reboot")
	upPrivileges()
	ExitWindowsEx(EWX_REBOOT, 0)
}
//关机
func shutdown(){
	log.Println("could host shutdown")
	upPrivileges()
	ExitWindowsEx(EWX_SHUTDOWN, 0)

}

//获取测试用例
func GetTestCase(url string){
	resp,err:=http.Get(url)
	checkError(err)
	log.Println("StatusCode : ",resp.StatusCode)
	defer resp.Body.Close()
	if resp.StatusCode==REBOOT{
		reboot()
		syscall.Exit(1)
	}else if resp.StatusCode==SHUTDOWN{
		shutdown()
		syscall.Exit(1)
	} else if resp.StatusCode==http.StatusOK{

		body, err := ioutil.ReadAll(resp.Body)
		checkError(err)
		//log.Println(string(body))
		err = json.Unmarshal(body, &G_TestConfig)
		log.Printf("testConfig=%+v\n", G_TestConfig)
	}

}

