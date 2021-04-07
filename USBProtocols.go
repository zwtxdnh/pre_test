package main

var key_map map[uint8]string= map[uint8]string{
	KEY_1:"1",
	KEY_2:"2",
	KEY_3:"3",
	KEY_4:"4",
	KEY_5:"5",
	KEY_6:"6",
	KEY_7:"7",
	KEY_8:"8",
	KEY_9:"9",
	KEY_0:"0",
	//
	KEY_Q:"q",
	KEY_W:"w",
	KEY_E:"e",
	KEY_R:"r",
	KEY_T:"t",
	KEY_Y:"y",
	KEY_U:"u",
	KEY_I:"i",
	KEY_O:"o",
	KEY_P:"p",
	//
	KEY_A:"a",
	KEY_S:"s",
	KEY_D:"d",
	KEY_F:"f",
	KEY_G:"g",
	KEY_H:"h",
	KEY_J:"j",
	KEY_K:"k",
	KEY_L:"l",

	//
	KEY_Z:"z",
	KEY_X:"x",
	KEY_C:"c",
	KEY_V:"v",
	KEY_B:"b",
	KEY_N:"n",
	KEY_M:"m",
	//
	KEY_F1:"f1",
	KEY_F2:"f2",
	KEY_F3:"f3",
	KEY_F4:"f4",
	KEY_F5:"f5",
	KEY_F6:"f6",
	KEY_F7:"f7",
	KEY_F8:"f8",
	KEY_F9:"f9",
	KEY_F10:"f10",
	KEY_F11:"f11",
	KEY_F12:"f12",
	// more
	//:"esc",


	KEY_DELETE:"delete",
	KEY_TAB:"tab",
	KEY_RIGHTCTRL:"ctrl",
	KEY_LEFTCTRL:"ctrl",
	KEY_SPACE:"space",
	KEY_RIGHTSHIFT:"shift",
	KEY_LEFTSHIFT:"shift",
	KEY_ENTER:"enter",
	KEY_LEFTALT:"alt",
	KEY_RIGHTALT:"alt",

	//
	KEY_UP:"up",
	KEY_DOWN:"down",
	KEY_LEFT:"left",
	KEY_RIGHT:"right",

}
//var mouse_map map[uint8]string= map[uint8]string{
//	LBTN:"mleft",
//	RBTN:"mright",
//	MBTN:"mcenter",
//}
var mouse_kind map[uint8]uint8= map[uint8]uint8{
	LBTN: 8,
	RBTN: 8,
	0:9,
}

var mouse_list = []uint8{LBTN,RBTN}

var key_list=[]uint8{KEY_1, KEY_2, KEY_3, KEY_4, KEY_5, KEY_6, KEY_7, KEY_8, KEY_9, KEY_0,
	KEY_A,KEY_B,KEY_C,KEY_D,KEY_E,KEY_F,KEY_G,
	KEY_H,KEY_I,KEY_J,KEY_K,KEY_L,KEY_M,KEY_N,
	KEY_O,KEY_P,/*KEY_Q*/KEY_R,KEY_S,KEY_T,KEY_U,
	KEY_V,KEY_W,KEY_X,KEY_Y,KEY_Z,KEY_TAB,KEY_SPACE,KEY_F1,KEY_F2,KEY_F3,
	KEY_F4,KEY_F5,KEY_F6,KEY_F7,KEY_F8,KEY_F9,KEY_F10,KEY_F11,KEY_F12,
	KEY_INSERT,KEY_HOME,KEY_PAGEUP,KEY_DELETE,KEY_END,KEY_PAGEDOWN,KEY_RIGHT,
	KEY_LEFT,KEY_DOWN,KEY_UP,KEY_ENTER }


var num_key_list = []uint8{KEY_1, KEY_2, KEY_3, KEY_4, KEY_5, KEY_6, KEY_7, KEY_8, KEY_9, KEY_0}

var all_char_list = []uint8{KEY_A,KEY_B,KEY_C,KEY_D,KEY_E,KEY_F,KEY_G,
	KEY_H,KEY_I,KEY_J,KEY_K,KEY_L,KEY_M,KEY_N,
	KEY_O,KEY_P,KEY_Q,KEY_R,KEY_S,KEY_T,KEY_U,
	KEY_V,KEY_W,KEY_X,KEY_Y,KEY_Z}

var extend_key_list = []uint8{KEY_TAB,KEY_SPACE,KEY_F1,KEY_F2,KEY_F3,
	KEY_F4,KEY_F5,KEY_F6,KEY_F7,KEY_F8,KEY_F9,KEY_F10,KEY_F11,KEY_F12,
	KEY_INSERT,KEY_HOME,KEY_PAGEUP,KEY_DELETE,KEY_END,KEY_PAGEDOWN,KEY_RIGHT,
	KEY_LEFT,KEY_DOWN,KEY_UP,KEY_ENTER ,KEY_LEFTSHIFT,KEY_LEFTCTRL ,
	KEY_LEFTALT ,KEY_RIGHTCTRL ,KEY_RIGHTSHIFT,KEY_RIGHTALT}

const keybord bool =true
const monse bool  = false

type USBProtocol struct {
	deviceFlag bool //true=key/false=mouse
	keyEvent uint8
	mouseEvent MouseEvent
}
type MouseEvent struct {
	mouseCode uint8//若是0就是只移动
	move2X int16
	move2Y int16
	screenSizeX int16
	screenSizeY int16
}

//mouse
const LBTN uint8= 0x1
const RBTN uint8= 0x2
const MBTN uint8= 0x4
//key
const  KEY_MOD_LCTRL  = 0x01
const  KEY_MOD_LSHIFT = 0x02
const KEY_MOD_LALT =  0x04
const KEY_MOD_LMETA = 0x08
const KEY_MOD_RCTRL = 0x10
const KEY_MOD_RSHIFT = 0x20
const KEY_MOD_RALT =  0x40
const KEY_MOD_RMETA = 0x80

const KEY_NONE = 0x00
const KEY_ERR_OVF = 0x01
const KEY_A = 0x04
const KEY_B = 0x05
const KEY_C = 0x06
const KEY_D = 0x07
const KEY_E = 0x08
const KEY_F = 0x09
const KEY_G = 0x0a
const KEY_H = 0x0b
const KEY_I = 0x0c
const KEY_J = 0x0d
const KEY_K = 0x0e
const KEY_L = 0x0f
const KEY_M = 0x10
const KEY_N = 0x11
const KEY_O = 0x12
const KEY_P = 0x13
const KEY_Q = 0x14
const KEY_R = 0x15
const KEY_S = 0x16
const KEY_T = 0x17
const KEY_U = 0x18
const KEY_V = 0x19
const KEY_W = 0x1a
const KEY_X = 0x1b
const KEY_Y = 0x1c
const KEY_Z = 0x1d

const KEY_1 = 0x1e
const KEY_2 = 0x1f
const KEY_3 = 0x20
const KEY_4 = 0x21
const KEY_5 = 0x22
const KEY_6 = 0x23
const KEY_7 = 0x24
const KEY_8 = 0x25
const KEY_9 = 0x26
const KEY_0 = 0x27

const KEY_ENTER = 0x28
const KEY_ESC = 0x29
const KEY_BACKSPACE = 0x2a
const KEY_TAB = 0x2b
const KEY_SPACE = 0x2c
const KEY_MINUS = 0x2d
const KEY_EQUAL = 0x2e
const KEY_LEFTBRACE = 0x2f
const KEY_RIGHTBRACE = 0x30
const KEY_BACKSLASH = 0x31
const KEY_HASHTILDE = 0x32
const KEY_SEMICOLON = 0x33
const KEY_APOSTROPHE = 0x34
const KEY_GRAVE = 0x35
const KEY_COMMA = 0x36
const KEY_DOT = 0x37
const KEY_SLASH = 0x38
const KEY_CAPSLOCK = 0x39

const KEY_F1 = 0x3a
const KEY_F2 = 0x3b
const KEY_F3 = 0x3c
const KEY_F4 = 0x3d
const KEY_F5 = 0x3e
const KEY_F6 = 0x3f
const KEY_F7 = 0x40
const KEY_F8 = 0x41
const KEY_F9 = 0x42
const KEY_F10 = 0x43
const KEY_F11 = 0x44
const KEY_F12 = 0x45

const KEY_SYSRQ = 0x46
const KEY_SCROLLLOCK = 0x47
const KEY_PAUSE = 0x48
const KEY_INSERT = 0x49
const KEY_HOME = 0x4a
const KEY_PAGEUP = 0x4b
const KEY_DELETE = 0x4c
const KEY_END = 0x4d
const KEY_PAGEDOWN = 0x4e
const KEY_RIGHT = 0x4f
const KEY_LEFT = 0x50
const KEY_DOWN = 0x51
const KEY_UP = 0x52

const KEY_NUMLOCK = 0x53
const KEY_KPSLASH = 0x54
const KEY_KPASTERISK = 0x55
const KEY_KPMINUS = 0x56
const KEY_KPPLUS = 0x57
const KEY_KPENTER = 0x58
const KEY_KP1 = 0x59
const KEY_KP2 = 0x5a
const KEY_KP3 = 0x5b
const KEY_KP4 = 0x5c
const KEY_KP5 = 0x5d
const KEY_KP6 = 0x5e
const KEY_KP7 = 0x5f
const KEY_KP8 = 0x60
const KEY_KP9 = 0x61
const KEY_KP0 = 0x62
const KEY_KPDOT = 0x63

const KEY_102ND = 0x64
const KEY_COMPOSE = 0x65
const KEY_POWER = 0x66
const KEY_KPEQUAL = 0x67

const KEY_F13 = 0x68
const KEY_F14 = 0x69
const KEY_F15 = 0x6a
const KEY_F16 = 0x6b
const KEY_F17 = 0x6c
const KEY_F18 = 0x6d
const KEY_F19 = 0x6e
const KEY_F20 = 0x6f
const KEY_F21 = 0x70
const KEY_F22 = 0x71
const KEY_F23 = 0x72
const KEY_F24 = 0x73

const KEY_OPEN = 0x74
const KEY_HELP = 0x75
const KEY_PROPS = 0x76
const KEY_FRONT = 0x77
const KEY_STOP = 0x78
const KEY_AGAIN = 0x79
const KEY_UNDO = 0x7a
const KEY_CUT = 0x7b
const KEY_COPY = 0x7c
const KEY_PASTE = 0x7d
const KEY_FIND = 0x7e
const KEY_MUTE = 0x7f
const KEY_VOLUMEUP = 0x80
const KEY_VOLUMEDOWN = 0x81

const KEY_KPCOMMA = 0x85

const KEY_RO = 0x87
const KEY_KATAKANAHIRAGANA = 0x88
const KEY_YEN = 0x89
const KEY_HENKAN = 0x8a
const KEY_MUHENKAN = 0x8b
const KEY_KPJPCOMMA = 0x8c

const KEY_HANGEUL = 0x90
const KEY_HANJA = 0x91
const KEY_KATAKANA = 0x92
const KEY_HIRAGANA = 0x93
const KEY_ZENKAKUHANKAKU = 0x94
const KEY_KPLEFTPAREN = 0xb6
const KEY_KPRIGHTPAREN = 0xb7
const KEY_LEFTCTRL = 0xe0
const KEY_LEFTSHIFT = 0xe1
const KEY_LEFTALT = 0xe2
const KEY_LEFTMETA = 0xe3
const KEY_RIGHTCTRL = 0xe4
const KEY_RIGHTSHIFT = 0xe5
const KEY_RIGHTALT = 0xe6
const KEY_RIGHTMETA = 0xe7

const KEY_MEDIA_PLAYPAUSE = 0xe8
const KEY_MEDIA_STOPCD = 0xe9
const KEY_MEDIA_PREVIOUSSONG = 0xea
const KEY_MEDIA_NEXTSONG = 0xeb
const KEY_MEDIA_EJECTCD = 0xec
const KEY_MEDIA_VOLUMEUP = 0xed
const KEY_MEDIA_VOLUMEDOWN = 0xee
const KEY_MEDIA_MUTE = 0xef
const KEY_MEDIA_WWW = 0xf0
const KEY_MEDIA_BACK = 0xf1
const KEY_MEDIA_FORWARD = 0xf2
const KEY_MEDIA_STOP = 0xf3
const KEY_MEDIA_FIND = 0xf4
const KEY_MEDIA_SCROLLUP = 0xf5
const KEY_MEDIA_SCROLLDOWN = 0xf6
const KEY_MEDIA_EDIT = 0xf7
const KEY_MEDIA_SLEEP = 0xf8
const KEY_MEDIA_COFFEE = 0xf9
const KEY_MEDIA_REFRESH = 0xfa
const KEY_MEDIA_CALC = 0xfb