package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MsgType byte
const(
	TYPE_PING MsgType = iota
	TYPE_PONG
	TYPE_CMD
	TYPE_MSG
)
const(
	HEADERLEN = 9
	MAGIC_CODE = 0x314518
)

type McuMsgHeader struct {
	MagicCode uint32
	Type		MsgType
	Length		uint32
}

type McuMsg struct{
	Header *McuMsgHeader
	Data []byte
}

func NewMsg(msgType MsgType,data string) McuMsg{
	return McuMsg{&McuMsgHeader{MagicCode:MAGIC_CODE,Type:msgType,Length:uint32(len(data))},[]byte(data)}
}

func (msg *McuMsg) Enpack() []byte{
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf,binary.LittleEndian,msg.Header.MagicCode)
	binary.Write(b_buf,binary.LittleEndian,msg.Header.Type)
	binary.Write(b_buf,binary.LittleEndian,uint32(len(msg.Data)))
	binary.Write(b_buf,binary.LittleEndian,msg.Data)
	return b_buf.Bytes()
}

func Depack(buf *[]byte) *McuMsg{
	header := new(McuMsgHeader)
	headerBuf := bytes.NewBuffer((*buf)[0:HEADERLEN])
	binary.Read(headerBuf,binary.LittleEndian,header)
	if verifyHeader(header){
		if uint32(len(*buf)) > (header.Length+HEADERLEN){
			data := (*buf)[HEADERLEN:header.Length+HEADERLEN]
			*buf = (*buf)[header.Length+HEADERLEN:]
			return &McuMsg{Header:header,Data:data}
		}else{
			return nil
		}
	}else{
		*buf = nil
		fmt.Println("非法数据")

	}
	return nil
}


func verifyHeader(header *McuMsgHeader) bool{
	if header == nil {
		return false
	}
	if MAGIC_CODE == header.MagicCode{
		return true
	}
	return false
}


