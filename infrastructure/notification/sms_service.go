package notification

import (
	"fmt"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/http"
	"github.com/machtwatch/catalystdk/go/log"
)

type ISMSService interface {
	Send(request SMSRequest) error
}

type SMSService struct {
	baseUrl string
}

type SMSRequest struct {
	Category    string `json:"category"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
}

// NewSMSService initialize sms service
func NewSMSService(baseUrl string) *SMSService {
	return &SMSService{
		baseUrl: baseUrl,
	}
}

// Send will call notification service send sms API
func (s SMSService) Send(request SMSRequest) error {
	client := http.DefaultHTTPClient()

	var res response.IntegrationServiceResponse
	_, err := client.R(). // response not used since it's auto unmarshalled
				SetHeader("X-Channel-Id", config.APP_NAME).
				SetBody(request).
				SetResult(&res).
				SetError(&res).
				Post(s.baseUrl + "/message/v1/sms/send")
	if err != nil {
		return err
	}

	if res.Code != response.SUCCESS {
		return fmt.Errorf("error on sending sms: %v", res.Errors)
	}

	log.Infof("sms sucessfully sent to %s", request.PhoneNumber)
	return nil
}
