package models

type (
	Company struct {
		ID          string `json:"id"`
		Name        string `json:"name" valid:"required,ascii"`
		Year        uint32 `json:"founded_at" valid:"optional,int"`
		Description string `json:"description" valid:"optional,ascii"`
		Country     string `json:"country" valid:"optional,ascii"`
		// Currency used in company filings.
		Currency string `json:"currency,omitempty"`
		// Listed exchange.
		Exchange string `json:"exchange,omitempty"`
		// Company name.
		// Company symbol/ticker as used on the listed exchange.
		Ticker string `json:"ticker,omitempty"`
		// IPO date.
		Ipo string `json:"ipo,omitempty"`
		// Market Capitalization.
		// Logo image.
		Logo string `json:"logo,omitempty"`
		// Company website.
		Weburl string `json:"weburl,omitempty"`
		// Finnhub industry classification.
		Industry string `json:"finnhubIndustry,omitempty"`
	}
	LikeUnlikeCompany struct {
		Name string `json:"name" valid:"required,ascii"`
	}
)
