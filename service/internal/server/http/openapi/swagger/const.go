// const.go
package swagger

import "io/fs"

const (
	defaultLoginURL = "http://10.70.20.95:8134/v1/login"
)

var (
	defaultRequiredRoles = []string{"INTERNAL_DOCS"}
)

type Config struct {
	Base          string
	DocsFS        fs.FS
	LoginURL      string
	CookieName    string
	CookiePath    string
	RequiredRoles []string
	ProjectPrefix string

	ServiceName string
}

func (c *Config) normalize() {
	c.Base = cleanBase(c.Base)
	if c.CookieName == "" {
		c.CookieName = "swagger_session"
	}
	if c.CookiePath == "" {
		if c.Base == "" {
			c.CookiePath = "/docs/"
		} else {
			c.CookiePath = c.Base + "/docs/"
		}
	}
	if c.LoginURL == "" {
		c.LoginURL = defaultLoginURL
	}
	if c.RequiredRoles == nil {
		c.RequiredRoles = defaultRequiredRoles
	}
	c.ProjectPrefix = cleanBase(c.ProjectPrefix)
}
