package ctx

type ErrorProto struct {
	Items []ErrorProtoItem
}

type ErrorProtoItem struct {
	Domain       string `json:"domain,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"location_type,omitempty"`
	ExtendedHelp string `json:"extended_help,omitempty"`
	SendReport   string `json:"send_report,omitempty"`
}

func NewErrorProto() *ErrorProto { return &ErrorProto{} }

func NewErrorProtoItem() *ErrorProtoItem { return &ErrorProtoItem{} }

func (e *ErrorProto) Add(d ErrorProtoItem) *ErrorProto {

	e.Items = append(e.Items, d)

	return e
}

func (e *ErrorProto) AsErrors() []interface{} {

	errs := make([]interface{}, len(e.Items))
	for i, v := range e.Items {
		errs[i] = v
	}
	return errs
}
