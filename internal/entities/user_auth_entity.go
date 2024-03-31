package entities

type UserCacheInfo struct {
	UserID    string `json:"userId"`
	EmailHash string `json:"emailHash"`
}
