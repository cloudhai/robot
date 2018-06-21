package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"log"
)

var upgrade = websocket.Upgrader{}
var realPath = "./static/"
func main()  {
	http.HandleFunc("/ws",echo)
	http.HandleFunc("/",home)
	err := http.ListenAndServe(":8800",nil)
	if err != nil{
		log.Fatal("listen err:",err)
	}
	fmt.Printf("finish")
}

func home(rw http.ResponseWriter,req *http.Request){
	path := req.URL.Path
	requstType := path[strings.LastIndex(path,"."):]
	switch requstType{
	case ".css":
		rw.Header().Set("Content-type","text/css")
	case ".js":
		rw.Header().Set("Content-type","text/javascript")
	default:

	}
	fin,err := os.Open(realPath+path)
	defer fin.Close()
	if err != nil {
		log.Println("static resource:",err)
		http.NotFound(rw,req)
	}
	fd,_ := ioutil.ReadAll(fin)
	rw.Write(fd)
}

func echo(rw http.ResponseWriter,req *http.Request){
	req.ParseForm()
	token := req.Form.Get("token")
	fmt.Printf(token)
	if token == ""{
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	conn,err := upgrade.Upgrade(rw,req,nil)
	if err != nil{
		log.Panic(err)
	}
	fmt.Println("get conn")
	//defer conn.Close()
	go func(){
		for {
			mt,msg,err := conn.ReadMessage()
			if err != nil{
				break
				log.Panic(err)
			}
			fmt.Printf("msg type：%d",mt)
			fmt.Println(string(msg))
			conn.WriteMessage(mt,[]byte("msg"))
		}
	}()
}
