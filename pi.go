package main

import (
    "io"
    "net/http"
    "log"
    "os"
    "fmt"
    "time"
    "github.com/bmizerany/pat"
    "github.com/hoisie/mustache"
)

// 默认的http端口
var listenPort int = 80
var (
    assetsDir           string
    controllerName      string
)

//TODO: 实现一个Controller的Interface，默认从对应的views目录加载index.html.mustache并Write

func Render(t string, context ...interface{}) string {
    return mustache.RenderFileInLayout("views/" + controllerName + "/" + t + ".html.mustache", 
        "views/layout/pi.html.mustache", context)
}

// hello world, the web server
func HelloController(w http.ResponseWriter, r *http.Request) {
    controllerName = "hello"
    var body string = Render("index", nil)
    io.WriteString(w, body)
}

// Serve static files
func AssetsServer(w http.ResponseWriter, r *http.Request) {
    filename := string(assetsDir + r.URL.Path[8:])
    // fmt.Println(filename)

    // Set MIME
    if filename[len(filename)-3:] == ".css" {
        w.Header().Set("Content-Type", "text/css")
    }

    // Set expire headers to now + 1 year
    yearLater := time.Now().AddDate(1, 0, 0)
    w.Header().Set("Expires", yearLater.Format(http.TimeFormat))

    http.ServeFile(w, r, filename)
}

// About
func AboutController(w http.ResponseWriter, r *http.Request) {
    controllerName = "about"
    io.WriteString(w, Render("index", nil))
}

// Contact
func ContactController(w http.ResponseWriter, r *http.Request) {
    controllerName = "contact"
    io.WriteString(w, Render("index", nil))
}

func main() {
    // 静态资源目录
    assetsDir = "./public/"

    m := pat.New()
    m.Get("/", http.HandlerFunc(HelloController))
    m.Get("/assets/:file", http.HandlerFunc(AssetsServer))
    m.Get("/about", http.HandlerFunc(AboutController))
    m.Get("/contact", http.HandlerFunc(ContactController))

    // Register this pat with the default serve mux so that other packages
    // may also be exported. (i.e. /debug/pprof/*)
    http.Handle("/", m)

    // 非生产环境http使用7000端口
    if os.Getenv("PI_ENV") != "production" {
        listenPort = 7000
    }
    var fullListenParam string = fmt.Sprintf(":%d", listenPort)
    err := http.ListenAndServe(fullListenParam, nil)

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}