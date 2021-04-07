package main

/*
#include <windows.h>
#include<stdio.h>

#define WAIT_FOR_EVENT 30000//等待时间
FILE *fp = NULL;
BOOL check_video(){
	fp = fopen("./log/client.log", "a");
	BOOL flag = TRUE;
	HANDLE   client_video_ack_handle_ = OpenEventA(EVENT_ALL_ACCESS, FALSE, "Global\\swyun_client_video_ack_event");
	HANDLE   service_video_ok_handle_ = OpenEventA(EVENT_ALL_ACCESS, FALSE, "Global\\service_video_ok_event");
	HANDLE handle_array[2] = { client_video_ack_handle_, service_video_ok_handle_ };
	if(client_video_ack_handle_==NULL)
	{
   		fprintf(fp, "get swyun_client_video_ack_event fail\n");
	}
	if(service_video_ok_handle_==NULL)
	{
		fprintf(fp,"get service_video_ok_event fail\n");
	}
	DWORD ret = WaitForMultipleObjects(2, handle_array, TRUE, WAIT_FOR_EVENT);//第三变量表示是否等待全部事件,第四变量表示等待时间
	if (ret == WAIT_OBJECT_0)
	{
		flag = TRUE;
	}
	else
	{
	DWORD res = GetLastError();
	fprintf(fp,"stream video fail,GetLastError = %lu, error code = %lu\n",ret,res);
	flag = FALSE;
	}
	CloseHandle(client_video_ack_handle_);
	CloseHandle(service_video_ok_handle_);
	fclose(fp);
	return flag;

}

BOOL check_audio()
{
	BOOL flag = TRUE;
	//HANDLE   box_audio_handle_ = OpenEventA(EVENT_ALL_ACCESS, TRUE, "Global\\swyun_client_box_audio_event");
	HANDLE   pc_audio_handle_ = OpenEventA(EVENT_ALL_ACCESS, TRUE, "Global\\swyun_client_pc_audio_event");
	//HANDLE handle_array[2] = { box_audio_handle_, pc_audio_handle_ };
	if (pc_audio_handle_ == NULL)
	{
		printf("get pc_audio_handle_ fail\n");
	}

	DWORD ret = WaitForSingleObject(pc_audio_handle_, WAIT_FOR_EVENT);//第三变量表示是否等待全部事件,第四变量表示等待时间
	if (ret == WAIT_OBJECT_0 || (ret == WAIT_OBJECT_0 + 1))
	{
		flag = TRUE;
	}
	else
	{
		DWORD res = GetLastError();
		printf("stream audio fail,GetLastError = %lu, error code = %lu\n", ret, res);
		flag = FALSE;
	}

	CloseHandle(pc_audio_handle_);
	return flag;
}

*/
import "C"
import (
	"log"
	"net/http"
	"strings"
)

func check_video()bool{
	if C.check_video()==0{

		http.Post(hostConfig.URL+"/report/error",
			"application/x-www-form-urlencoded",
			strings.NewReader("error=check video fail"))

		log.Println("check video fail")

		return false
	}else {
		log.Println("check video ok")
		return true
	}
}

func check_audio()bool{

	if C.check_audio()==0{
		log.Fatal("check audio fail")
		return false
	}else {
		return true
	}

}