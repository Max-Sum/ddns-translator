package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func HandleDnspod(w http.ResponseWriter, req *http.Request) {
	// Prepare
	query := req.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	domain := query.Get("domain_id")
	record := query.Get("record_id")
	sub_domain := query.Get("sub_domain")
	if id == "" || token == "" || domain == "" || record == "" || sub_domain == "" {
		w.Write([]byte("id: " + id + "\n"))
		w.Write([]byte("token: " + token + "\n"))
		w.Write([]byte("domain_id: " + domain + "\n"))
		w.Write([]byte("record_id: " + record + "\n"))
		w.Write([]byte("sub_domain: " + sub_domain + "\n"))
		w.Write([]byte("Error: missing id / token / domain_id / record_id / sub_domain arg.\n"))
		return
	}
	record_line := query.Get("record_line")
	if record_line == "" {
		record_line = "默认"
	}
	ip := query.Get("ip")
	if ip == "" {
		ip = ClientIP(req)
	}
	w.Write([]byte(ip))
	// Set Params
	v := url.Values{}
    v.Set("login_token", id + "," + token)
	v.Set("format", "json")
	v.Set("domain_id", domain)
	v.Set("record_id", record)
	v.Set("record_line", record_line)
	v.Set("sub_domain", sub_domain)
	v.Set("value", ip)
	// Post
    body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
    client := &http.Client{}//客户端,被Get,Head以及Post使用
    request, err := http.NewRequest("POST", "https://dnsapi.cn/Record.Ddns", body)
    if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "text/json")
	resp, err := client.Do(request)
    if err != nil {
		w.Write([]byte(err.Error()))
		return
    }
    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)
	w.Write(respBody)
}

func main() {
	http.HandleFunc("/dnspod", HandleDnspod)
    http.ListenAndServe(":80", nil)
}