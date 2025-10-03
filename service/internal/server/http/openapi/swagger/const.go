package swagger

const (
	loginURL   = "http://10.70.20.90:8134/v1/login" // TODO: change to the correct login URL
	cookieName = "swagger_session"
	cookiePath = "/"
)

var (
	requiredRoles = []string{"INTERNAL_DOCS"}
)
