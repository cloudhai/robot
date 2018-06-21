package main

import (
	"net"
	"fmt"
	"os"
	server2 "robot/server"
)

type Info struct{
	Message string "msg"

}

func main(){

	//msg := &entity.CmdMsg{
	//	Magic:0x9023ae43,
	//	Type:entity.MsgType_CMD,
	//	Data:[]byte("this is test foarmat"),
	//}
	//msg := &protocol.McuMsg{Header:&protocol.McuMsgHeader{MagicCode:protocol.MAGIC_CODE,Tyte:protocol.CMD,},Data:[]byte("1234567890")}
	//msg2 := &protocol.McuMsg{Header:&protocol.McuMsgHeader{MagicCode:protocol.MAGIC_CODE,Tyte:protocol.MSG,},Data:[]byte("this is a text for test of protocol")}
	//buf := msg.Enpack()
	//buf = append(buf,msg2.Enpack()...)
	//ioutil.WriteFile("msg.dat",buf,os.ModePerm)

	//data,_ := ioutil.ReadFile("./msg.dat")
	//msg:= protocol.Depack(&data)
	//if msg != nil{
	//	fmt.Println(string(msg.Data))
	//}
	//if len(data) > protocol.HEADERLEN{
	//	msg2:= protocol.Depack(&data)
	//	if msg2 != nil{
	//		fmt.Println(string(msg2.Data))
	//	}
	//}
	//fmt.Printf("magic code is :%x",msg.Magic)
	svr,err := net.ResolveTCPAddr("tcp","127.0.0.1:8322")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	listen,err := net.ListenTCP("tcp",svr)
	if err != nil{
		fmt.Println(err)
	}
	defer listen.Close()
	for{
		conn,err := listen.Accept()
		if err != nil{
			fmt.Println(err)
		}
		server := server2.NewConnection(conn,0xff)
		server.Run()
	}
}
