package ctx

import "fmt"

type GErrors []*GError

type GError struct {
	Code         uint   `json:"code,omitempty"`
	Err          error  `json:"-"`
	Domain       string `json:"domain,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"locationType,omitempty"`
	ExtendedHelp string `json:"extendedHelp,omitempty"`
	SendReport   string `json:"sendReport,omitempty"`
}

func (g *GError) Error() string {
	if g.Err != nil {
		return fmt.Sprintf("%s:%s %s, Err:%s", g.Domain, g.Reason, g.Message, g.Err.Error())
	}
	return fmt.Sprintf("%s:%s %s", g.Domain, g.Reason, g.Message)
}

func (g *GError) AppendDomain(domain string) error {
	g.Domain = domain + "." + g.Domain
	return g
}

func (c *GErrors) Append(gErr *GError) *GErrors {
	*c = append(*c, gErr)
	return c
}

func (c *GErrors) AppendDomain(domain string) *GErrors {
	for idx := range *c {
		(*c)[idx].Domain = domain + "." + (*c)[idx].Domain
	}
	return c
}

func (c *GErrors) Empty() bool {
	return len(*c) == 0
}

func (c *GErrors) Code() uint {
	if c.Empty() {
		return 0
	}
	return (*c)[0].Code
}

func (c *GErrors) Message() string {
	if c.Empty() {
		return ""
	}
	return (*c)[0].Message
}

func NewGErrors() *GErrors { return &GErrors{} }
