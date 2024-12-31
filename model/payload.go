package model

type TokenDetail struct {
	Email        string
	ProfileId    uint
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}
type TokenLoadResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	Profile      interface{} `json:"profile"`
	LoginType    string      `json:"loginType"`
}
