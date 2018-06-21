package server

import (
	"net"
	"robot/protocol"
	"fmt"
	"time"
	"io"
)

const(
	STATUS_NONE = iota
	STATUS_CONN
	STATUS_DISC
)

const(
	MAX_BUF_LENGTH = 1024
	SEND_TIMEOUT	= 2
)

type Connection struct {
	conn 			net.Conn
	status 			int
	connId 			int
	sendMsgQueue	chan []byte
	sendTimeOut		int
	eventQueue		string
	buf 			[]byte
}

func NewConnection(c net.Conn,sendBufSize int) *Connection{
	return &Connection{
		conn:c,
		status:STATUS_CONN,
		connId:0,
		sendMsgQueue:make(chan []byte,sendBufSize),

	}
}

func (this *Connection)Close(){
	if this.status != STATUS_CONN{
		return
	}
	this.conn.Close()
	this.status = STATUS_DISC
}

func (this *Connection) unpack(buf []byte){
	if buf != nil{
		this.buf = append(this.buf,buf...)
	}
	if len(this.buf) > 9{
		msg := protocol.Depack(&this.buf)
		if msg == nil{
			return
		}
		go func(msg *protocol.McuMsg){
			//handler the msg
			fmt.Println("handler the msg")
			if msg.Header.Type == protocol.TYPE_CMD || msg.Header.Type == protocol.TYPE_MSG{
				fmt.Println(string(msg.Data))
				this.Send(protocol.NewMsg(protocol.TYPE_CMD,"getip"))
			}

		}(msg)
		if len(this.buf) > protocol.HEADERLEN {
			this.unpack(nil)
		}
	}


}

func (this *Connection) Send(msg protocol.McuMsg){
	if this.status != STATUS_CONN{
		fmt.Println("status is not connect")
		return
	}
	buf := msg.Enpack()

	select {
		case this.sendMsgQueue <- buf:
		case <-time.After(time.Duration(SEND_TIMEOUT*time.Second)):
			this.Close()
	}
}

func (this *Connection) Run(){
	go func(){
		for{
			select {
				case msg,ok := <-this.sendMsgQueue:
					if ok{
						if this.status == STATUS_CONN{
							_,err := this.conn.Write(msg)
							if err != nil {

							}
						}
					}
			}
		}
	}()

	go func(){
		for{
			buf := make([]byte,MAX_BUF_LENGTH)
			//buf := []byte{}
			n,err := this.conn.Read(buf)
			if err == io.EOF{
				fmt.Println("连接关闭")
				this.Close()
				return
			}
			if err != nil{
				fmt.Println(err)
				return
			}
			if n > 0 {
				this.unpack(buf[:n])
			}
		}
	}()
}
