package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type (
	Company struct {
		ID          string          `json:"id"`
		Name        string          `json:"name" valid:"required,ascii"`
		IPO         time.Time       `json:"ipo" valid:"optional,int" db:"ipo"`
		Description string          `json:"description" valid:"optional,ascii"`
		Country     string          `json:"country" valid:"optional,ascii"`
		Ticker      string          `json:"ticker,omitempty"`
		Logo        string          `json:"logo,omitempty"`
		Weburl      string          `json:"weburl,omitempty"`
		Attributes  AttributesCmpny `json:"attributes,omitempty"`
	}
	LikeUnlikeCompany struct {
		Ticker string `json:"ticker" valid:"required,ascii"`
	}
	AttributesCmpny struct {
		Currency string `json:"currency,omitempty"`
		Exchange string `json:"exchange,omitempty"`
		Industry string `json:"finnhubIndustry,omitempty"`
	}
	CompanyByTicker struct {
		SrackingId string `json:"trackingId"`
		Status     string `json:"status"`
		Payload    struct {
			Total       int `json:"total"`
			Instruments []struct {
				Figi     string `json:"figi"`
				Ticker   string `json:"ticker"`
				Isin     string `json:"isin"`
				Currency string `json:"currency"`
				Name     string `json:"name"`
				Type     string `json:"type"`
			} `json:"instruments"`
		} `json:"payload"`
	}
)

//для json scan
func (pc *AttributesCmpny) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		_ = json.Unmarshal(v, &pc)
		return nil
	case string:
		_ = json.Unmarshal([]byte(v), &pc)
		return nil
	case nil:
		pc = &AttributesCmpny{}
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (pc *AttributesCmpny) Value() (driver.Value, error) {
	return json.Marshal(pc)
}
