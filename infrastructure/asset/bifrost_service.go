package asset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
)

type IBifrostService interface {
	Upload(src multipart.File, header *multipart.FileHeader) (BifrostResponse, error)
}

type BifrostService struct {
	baseUrl string
}

type BifrostResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    BifrostData `json:"data"`
	Errors  []string    `json:"errors"`
}

type BifrostData struct {
	FileType      string `json:"file_type"`
	FileName      string `json:"file_name"`
	FileSize      int    `json:"file_size"`
	FileExtension string `json:"extension"`
	Url           string `json:"url"`
}

func NewBifrostService(baseUrl string) *BifrostService {
	return &BifrostService{
		baseUrl: baseUrl,
	}
}

func (b BifrostService) Upload(src multipart.File, header *multipart.FileHeader) (data BifrostResponse, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	dst, err := writer.CreatePart(header.Header)
	if err != nil {
		return data, err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return data, err
	}

	err = writer.Close()
	if err != nil {
		return data, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/asset/v1/upload", b.baseUrl), payload)
	if err != nil {
		return data, err
	}

	req.Header.Add("X-Tenant-Code", config.BIFROST_TENANT_CODE)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	if data.Code != "SUCCESS" {
		return data, fmt.Errorf(data.Message)
	}

	return data, nil
}
