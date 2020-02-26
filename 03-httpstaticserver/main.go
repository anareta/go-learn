package main

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", htmlHandler)

	// サーバーを起動
	// http://localhost:8989
	http.ListenAndServe(":8989", nil)
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {

	mem := sigar.Mem{}
	mem.Get()

	// メモリ使用量の表示
	str := fmt.Sprintf("Mem(GB): total=%.2f", formatGB(mem.Total))
	str = str + fmt.Sprintf(", used=%.2f", formatGB(mem.Used))
	str = str + fmt.Sprintf(", free=%.2f", formatGB(mem.Free))

	// テンプレートをパース
	t := template.Must(template.ParseFiles("templates/template000.html.tpl"))

	// HTMLに入れるオブジェクト
	dat := struct {
		Time string
		Mem  string
	}{
		Mem:  str,
		Time: time.Now().Format("01/02 15:04:05"),
	}

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "template000.html.tpl", dat); err != nil {
		log.Fatal(err)
	}
}

func formatGB(val uint64) float64 {
	return float64(val) / 1024 / 1024 / 1024
}
