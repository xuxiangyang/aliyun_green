package aliyun_green

type Client struct {
	AccessKeyID     string
	AccessKeySecret string
	BizType         string
}

const (
	PASS   = "pass"
	REVIEW = "review"
	BLOCK  = "block"
)

func NewClient(accessKeyID string, accessKeySecret string) *Client {
	return &Client{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	}
}

func NewClientWithBizType(accessKeyID string, accessKeySecret string, bizType string) *Client {
	return &Client{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		BizType:         bizType,
	}
}
