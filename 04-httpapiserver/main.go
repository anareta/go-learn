package main

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/mitchellh/go-ps"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", htmlHandler)
	http.HandleFunc("/api", apiHandler)

	// サーバーを起動
	// http://localhost:8989
	http.ListenAndServe(":8989", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	processes, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	// 重複削除
	m := make(map[string]bool)
	uniq := []ps.Process{}

	for _, ele := range processes {
		if !m[ele.Executable()] {
			m[ele.Executable()] = true
			uniq = append(uniq, ele)
		}
	}

	// メモリ使用量（全体）
	mem := sigar.Mem{}
	mem.Get()

	memAll := fmt.Sprintf("total=%.2f [GB],", formatGB(mem.Total))
	memAll = memAll + fmt.Sprintf("   used=%.2f  [GB],", formatGB(mem.Used))
	memAll = memAll + fmt.Sprintf("   free=%.2f  [GB]", formatGB(mem.Free))

	message := memAll + "<br><br>"
	for _, p := range uniq {
		if p.Pid() == 0 {
			continue
		}
		mem := sigar.ProcMem{}
		if err := mem.Get(p.Pid()); err != nil {
			continue
		}

		message += p.Executable() + " : " + fmt.Sprintf("%.2f", float64(mem.Size)/1024/1024) + "[MB]<br>"
	}

	fmt.Fprintln(w, message)

}

func htmlHandler(w http.ResponseWriter, r *http.Request) {

	// ホスト名の取得
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// テンプレートをパース
	t := template.Must(template.ParseFiles("templates/template.html.tpl"))

	// HTMLに入れるオブジェクト
	dat := struct {
		HostName string
	}{
		HostName: name,
	}

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "template.html.tpl", dat); err != nil {
		log.Fatal(err)
	}
}

func formatGB(val uint64) float64 {
	return float64(val) / 1024 / 1024 / 1024
}
