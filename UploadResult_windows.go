package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getWorkLog(logdir string)*os.File{
	files, _ := ioutil.ReadDir(logdir)
	for _, f := range files {

		if strings.Contains(f.Name(),"desktop_worker"){
			log.Println("find worker log name : ",f.Name())
			fi, err := os.OpenFile(logdir+"/"+f.Name() ,os.O_RDONLY, 0766)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			//temp, _ := ioutil.ReadAll(fi)
			//utf8String:=UTF16BytesToString(temp,binary.LittleEndian)
			////log.Println(utf8String)
			//fi.Close()
			//fi, err = os.OpenFile(logdir+"/"+f.Name() ,os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
			//n, _ := fi.Seek(0, os.SEEK_END)
			//_, err = fi.WriteAt([]byte(utf8String), n)
			return fi
		}
	}
	return nil
}

func uploadResult(){
	var result TestResult

	fi:=getWorkLog(hostConfig.WorkLogPath)
	defer fi.Close()
	if fi==nil{
		log.Fatal("open desktop_worker log fail")
	}

	buf := bufio.NewScanner(fi)
	var bo binary.ByteOrder
	bo=binary.LittleEndian
	splitFunc, orderFunc := ScanUTF16LinesFunc(bo)
	buf.Split(splitFunc)
	for {
		if !buf.Scan() {
			break
		}
		b := buf.Bytes()

		line := UTF16BytesToString(b, orderFunc())
		//line := strings.TrimSpace(string(bline))
		//log.Println(string(bline))
		//line=UTF16BytesToString(line,binary.LittleEndian)
		//log.Println(line)
		encodeReg:=regexp.MustCompile(`encode\[avg:\d+us, max:\d+us, slow:(\d+), total_slow:\d+\], decode`)
		encodeData:=encodeReg.FindAllStringSubmatch(line,-1)
		if encodeData != nil{
			//fmt.Printf("%+v\n",encodeData[0][1])
			slow,_:=strconv.Atoi(encodeData[0][1])
			result.EncodeSlow+=slow
			//continue
		}
		decodeReg:=regexp.MustCompile(`decode\[avg:\d+us, max:\d+us, slow:(\d+), total_slow:\d+\], net`)
		decodeData:=decodeReg.FindAllStringSubmatch(line,-1)
		if decodeData != nil{
			//fmt.Printf("%+v\n",decodeData)
			slow,_:=strconv.Atoi(decodeData[0][1])
			result.DecodeSlow+=slow
			//continue
		}
		netReg:=regexp.MustCompile(`net\[avg:\d+us, max:\d+us, slow:(\d+), total_slow:\d+], delay`)
		netData:=netReg.FindAllStringSubmatch(line,-1)
		if netData != nil{
			slow,_:=strconv.Atoi(netData[0][1])
			result.NetSlow+=slow
			//continue
		}
		fpsReg:=regexp.MustCompile(`capture_frame:(\d+), play_stream:(\d+), render:(\d+)`)
		fpsData:=fpsReg.FindAllStringSubmatch(line,-1)
		if fpsData != nil{
			capture,_:=strconv.Atoi(fpsData[0][1])
			render,_:=strconv.Atoi(fpsData[0][3])
			fpsloss:=capture-render

			if fpsloss==1{
				result.Loss1+=1
			}else if fpsloss==2{
				result.Loss2+=1
			}else if fpsloss==3{
				result.Loss3+=1
			}else if fpsloss==4{
				result.Loss4+=1
			}else if fpsloss==5{
				result.Loss5+=1
			}else if fpsloss>5{
				result.LossAbove5+=1
			}

			//log.Printf("%+v\n",fpsData)
			//continue
		}
	}

	result.UsbResult=strconv.FormatFloat(usbAccuracyRate,'f',6,64)
	result.TimeStamp=time.Now().Format("2006/1/2 15:04:05")
	log.Printf("analyze log ended, result %+v: ",result)

	js,err :=json.Marshal(result)
	if err!=nil{
		log.Fatal(err.Error())
	}

	_,err = http.Post(hostConfig.URL+"/report/streamResult",
		"application/x-www-form-urlencoded",
		bytes.NewBuffer(js))
	if err != nil {
		log.Println(err)
	}

	log.Println("post result to server")
}