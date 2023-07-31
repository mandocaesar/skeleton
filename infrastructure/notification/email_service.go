package notification

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/http"
)

type INotificationEmail interface {
	Send(req EmailRequest, tmpl EmailTemplate) error
}

type EmailService struct {
	embed.FS
	baseUrl    string
	prefixPath string
}

func NewEmailService(baseUrl string, prefixPath string) *EmailService {
	return &EmailService{
		baseUrl:    baseUrl,
		prefixPath: prefixPath,
	}
}

type EmailTemplate struct {
	File     string
	Variable interface{}
}

type EmailRequest struct {
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
	From       *EmailRecipient   `json:"from"`
	To         []EmailRecipient  `json:"to"`
	ReplyTo    *[]EmailRecipient `json:"reply_to"`
	Cc         *[]EmailRecipient `json:"cc"`
	Attachment *[]string         `json:"attachment"`
}

type EmailRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (e EmailService) Send(req EmailRequest, tmpl EmailTemplate) error {
	body, err := tmpl.Parse(e.prefixPath)
	if err != nil {
		return err
	}

	req.Body = body.String()
	req.From = &EmailRecipient{
		Name:  config.NOTIFICATION_SERVICE_SENDER_EMAIL,
		Email: config.NOTIFICATION_SERVICE_SENDER_EMAIL,
	}

	var res response.IntegrationServiceResponse

	client := http.DefaultHTTPClient()
	_, err = client.R().
		SetHeader("X-Channel-Id", config.APP_NAME).
		SetBody(req).
		SetResult(&res).
		SetError(&res).
		Post(e.baseUrl + "/message/v1/email/send")
	if err != nil {
		return err
	}

	if res.Code != response.SUCCESS {
		return fmt.Errorf("error on sending email: %v", res.Errors)
	}

	return nil
}

func (t EmailTemplate) Parse(prefixPath string) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/files/template/email/%s", prefixPath, t.File))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	err = tmpl.ExecuteTemplate(buf, t.File, t.Variable)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
