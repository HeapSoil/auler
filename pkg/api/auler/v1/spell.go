package v1 


// SpellInfo 指定了咒语的详细信息.
type SpellInfo struct {
	Username  string `json:"username,omitempty"`
	SpellID    string `json:"spellID,omitempty"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// 创建咒语: 指定请求参数和返回参数
// CreateSpellRequest 指定了 `POST /v1/spells` 接口的请求参数.
type CreateSpellRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}
// CreateSpellResponse 指定了 `POST /v1/spells` 接口的返回参数.
type CreateSpellResponse struct {
	SpellID string `json:"spellID"`
}


// 获取咒语：指定返回参数
// GetSpellResponse 指定了 `GET /v1/spells/{spellID}` 接口的返回参数.
type GetSpellResponse SpellInfo

// 罗列用户咒语： 指定请求和返回参数
// ListSpellRequest 指定了 `GET /v1/spells` 接口的请求参数.
type ListSpellRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}
// ListSpellsResponse 指定了 `GET /v1/spells` 接口的返回参数.
type ListSpellResponse struct {
	TotalCount int64       `json:"totalCount"`
	Spells      []*SpellInfo `json:"Spells"`
}

// 更新咒语：指定请求参数
// UpdateSpellRequest 指定了 `PUT /v1/spells` 接口的请求参数.
type UpdateSpellRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}
