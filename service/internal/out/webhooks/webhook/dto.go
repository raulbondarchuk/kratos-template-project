package webhook

type Route1Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Route1Response struct {
	Token string `json:"token"`
}

type Route2Request struct {
	Username string `json:"username"`
}

type Route2Response struct {
	Username string `json:"username"`
}
