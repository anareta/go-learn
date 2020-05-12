// longpollを使ってServer->Client通信やる
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jcuga/golongpoll"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type messages struct {
	Message []message
}

type message struct {
	Body string
	Time string
}

// ロングポーリング
// This launches a goroutine and creates channels for all the plumbing
var manager, _ = golongpoll.StartLongpoll(golongpoll.Options{})

var messageDatabase = messages{Message: []message{}}

func main() {
	fmt.Println("The Server runs with http://localhost:8989")

	// Expose events to browsers
	// See subsection on how to interact with the subscription handler
	http.HandleFunc("/", htmlHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/postMessage", postHandler)
	http.HandleFunc("/events", manager.SubscriptionHandler)

	// https://tech-up.hatenablog.com/entry/2018/12/28/120517
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("templates/resource"))))

	// サーバーを起動
	// http://localhost:8989
	http.ListenAndServe(":8989", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// 情報取得API
	// encode json
	responseJSON, _ := json.Marshal(messageDatabase)
	fmt.Fprintln(w, string(responseJSON))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// 投稿API
	var _ = r.ParseForm()
	var update bool = false

	for k := range r.Form {
		messageDatabase.Message = append(messageDatabase.Message, createMessage(k))
		update = true

		if len(messageDatabase.Message) > 50 {
			// メッセージ件数が上限を超える場合は古いデータから消す
			unset := func(s []message, i int) []message {
				if i >= len(s) {
					return s
				}
				return append(s[:i], s[i+1:]...)
			}

			messageDatabase.Message = unset(messageDatabase.Message, 0)
		}
	}

	if update {
		manager.Publish("update", "message updated")
	}
}

func createMessage(m string) message {
	t := time.Now()
	const timelayout = "15:04:05"
	msg := message{Body: m, Time: t.Format(timelayout)}
	return msg
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
