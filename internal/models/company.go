package models

type (
	Company struct {
		ID          string `json:"id"`
		Name        string `json:"name" valid:"required,ascii"`
		Year        uint32 `json:"founded_at" valid:"optional,int"`
		Description string `json:"description" valid:"optional,ascii"`
	}
	LikeUnlikeCompany struct {
		Name string `json:"name" valid:"required,ascii"`
	}
)
