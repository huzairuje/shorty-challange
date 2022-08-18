package tiny_url

import (
	"errors"
	"time"

	"test_amartha_muhammad_huzair/pkg/utils"

	"github.com/lucasjones/reggen"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

var ListAllTinyUrl []*Data

func (s Service) GetAllData() []*Data {
	return ListAllTinyUrl
}

func (s Service) GetSingleData(shortCode string) (*Data, error) {
	var data *Data
	if len(ListAllTinyUrl) > 0 {
		for _, v := range ListAllTinyUrl {
			if v.ShortCode == shortCode {
				data = &Data{
					ShortCode:     v.ShortCode,
					Url:           v.Url,
					StartDate:     v.StartDate,
					LastSeenDate:  v.LastSeenDate,
					RedirectCount: v.RedirectCount,
				}
				return data, nil
			}
		}

		if data == nil {
			newErr := errors.New(utils.ShortCodeIsNotExist)
			return nil, newErr
		}
	}
	return nil, nil
}

func (s Service) UpdateStat(shortCode string) {
	timeNow := time.Now().Format(time.RFC3339)
	for index, value := range ListAllTinyUrl {
		if shortCode == value.ShortCode {
			ListAllTinyUrl[index].RedirectCount = ListAllTinyUrl[index].RedirectCount + 1
			ListAllTinyUrl[index].LastSeenDate = timeNow
		}
	}
}

func (s Service) CreateData(dataReq Data) (string, error) {
	var err error

	if dataReq.ShortCode == "" {
		dataReq.ShortCode, err = s.generateShortCode()
		if err != nil {
			return "", err
		}

		if !utils.IsValidShortCode(dataReq.ShortCode) {
			newErr := errors.New("shortcode generated is not valid, you can try again")
			return "", newErr
		}
	}

	timeNow := time.Now().Format(time.RFC3339)
	newData := &Data{
		ShortCode:     dataReq.ShortCode,
		Url:           dataReq.Url,
		StartDate:     timeNow,
		RedirectCount: 0,
	}

	ListAllTinyUrl = append(ListAllTinyUrl, newData)

	return newData.ShortCode, nil
}

func (s Service) generateShortCode() (string, error) {
	strFinal, err := reggen.Generate("^[0-9a-zA-Z_]{6}$", 1)
	if err != nil {
		return "", err
	}
	return strFinal, nil
}
