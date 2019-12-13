package gateway

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"
)

type Proxy struct {
	Port           int
	TriggerService string
}

func StartProxy(proxy *Proxy) error {
	http.HandleFunc("/", proxy.Proxy)
	http.HandleFunc("/health/liveness", proxy.LivenessProbe)
	http.HandleFunc("/health/readiness", proxy.ReadinessProbe)

	go func() {
		glog.Infof("api server is started")
		if err := http.ListenAndServe(fmt.Sprintf(":%d", proxy.Port), nil); err != nil {
			glog.Errorf("API Listen error:%s ", err.Error())
		}
	}()

	return nil
}

func (p *Proxy) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "liveness")
}

func (p *Proxy) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "readiness")
}

func (p *Proxy) Proxy(w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}

	//log.Printf("%v %v", r.Method, r.RequestURI)
	req, err = http.NewRequest(r.Method, fmt.Sprintf("http://%s:8080%s", p.TriggerService, r.RequestURI), r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}

	// 把 URL Query 参数转换成 Header，给 Tekton Trigger 使用
	queryParams := r.URL.Query()
	if queryParams != nil {
		for key, _ := range queryParams {
			req.Header.Set(key, queryParams.Get(key))
		}
	}

	req.Host = p.TriggerService
	req.Form = r.Form
	req.PostForm = r.PostForm
	resp, err = client.Do(req)
	r.Body.Close()

	// combined for GET/POST
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &HttpConnection{r, resp}

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()

	PrintHTTP(conn)
	//connChannel <- &HttpConnection{r,resp}
}

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

func PrintHTTP(conn *HttpConnection) {
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println("==============================")
	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println(conn.Response.Body)
	fmt.Println("==============================")
}
