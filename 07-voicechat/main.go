// WebRTCでボイスチャット作る
package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	TYPE_CONNECTED = "connected"
	TYPE_ENTER     = "enter"
	TYPE_LEAVE     = "leave"
)

func main() {
	http.HandleFunc("/", htmlHandler)

	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("templates/resource"))))

	/*
			conns = make(map[string]*Connection)
		http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(wsHandler)}
			s.ServeHTTP(w, req)
		})
		if err := http.ListenAndServe(SERVE_WS, nil); err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	*/

	conns = make(map[string]*Connection)
	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(wsHandler)}
		s.ServeHTTP(w, req)
	})

	// サーバーを起動
	// http://localhost:8989
	fmt.Println("The Server runs with http://localhost:8989")
	http.ListenAndServe(":8989", nil)
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {

	// ホスト名の取得
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// テンプレートをパース
	t := template.Must(template.ParseFiles("templates/template.html"))

	// HTMLに入れるオブジェクト
	dat := struct {
		HostName string
	}{
		HostName: name,
	}

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "template.html", dat); err != nil {
		log.Fatal(err)
	}
}

type Connection struct {
	Channel     string
	Conn        *websocket.Conn
	Room        string
	sendChannel chan Message
}

func NewConnection(ws *websocket.Conn) *Connection {
	var id, _ = uuid.NewV4()
	channel := id.String()
	connection := &Connection{
		Channel:     channel,
		Conn:        ws,
		sendChannel: make(chan Message)}
	go channelHandler(connection)
	msg := Message{
		Type: TYPE_CONNECTED,
		Msg:  channel}
	connection.SendMessage(msg, connection)
	return connection
}

func channelHandler(self *Connection) {
	for msg := range self.sendChannel {
		fmt.Printf("Sending %s channel:%s\n", msg.Type, self.Channel)
		websocket.JSON.Send(self.Conn, msg)
	}
}

func (self *Connection) EnterRoom(roomname string) {
	self.Room = roomname
}

func (self *Connection) LeaveRoom(roomname string) {
	self.Room = ""
}
func (self *Connection) IsJoiningRoom(targetRoom string) bool {
	return self.Room == targetRoom
}

func (self *Connection) SendMessage(msg Message, from *Connection) {
	msg.From = from.Channel
	self.sendChannel <- msg
}

type Message struct {
	Type string
	Msg  string
	From string
}

const (
	DestTypeUnicast   = "uni"
	DestTypeBroadcast = "bro"
	DestTypeRoom      = "room"
)

type MessageFrame struct {
	Type    string
	Dest    string
	Message Message
}

var conns map[string]*Connection

func wsHandler(ws *websocket.Conn) {
	conn := NewConnection(ws)
	conns[conn.Channel] = conn
	wsMsgHandler(conn)
}

func wsMsgHandler(conn *Connection) {
	for {
		var frame MessageFrame
		err := websocket.JSON.Receive(conn.Conn, &frame)
		if err != nil {
			fmt.Println(err)
			onDisconnected(conn)
			return
		}

		fmt.Printf("Received: %s from[%s]\n", frame.Type, conn.Channel)

		if frame.Message.Type == TYPE_ENTER {
			onEnter(conn, frame)
		} else {
			frame.Message.From = conn.Channel
			switch frame.Type {
			case DestTypeUnicast:
				conns[frame.Dest].SendMessage(frame.Message, conn)
			default:
				broadcastRoom(frame.Message, frame.Dest, conn)
			}
		}
	}
}

func onEnter(conn *Connection, frame MessageFrame) {
	fmt.Printf("channel[%s] entered room[%s]\n", conn.Channel, frame.Dest)
	conn.EnterRoom(frame.Dest)
	broadcastRoom(frame.Message, frame.Dest, conn)
}

func onDisconnected(conn *Connection) {
	fmt.Printf("channel[%s] dicconnected\n", conn.Channel)
	delete(conns, conn.Channel)
	msg := Message{
		Type: TYPE_LEAVE,
		Msg:  conn.Channel,
		From: conn.Channel,
	}
	broadcastRoom(msg, conn.Room, conn)
}

func broadcastRoom(msg Message, room string, from *Connection) {
	fmt.Printf("------ broadcast room:[%s] from:[%s] type[%s] --------\n", room, from.Channel, msg.Type)
	for _, conn := range conns {
		if conn.IsJoiningRoom(room) {
			conn.SendMessage(msg, from)
		}
	}
}
