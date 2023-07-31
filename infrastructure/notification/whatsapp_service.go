package notification

import (
	"fmt"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/http"
	"github.com/machtwatch/catalystdk/go/log"
)

type WhatsAppService struct {
	baseUrl string
}

type WhatsAppOTPRequest struct {
	PhoneCode   string `json:"phone_code"`
	PhoneNumber string `json:"phone_number"`
	OTPCode     string `json:"otp_code"`
}

// NewWhatsAppService initialize whatsapp service
func NewWhatsAppService(baseUrl string) *WhatsAppService {
	return &WhatsAppService{
		baseUrl: baseUrl,
	}
}

// SendOTP will call notification service send whatsapp OTP API
func (w WhatsAppService) SendOTP(request WhatsAppOTPRequest) error {
	client := http.DefaultHTTPClient()

	var res response.IntegrationServiceResponse
	_, err := client.R(). // response not used since it's auto unmarshalled
				SetHeader("X-Channel-Id", config.APP_NAME).
				SetBody(request).
				SetResult(&res).
				SetError(&res).
				Post(w.baseUrl + "/message/v1/whatsapp/otp")
	if err != nil {
		return err
	}

	if res.Code != response.SUCCESS {
		return fmt.Errorf("error on sending whatsapp otp: %v", res.Errors)
	}

	log.Infof("whatsapp otp sent to %v", request.PhoneNumber)
	return nil
}
