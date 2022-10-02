package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

var addr = flag.String("addr", ":9999", "http service address")

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	log.SetFlags(0)
	flag.Parse()
	http.Handle("/", http.HandlerFunc(QR))
	url := fmt.Sprintf("http://localhost%v", *addr)
	log.Printf("url = %#v", url)
	openBrowser(url)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func QR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Please open %v in your browser", *addr)
	}
	if err != nil {
		log.Fatal(err)
	}
}

const templateStr = `
<html>

<head>
    <title>QR Link Generator</title>
</head>

<body>
    {{if .}}
    <img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
    <br />
    {{.}}
    <br />
    <br />
    {{end}}
    <form action="/" name="f" method="GET">
        <input maxlength="1024" size="70" name="s" value="" title="Text to QR Encode" />
        <input type="submit" value="Show QR" name="qr" />
    </form>
</body>

</html>
`
