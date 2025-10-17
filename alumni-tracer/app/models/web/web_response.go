package web

type PaginationResponse struct{
	Data interface{} `json:"data"`
	Meta MetaInfo
}