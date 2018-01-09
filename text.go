package aliyun_green

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type AntispamResult struct {
	Scene      string  `json:"scene"`
	Suggestion string  `json:"suggestion"`
	Label      string  `json:"label"`
	Rate       float32 `json:"rate"`
}

type AntispamResponse struct {
	Code int                     `json:"code"`
	Data []*AntispamResponseData `json:data`
}

type AntispamResponseData struct {
	Code    int                   `json:"code"`
	Message string                `json:"msg"`
	DataID  string                `json:"dataId"`
	TaskID  string                `json:"taskId"`
	Content string                `json:"content"`
	Data    *AntispamResponseData `json:"data"`
	Results []*AntispamResult     `json:"results"`
}

func (this *Client) Antispam(text string) (*AntispamResult, error) {
	data := map[string]interface{}{
		"scenes": []string{"antispam"},
		"tasks": []interface{}{
			map[string]string{
				"content": text,
			},
		},
	}

	if len(this.BizType) > 0 {
		data["bizType"] = this.BizType
	}

	response, err := this.Post("/green/text/scan", data)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, &ErrorResponse{StatusCode: response.StatusCode, Body: string(body)}
	}

	responseStruct := AntispamResponse{}
	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return nil, err
	}

	if len(responseStruct.Data) == 0 {
		return nil, errors.New("Blank data with response: " + string(body))
	}
	responseData := responseStruct.Data[0]

	if responseData.Code == 586 {
		return nil, &AlgoFailed{Body: responseData.Message}
	}

	if responseData.Code != 200 {
		return nil, &ErrorResponse{StatusCode: responseData.Code, Body: responseData.Message}
	}

	if len(responseData.Results) == 0 {
		return nil, errors.New("Blank Results with response: " + string(body))
	}

	return responseData.Results[0], nil
}
