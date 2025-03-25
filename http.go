package main

import (
	"fmt"
	"net/http"
	"net/url"
)

// 处理请求，解析并打印参数（保留 text 的编码值）
func sendenfy2Handler(w http.ResponseWriter, r *http.Request) {
	// 仅允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取原始未解码的查询字符串
	rawQuery := r.URL.RawQuery

	// 手动解析查询参数（避免自动解码）
	params, err := url.ParseQuery(rawQuery)
	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// 提取参数（直接获取原始值）
	phone := params.Get("phone")
	text := params.Get("text") // 保留编码后的值（如 "%E4%BD%A0%E5%A5%BD"）
	sndchl := params.Get("sndchl")
	source := params.Get("source")

	// 打印参数（控制台输出）
	fmt.Printf("Received request:\n")
	fmt.Printf("phone : %s\n", phone)
	fmt.Printf("text  : %s\n", text) // 直接输出编码后的字符串
	fmt.Printf("sndchl: %s\n", sndchl)
	fmt.Printf("source: %s\n\n", source)

	// 返回响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Parameters received"))
}

func main() {
	// 注册路由
	http.HandleFunc("/sendenfy2", sendenfy2Handler)

	// 启动服务器
	fmt.Println("Server listening on :9999")
	http.ListenAndServe(":9999", nil)
}
