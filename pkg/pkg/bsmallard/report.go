package bsmallard

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	MallardClient struct {
		Option *Option
	}
)

var (
	MallardIns *MallardClient
)

// NewClient connects clients via option
// if one client connects error will throw panic
func MustNewClient(mallardOption *Option) (*MallardClient, error) {
	if !mallardOption.Switch {
		return nil, errors.New("mallard switch is close")
	}
	if GetMallardIns() == nil {
		MallardIns = &MallardClient{Option: mallardOption}
		if MallardIns.Option.ReportHttpTransport == nil {
			MallardIns.Option.ReportHttpTransport = &http.Transport{
				ResponseHeaderTimeout: time.Minute * 5,
				IdleConnTimeout:       time.Minute * 5,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:          50,
			}
		}
		MallardIns.Option.ReportHttpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return MallardIns, nil
}

func GetMallardIns() *MallardClient {
	return MallardIns
}

func (m MallardClient) Report(mallardDataSlice []*MetricRaw) error {
	mallardDataByte, mallardErr := m.checkReport(mallardDataSlice)
	if mallardErr != nil {
		return errors.New("Failed to Marshal mallardDataSlice. -" + mallardErr.Error())
	}

	client := &http.Client{Transport: m.Option.ReportHttpTransport}
	res, resErr := client.Post(m.Option.ReportUrl, "application/json", bytes.NewBuffer(mallardDataByte))
	if resErr != nil {
		return errors.New("Failed to post mallardDataByte. -" + resErr.Error())
	}

	defer res.Body.Close()
	body, bodyError := ioutil.ReadAll(res.Body)
	if bodyError != nil {
		return errors.New("mallard api request error : " + bodyError.Error())

	}
	if res.StatusCode >= 300 {
		return errors.New("mallard api response failed url " + m.Option.ReportUrl + "status" + strconv.Itoa(res.StatusCode) + "body info: " + string(body))
	}
	return nil
}

func (m MallardClient) checkReport(mallardDataSlice []*MetricRaw) ([]byte, error) {
	for _, val := range mallardDataSlice {
		if val.Name == "" {
			return nil, errors.New("the metric name is not allowed to be empty")
		}
		if strings.Contains(m.Option.ReportUrl, "open/metric") && val.Endpoint == "" {
			return nil, errors.New("the endpoint is not allowed to be empty. - metric name: " + val.Name)
		}
	}
	mallardDataByte, mallardErr := json.Marshal(mallardDataSlice)
	return mallardDataByte, mallardErr

}
