package Proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/elazarl/goproxy"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 全局正则表达式切片
var (
	blackListRegexes []*regexp.Regexp
	whiteListRegexes []*regexp.Regexp
	// 黑白名单日志记录器
	bwLogger *log.Logger
)

// loadRegexList 从指定文件加载正则表达式，每行一个（忽略空行），编译后返回正则表达式切片
func loadRegexList(filename string) ([]*regexp.Regexp, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var regexes []*regexp.Regexp
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		re, err := regexp.Compile(line)
		if err != nil {
			log.Printf("编译正则表达式 [%s] 失败: %v", line, err)
			continue
		}
		regexes = append(regexes, re)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return regexes, nil
}

func StartProxy() {
	// --------------------------
	// 初始化黑白名单日志记录器
	// --------------------------
	logFile, err := os.OpenFile("black_white_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
	}
	// 确保程序退出前关闭日志文件
	defer logFile.Close()
	bwLogger = log.New(logFile, "", log.LstdFlags)

	// --------------------------
	// 加载黑名单和白名单正则表达式
	// --------------------------
	blackListRegexes, err = loadRegexList("blacklist.txt")
	if err != nil {
		log.Fatalf("加载黑名单失败: %v", err)
	}
	whiteListRegexes, err = loadRegexList("whitelist.txt")
	if err != nil {
		log.Fatalf("加载白名单失败: %v", err)
	}

	// --------------------------
	// 加载自定义 CA 证书和解密私钥
	// --------------------------
	//证书
	caCertFile := "cert.pem"
	//未加密私钥
	caKeyFile := "key.unencrypted.pem"

	ca, err := tls.LoadX509KeyPair(caCertFile, caKeyFile)
	if err != nil {
		log.Fatalf("加载 CA 证书失败: %v", err)
	}
	// 配置 goproxy 使用自定义 CA
	goproxy.GoproxyCa = ca

	// --------------------------
	// 创建 goproxy 实例
	// --------------------------
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// 拦截 HTTP 请求：读取载荷，检查黑白名单匹配并处理
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		var payload string
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				log.Printf("读取请求 body 失败: %v", err)
			} else {
				payload = string(bodyBytes)
				log.Printf("HTTP 请求 [%s %s] 载荷:\n%s", req.Method, req.URL, payload)
			}
			// 重建 req.Body，确保后续处理正常
			req.Body = io.NopCloser(bytes.NewBuffer([]byte(payload)))
		}

		// 优先检查白名单：匹配则直接放行
		for _, re := range whiteListRegexes {
			if re.MatchString(payload) {
				bwLogger.Printf("白名单匹配（请求）: %s %s, 正则：%s", req.Method, req.URL, re.String())
				return req, nil
			}
		}
		// 检查黑名单：匹配则丢弃报文（返回 403 响应）
		for _, re := range blackListRegexes {
			if re.MatchString(payload) {
				bwLogger.Printf("黑名单匹配（请求）: %s %s, 正则：%s", req.Method, req.URL, re.String())
				resp := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Blocked by proxy filter (blacklist matched)")
				return req, resp
			}
		}
		return req, nil
	})

	// 拦截 HTTP 响应：读取载荷，检查黑白名单匹配并处理
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var payload string
		if resp.Body != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("读取响应 body 失败: %v", err)
			} else {
				payload = string(bodyBytes)
				log.Printf("HTTP 响应 [%s] 载荷:\n%s", resp.Request.URL, payload)
			}
			// 重建 resp.Body，确保客户端正常接收数据
			resp.Body = io.NopCloser(bytes.NewBuffer([]byte(payload)))
		}

		// 优先检查白名单：匹配则直接放行
		for _, re := range whiteListRegexes {
			if re.MatchString(payload) {
				bwLogger.Printf("白名单匹配（响应）: %s, 正则：%s", resp.Request.URL, re.String())
				return resp
			}
		}
		// 检查黑名单：匹配则丢弃报文（返回 403 响应）
		for _, re := range blackListRegexes {
			if re.MatchString(payload) {
				bwLogger.Printf("黑名单匹配（响应）: %s, 正则：%s", resp.Request.URL, re.String())
				newResp := goproxy.NewResponse(resp.Request, goproxy.ContentTypeText, http.StatusForbidden, "Blocked by proxy filter (blacklist matched)")
				return newResp
			}
		}
		return resp
	})

	// 启用 HTTPS MITM 模式（解密 HTTPS 流量）
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	// 启动代理服务器
	log.Println("代理服务器运行在 :8800")
	log.Fatal(http.ListenAndServe(":8800", proxy))
}
