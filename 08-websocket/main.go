// websocketを使ってServer->Client通信やる
// https://qiita.com/vitor/items/4a257cc24f6a07e6e118
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
)

var clients = make(map[*websocket.Conn]bool) // 接続されるクライアント
var broadcast = make(chan Message)           // メッセージブロードキャストチャネル

// アップグレーダ
var upgrader = websocket.Upgrader{}

// メッセージ用構造体
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

type MessageArray struct {
	Message []Message
}

// メッセージを貯める
var messageDatabase = MessageArray{Message: []Message{}}

func main() {

	log.Println("The Server runs with http://localhost:8990")

	// ファイルサーバーを立ち上げる
	http.HandleFunc("/", htmlHandler)
	http.HandleFunc("/getAll", getAllHandler)
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("templates/resource"))))

	// websockerへのルーティングを紐づけ
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	// サーバーをlocalhostのポート8990で立ち上げる
	err := http.ListenAndServe(":8990", nil)

	// エラーがあった場合ロギングする
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 送られてきたGETリクエストをwebsocketにアップグレード
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 関数が終わった際に必ずwebsocketnのコネクションを閉じる
	defer ws.Close()

	// クライアントを新しく登録
	clients[ws] = true

	for {
		var msg Message
		// 新しいメッセージをJSONとして読み込みMessageオブジェクトにマッピングする
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// データベースに登録
		messageDatabase.Message = append(messageDatabase.Message, msg)

		// 新しく受信されたメッセージをブロードキャストチャネルに送る
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// ブロードキャストチャネルから次のメッセージを受け取る
		msg := <-broadcast
		// 現在接続しているクライアント全てにメッセージを送信する
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
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

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	// 情報取得API
	responseJSON, _ := json.Marshal(messageDatabase)
	fmt.Fprintln(w, string(responseJSON))
}
