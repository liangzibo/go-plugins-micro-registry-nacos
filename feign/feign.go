package feign

import (
	"errors"
	"github.com/ddliu/go-httpclient"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"io/ioutil"
	"net/url"
	"time"
)

var hp *httpclient.HttpClient

func init() {
	NewClient()
}

func NewClient() *httpclient.HttpClient {
	if hp == nil {
		hp = httpclient.NewHttpClient().Defaults(httpclient.Map{
			httpclient.OPT_USERAGENT: "go-feign",
		})
		hp.WithOption(httpclient.OPT_CONNECTTIMEOUT, 30)
		//httpclient.WithOption(httpclient.OPT_CONNECTTIMEOUT_MS,30)
		hp.WithOption(httpclient.OPT_TIMEOUT, 10)
	}
	return hp
}
//转换
func ConvertServiceToUrl(address string) string {
	return "http://" + address
}
//组装
func AssembleUrl(service, path string) string {
	return ConvertServiceToUrl(service) + path
}
//普通 GET
func Get(url string, headers map[string]string) (*httpclient.Response, error) {
	return NewClient().WithHeaders(headers).Get(url)
}
//普通
func Post(url string, param url.Values, headers map[string]string) (*httpclient.Response, error) {
	return NewClient().WithHeaders(headers).Post(url, param)
}
//普通
func PostJson(url string, data interface{}, headers map[string]string) (*httpclient.Response, error) {
	return NewClient().WithHeaders(headers).PostJson(url, data)
}
//普通
func PutJson(url string, data interface{}, headers map[string]string) (*httpclient.Response, error) {
	return NewClient().WithHeaders(headers).PutJson(url, data)
}
//普通
func Delete(url string, headers map[string]string, params ...interface{}) (*httpclient.Response, error) {
	return NewClient().WithHeaders(headers).Delete(url, params)
}
// 获取服务
func GetServiceAddr(r registry.Registry, serviceName string) (string, error) {
	var retryCount int
	var address string
	for {
		servers, err := r.GetService(serviceName)
		if err != nil {
			return "", err
		}
		var services []*registry.Service
		for _, value := range servers {
			//fmt.Println(value.Name, ":", value.Version)
			services = append(services, value)
		}
		next := selector.RoundRobin(services)
		if node, err := next(); err == nil {
			address = node.Address
		}
		if len(address) > 0 {
			//fmt.Println("address:", address)
			return address, nil
		}
		retryCount++
		time.Sleep(time.Second * 1)
		if retryCount >= 5 {
			return "", errors.New("服务不存在")
		}
	}
}
// 获取服务 GET
func GetFeign(r registry.Registry, serviceName, path string, headers map[string]string) (string, error) {
	addr, err := GetServiceAddr(r, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) <= 0 {
		return "", errors.New("服务不存在")
	}
	url := AssembleUrl(addr, path)
	rsp, err := Get(url, headers)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.New("转换错误")
	}
	return string(bytes), nil
}
// 获取服务 POST
func PostFeign(r registry.Registry, serviceName, path string, param url.Values, headers map[string]string) (string, error) {
	addr, err := GetServiceAddr(r, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) <= 0 {
		return "", errors.New("服务不存在")
	}
	url := AssembleUrl(serviceName, path)
	rsp, err := Post(url, param, headers)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.New("转换错误")
	}
	return string(bytes), nil
}
// 获取服务 POSTJSON
func PostJsonFeign(r registry.Registry, serviceName, path string, data interface{}, headers map[string]string) (string, error) {
	addr, err := GetServiceAddr(r, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) <= 0 {
		return "", errors.New("服务不存在")
	}
	url := AssembleUrl(addr, path)
	rsp, err := PostJson(url, data, headers)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.New("转换错误")
	}
	return string(bytes), nil
}
// 获取服务 PUTJSON
func PutJsonFeign(r registry.Registry, serviceName, path string, data interface{}, headers map[string]string) (string, error) {
	addr, err := GetServiceAddr(r, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) <= 0 {
		return "", errors.New("服务不存在")
	}
	url := AssembleUrl(addr, path)
	rsp, err := PutJson(url, data, headers)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.New("转换错误")
	}
	return string(bytes), nil
}
// 获取服务 DELETE
func DeleteFeign(r registry.Registry, serviceName, path string, headers map[string]string, params ...interface{}) (string, error) {
	addr, err := GetServiceAddr(r, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) <= 0 {
		return "", errors.New("服务不存在")
	}
	url := AssembleUrl(addr, path)
	rsp, err := Delete(url, headers, params)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.New("转换错误")
	}
	return string(bytes), nil
}
