package integration

type IXMSService interface {
}

type XMSService struct {
	baseUrl string
}

func NewXMSService(baseUrl string) IXMSService {
	return &XMSService{
		baseUrl: baseUrl,
	}
}
