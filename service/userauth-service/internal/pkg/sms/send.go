package sms

import (
	"fmt"
	"os"
	"userauth-service/internal/conf"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SendSms 发送短信验证码
func SendSms(smsConf *conf.Sms, mobile, code string) error {
	// 从环境变量读取，优先于配置
	apiKey := smsConf.ApiKey
	apiSecret := smsConf.ApiSecret
	
	if v := os.Getenv("SMS_API_KEY"); v != "" {
		apiKey = v
	}
	if v := os.Getenv("SMS_API_SECRET"); v != "" {
		apiSecret = v
	}

	if apiKey == "" || apiSecret == "" {
		return fmt.Errorf("SMS API credentials not configured")
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(apiKey),
		AccessKeySecret: tea.String(apiSecret),
		RegionId:        tea.String(smsConf.RegionId),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	
	client, err := dysmsapi.NewClient(config)
	if err != nil {
		return err
	}
	
	request := &dysmsapi.SendSmsRequest{}
	request.SetTemplateCode(smsConf.TemplateCode)
	request.SetTemplateParam(fmt.Sprintf("{\"code\":\"%s\"}", code))
	request.SetPhoneNumbers(mobile)
	request.SetSignName(smsConf.SignName)
	
	_, err = client.SendSms(request)
	return err
}
