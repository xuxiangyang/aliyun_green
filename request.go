package aliyun_green

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

const (
	HOST = "http://green.cn-shanghai.aliyuncs.com"
)

type ErrorResponse struct {
	StatusCode int
	Body       string
}

func (this *Client) Post(path string, data map[string]interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", HOST+path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	req.Header.Set("x-acs-version", "2017-01-12")
	req.Header.Set("x-acs-signature-nonce", uuid.NewV4().String())
	req.Header.Set("x-acs-signature-version", "1.0")
	req.Header.Set("x-acs-signature-method", "HMAC-SHA1")
	req.Header.Set("Content-MD5", this.computeBodyMD5(body))
	req.Header.Set("Authorization", "acs "+this.AccessKeyID+":"+this.computeSignature("POST", path, body, req.Header))
	return (&http.Client{}).Do(req)
}

func (this *Client) computeSignature(method string, path string, data []byte, header http.Header) string {
	str := method + "\n" + "application/json" + "\n" + header.Get("Content-MD5") + "\n" + "application/json" + "\n" + header.Get("Date") + "\n" + this.serializeHeader(header) + "\n" + path
	mac := hmac.New(sha1.New, []byte(this.AccessKeySecret))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))

}

func (this *Client) computeBodyMD5(body []byte) string {
	var md5Ctx = md5.New()
	md5Ctx.Write(body)
	return base64.StdEncoding.EncodeToString(md5Ctx.Sum(nil))
}

func (this *Client) serializeHeader(header http.Header) string {
	return "x-acs-signature-method:" + header.Get("x-acs-signature-method") + "\n" + "x-acs-signature-nonce:" + header.Get("x-acs-signature-nonce") + "\n" + "x-acs-signature-version:" + header.Get("x-acs-signature-version") + "\n" + "x-acs-version:" + header.Get("x-acs-version")
}

func (this *ErrorResponse) Error() string {
	return fmt.Sprintf("AliyunGreen response with %d status, body: %s", this.StatusCode, this.Body)
}
