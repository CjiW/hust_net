package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"math/big"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"strings"
)

type Paras struct {
	Userid string `json:"userid"`
	Passwd string `json:"passwd"`
}

// ReadConf 读取配置
func ReadConf(filename string) *Paras {
	paras := Paras{}
	data, _ := os.ReadFile(filename)
	json.Unmarshal(data, &paras)
	return &paras
}

// GetQueryStr 从 123.123.123.123 获取状态信息
func GetQueryStr() string {
	url := "http://123.123.123.123"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	body := string(data)
	ret := strings.Split(strings.Split(body, "'")[1], "?")[1]
	return ret
}

type PubKey struct {
	PublicKeyExponent string `json:"publicKeyExponent"`
	PublicKeyModulus  string `json:"publicKeyModulus"`
}

// GetPubKey 通过请求 192.168.50.3:8080 获取公钥信息
func GetPubKey(queryStr string) PubKey {
	url := "http://192.168.50.3:8080/eportal/InterFace.do?method=pageInfo"
	client := &http.Client{}
	ule := "queryString=" + url2.QueryEscape(queryStr)
	bd := strings.NewReader(ule)
	req, _ := http.NewRequest("POST", url, bd)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, _ := client.Do(req)
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	var pk PubKey
	json.Unmarshal(data, &pk)
	return pk
}

type LoginPara struct {
	Method          string `url:"method"`
	Userid          string `url:"userId"`
	Passwd          string `url:"password"`
	QueryStr        string `url:"queryString"`
	PasswordEncrypt string `url:"passwordEncrypt"`
	Service         string `url:"service"`
	OperatorPwd     string `url:"operatorPwd" `
	Validcode       string `url:"validcode"`
}

// Login 登录请求，参数有 账户信息，设备信息等
func Login(paras *LoginPara) {
	para, _ := query.Values(paras)
	u := "http://192.168.50.3:8080/eportal/InterFace.do?" + para.Encode()
	client := &http.Client{}
	req, _ := http.NewRequest("POST", u, nil)

	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}

func main() {
	paras := ReadConf("./login.json")
	queryStr := GetQueryStr()
	pk := GetPubKey(queryStr)
	E, _ := strconv.ParseInt(pk.PublicKeyExponent, 16, 64)

	parseQuery, _ := url2.ParseQuery(queryStr)
	msg := paras.Passwd + ">" + parseQuery["mac"][0]
	encryptData := EncryptData(pk.PublicKeyModulus, E, msg)
	loginPara := LoginPara{
		Method:          "login",
		Userid:          paras.Userid,
		Passwd:          encryptData,
		QueryStr:        queryStr,
		PasswordEncrypt: "true",
		Service:         "",
		OperatorPwd:     "",
		Validcode:       "",
	}
	Login(&loginPara)
}

// EncryptData RSA 公钥加密 res = msg ^ e mod m
func EncryptData(modulus string, e int64, msg string) string {
	N := new(big.Int)
	N.SetString(modulus, 16)
	E := big.NewInt(e)
	msgNum := big.NewInt(0)
	msgLen := len(msg)
	for i := 0; i < msgLen; i++ {
		ch := big.NewInt(int64(msg[msgLen-i-1]))
		ch.Lsh(ch, uint(8*i))
		msgNum.Add(msgNum, ch)
	}
	return fmt.Sprintf("%x", msgNum.Exp(msgNum, E, N))
}
