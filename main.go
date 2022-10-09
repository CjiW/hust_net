package main

import (
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"os"
)

type Paras struct {
	Method          string `url:"method"`
	Userid          string `url:"userId" json:"userid"`
	Passwd          string `url:"password" json:"passwd"`
	QueryStr        string `url:"queryString"`
	PasswordEncrypt string `url:"passwordEncrypt"`
	Service         string `url:"service"`
	OperatorPwd     string `url:"operatorPwd" `
	Validcode       string `url:"validcode"`
}

func ReadConf(filename string) *Paras {
	paras := Paras{}
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = json.Unmarshal(data, &paras)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	paras.Method = "login"
	paras.PasswordEncrypt = "false"
	return &paras
}

// GetQueryStr 通过请求 123.123.123.123 获取设备信息
func GetQueryStr() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://123.123.123.123:80", nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	tmp1 := string(body)
	cond := regexp2.MustCompile("(?<=(\\?))(.*)(?=')", 0)
	mat, err := cond.FindStringMatch(tmp1)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return mat.String()
}

// Login 登录请求，参数有 账户信息，设备信息等
func Login(paras *Paras) {
	para, _ := query.Values(paras)
	u := "http://172.18.18.60:8080/eportal/InterFace.do?" + para.Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
func main() {
	paras := ReadConf("./login.json")
	if paras == nil {
		return
	}
	paras.QueryStr = GetQueryStr()
	if len(paras.QueryStr) == 0 {
		return
	}
	Login(paras)
}
