package webullconnector

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

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

type MfaReq struct {
	Account     string `json:"account"`
	AccountType string `json:"accountType"`
	CodeType    int    `json:"codeType"`
}

type LoginResp struct {
	ExtInfo ExtraInfo
}

type ExtraInfo struct {
	DesensitizedPhone string
	AccountType       string
	EquipmentCheck    string
}

func buildHeaders(request *http.Request, deviceId string) {

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

}

func Login(username string, password string, deviceId string) (LoginResp, error) {

	loginUrl := "https://userapi.webull.com/api/passport/login/v5/account"

	text := "wl_app-a&b@!423^" + password
	hash := md5.Sum([]byte(text))
	pwdHash := hex.EncodeToString(hash[:])

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

	request, newReqErr := http.NewRequest("POST", loginUrl, bytes.NewBuffer(body))
	if newReqErr != nil {
		return LoginResp{}, newReqErr
	}

	buildHeaders(request, deviceId)

	client := &http.Client{Timeout: 30 * time.Second}
	response, reqErr := client.Do(request)
	if reqErr != nil {
		return LoginResp{}, reqErr
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return LoginResp{}, errors.New(response.Status)
	}

	// fmt.Println("response Status:", response.Status)
	// fmt.Println("response Headers:", response.Header)

	respObj := LoginResp{}

	respErr := json.NewDecoder(response.Body).Decode(&respObj)

	if respErr != nil {
		return LoginResp{}, respErr
	}

	// fmt.Print(resp)

	// resbody, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("response Body:", string(resbody))

	return respObj, nil
}

func GetMfa(username string, deviceId string) error {
	mfaUrl := "https://userapi.webullfintech.com/api/passport/v2/verificationCode/send"

	reqBody := MfaReq{
		Account:     username,
		AccountType: "2",
		CodeType:    5,
	}
	body, _ := json.Marshal(reqBody)

	request, newReqErr := http.NewRequest("POST", mfaUrl, bytes.NewBuffer(body))
	if newReqErr != nil {
		return errors.New("Error getting Mfa" + newReqErr.Error())
	}

	buildHeaders(request, deviceId)

	client := &http.Client{Timeout: 30 * time.Second}
	response, reqErr := client.Do(request)
	if reqErr != nil {
		return errors.New("Error getting Mfa" + reqErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New("Error getting Mfa" + response.Status)
	}

	return nil
}
