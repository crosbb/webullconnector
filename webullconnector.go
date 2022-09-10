package webullconnector

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LoginReq struct {
	Account     string `json:"account"`
	AccountType string `json:"accountType"`
	DeviceId    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	Grade       int    `json:"grade"`
	Password    string `json:"pwd"`
	RegionId    int    `json:"regionId"`
}

type LoginResp struct {
	ExtInfo ExtraInfo
}

type ExtraInfo struct {
	DesensitizedPhone string
	AccountType       string
	EquipmentCheck    string
}

func login(username string, password string, deviceId string) {

	loginUrl := "https://userapi.webull.com/api/passport/login/v5/account"

	text := "wl_app-a&b@!423^" + password
	hash := md5.Sum([]byte(text))
	pwdHash := hex.EncodeToString(hash[:])

	print(pwdHash)

	reqBody := LoginReq{
		Account:     username,
		AccountType: "2",
		DeviceId:    deviceId,
		DeviceName:  "StockImporter",
		Grade:       1,
		Password:    pwdHash,
		RegionId:    6,
	}

	body, _ := json.Marshal(reqBody)

	request, error := http.NewRequest("POST", loginUrl, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:99.0) Gecko/20100101 Firefox/99.0")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Accept-Language", "en-US,en;q=0.5")
	request.Header.Add("platform", "web")
	request.Header.Add("hl", "en")
	request.Header.Add("os", "web")
	request.Header.Add("osv", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:99.0) Gecko/20100101 Firefox/99.0")
	request.Header.Add("app", "global")
	request.Header.Add("appid", "webull-webapp")
	request.Header.Add("ver", "3.39.18")
	request.Header.Add("lzone", "dc_core_r001")
	request.Header.Add("ph", "MacOS Firefox")
	request.Header.Add("locale", "eng")
	request.Header.Add("device-type", "Web")
	request.Header.Add("did", deviceId)

	client := &http.Client{Timeout: 30 * time.Second}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	respObj := LoginResp{}

	resp := json.NewDecoder(response.Body).Decode(&respObj)

	fmt.Print(resp)

	resbody, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(resbody))

}
