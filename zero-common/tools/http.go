package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Post(url string, params any) ([]byte, error) {
	marshal, _ := json.Marshal(params)
	s := string(marshal)
	reqBody := strings.NewReader(s)
	httpReq, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()
	rspBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}

func PostWithHeader(path string, params any, m map[string]string, proxy string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	marshal, _ := json.Marshal(params)
	s := string(marshal)
	reqBody := strings.NewReader(s)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}
	if m != nil {
		for k, v := range m {
			httpReq.Header.Set(k, v)
		}
	}
	//httpReq.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	if proxy != "" {
		proxyAddress, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddress),
			},
		}
	}
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()
	rspBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}

func GetWithHeader(path string, m map[string]string, proxy string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if m != nil {
		for k, v := range m {
			httpReq.Header.Set(k, v)
		}
	}
	httpReq.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	if proxy != "" {
		proxyAddress, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddress),
			},
		}
	}
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()
	rspBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}
