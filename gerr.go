package ctx

type GErrors []GError

type GError struct {
	Code         int    `json:"-"`
	Domain       string `json:"domain,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"locationType,omitempty"`
	ExtendedHelp string `json:"extendedHelp,omitempty"`
	SendReport   string `json:"sendReport,omitempty"`
}

func (c GErrors) Append(gerr GError) GErrors {
	c = append(c, gerr)
	return c
}

func NewGErrors() GErrors { return GErrors{} }
