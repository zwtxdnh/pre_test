package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)


//等待一段时间后,执行结束任务的函数,可以指定并行任务
func WaitForKill(simultaneousTask func(chan bool),t int,kill func()){
	var taskch = make(chan bool)
	if simultaneousTask != nil{
		//执行异步任务
		go simultaneousTask(taskch)
	}
	//log.Println("wait time : ",t," Minute")
	time.Sleep(time.Minute*time.Duration(t))

	if simultaneousTask!=nil{
		taskch<-true
	}

	log.Println("run kill func")
	kill()

}

func CreateDir(path string){
	if _, err := os.Stat(path); err != nil {
		//fmt.Println("path not exists ", path)
		err := os.MkdirAll(path, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return
		}
	}
}

func setLogOutput(path string)*os.File{
	logFile, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout,logFile)
	log.SetOutput(mw)
	return logFile
}

func UTF16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

const (
	unknownEndian = iota
	bigEndian
	littleEndian


)

func ScanUTF16LinesFunc(byteOrder binary.ByteOrder) (bufio.SplitFunc, func() binary.ByteOrder) {

	// Function closure variables
	var endian = unknownEndian
	switch byteOrder {
	case binary.BigEndian:
		endian = bigEndian
	case binary.LittleEndian:
		endian = littleEndian
	}
	const bom = 0xFEFF
	var checkBOM bool = endian == unknownEndian

	// Scanner split function
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if checkBOM {
			checkBOM = false
			if len(data) > 1 {
				switch uint16(bom) {
				case uint16(data[0])<<8 | uint16(data[1]):
					endian = bigEndian
					return 2, nil, nil
				case uint16(data[1])<<8 | uint16(data[0]):
					endian = littleEndian
					return 2, nil, nil
				}
			}
		}

		// Scan for newline-terminated lines.
		i := 0
		for {
			j := bytes.IndexByte(data[i:], '\n')
			if j < 0 {
				break
			}
			i += j
			switch e := i % 2; e {
			case 1: // UTF-16BE
				if endian != littleEndian {
					if i > 1 {
						if data[i-1] == '\x00' {
							endian = bigEndian
							// We have a full newline-terminated line.
							return i + 1, dropCRBE(data[0 : i-1]), nil
						}
					}
				}
			case 0: // UTF-16LE
				if endian != bigEndian {
					if i+1 < len(data) {
						i++
						if data[i] == '\x00' {
							endian = littleEndian
							// We have a full newline-terminated line.
							return i + 1, dropCRLE(data[0 : i-1]), nil
						}
					}
				}
			}
			i++
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			// drop CR.
			advance = len(data)
			switch endian {
			case bigEndian:
				data = dropCRBE(data)
			case littleEndian:
				data = dropCRLE(data)
			default:
				data, endian = dropCR(data)
			}
			if endian == unknownEndian {
				if runtime.GOOS == "windows" {
					endian = littleEndian
				} else {
					endian = bigEndian
				}
			}
			return advance, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}

	// Endian byte order function
	orderFunc := func() (byteOrder binary.ByteOrder) {
		switch endian {
		case bigEndian:
			byteOrder = binary.BigEndian
		case littleEndian:
			byteOrder = binary.LittleEndian
		}
		return byteOrder
	}

	return splitFunc, orderFunc
}

// dropCREndian drops a terminal \r from the endian data.
func dropCREndian(data []byte, t1, t2 byte) []byte {
	if len(data) > 1 {
		if data[len(data)-2] == t1 && data[len(data)-1] == t2 {
			return data[0 : len(data)-2]
		}
	}
	return data
}

// dropCRBE drops a terminal \r from the big endian data.
func dropCRBE(data []byte) []byte {
	return dropCREndian(data, '\x00', '\r')
}

// dropCRLE drops a terminal \r from the little endian data.
func dropCRLE(data []byte) []byte {
	return dropCREndian(data, '\r', '\x00')
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) ([]byte, int) {
	var endian = unknownEndian
	switch ld := len(data); {
	case ld != len(dropCRLE(data)):
		endian = littleEndian
	case ld != len(dropCRBE(data)):
		endian = bigEndian
	}
	return data, endian
}

func USBStructToBytes(u *USBProtocol) []byte {
	var x reflect.SliceHeader
	x.Len =  int(unsafe.Sizeof(USBProtocol{}))
	x.Cap = int(unsafe.Sizeof(USBProtocol{}))
	x.Data = uintptr(unsafe.Pointer(u))
	return *(*[]byte)(unsafe.Pointer(&x))
}


const (
	REBOOT = 230
	SHUTDOWN  = 231
)
type Error struct {
	ErrorID   int    `db:"err_id"`
	Timestamp string `db:"err_timestamp"`
	IP      string `db:"err_ip"`
	Msg    string `db:"err_msg"`
}

var G_TestConfig TestConfig
//测试配置
type TestConfig struct {
	CloudHostTest CloudHostTest `json:"cloud_host_test"`
	RebootTest RebootTest `json:"reboot_test"`
}
type GameCase struct {
	Active bool `json:"active"`
	Time int `json:"time"`
	Order int `json:"order"`
}
type VideoCase struct {
	Active bool `json:"active"`
	Time int `json:"time"`
	Order int `json:"order"`
}
type CloudHostTest struct {
	GameCase GameCase `json:"game_case"`
	VideoCase VideoCase `json:"video_case"`
	Loop int `json:"loop"`
}
type RebootTest struct {
	ID int `json:"id"`
	Active bool `json:"active"`
	TotalReboot int `json:"total_reboot"`
	AlreadyReboot int `json:"already_reboot"`
}


//测试结果

type TestResult struct {
	IP string `db:"res_ip"`
	TimeStamp string `db:"res_timestamp"`
	Loss1 int `db:"res_lost1"`
	Loss2 int `db:"res_lost2"`
	Loss3 int `db:"res_lost3"`
	Loss4 int `db:"res_lost4"`
	Loss5 int `db:"res_lost5"`
	LossAbove5 int `db:"res_lost_upon5"`
	NetSlow int `db:"res_timeout"`
	DecodeSlow int `db:"res_decode_timeout"`
	EncodeSlow int `db:"res_encode_timeout"`
	UsbResult string `db:"res_usb_test"`
	RebootSuccessRate string `db:"res_reboot_count"`
}


type HostConfig struct {
	URL string `json:"url"`
	MarkPath string `json:"markPath"`
	WorkLogPath string `json:"workLogPath"`
}


////#include<sw_cpc_sdk.h>


//var swdll *syscall.DLL

//func initbox(){
//
//	init_box,_:=swdll.FindProc("init_box")
//	init_box.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr("./res/desktop_client/desktop_client.exe"))))
//}
//
//func startstream(){
//	start_stream_box,_:=swdll.FindProc("start_stream_box")
//	var sbc C.StartBoxCmd
//	sbc.ip=C.CString("192.168.29.194")
//	sbc.port=C.ushort(8800)
//	sbc.token=C.CString("123")
//	sbc.business_token=unsafe.Pointer(syscall.StringBytePtr("123"))
//	sbc.business_token_size=C.uint(3)
//	sbc.access=C.uint(0xFFFF)
//	log.Println("call start_stream_box")
//	start_stream_box.Call(uintptr(unsafe.Pointer(&sbc)))
//
//}
