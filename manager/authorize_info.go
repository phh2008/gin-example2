package manager

type AuthorizeInfo struct {
	Uid         string   `json:"uid"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}
