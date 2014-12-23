package platypus

type Container struct {
	Header string `xml:"header"`
	Body   Body   `xml:"body"`
}

type Body struct {
	Data DataBlock `xml:"data_block"`
}

type DataBlock struct {
	Protocol     string     `xml:"protocol"`
	Object       string     `xml:"object"`
	Action       string     `xml:"action"`
	Username     string     `xml:"username"`
	Password     string     `xml:"password"`
	Logintype    string     `xml:"logintype"`
	Properties   string     `xml:"properties"`
	Parameters   Parameters `xml:"parameters"`
	ResponseCode string     `xml:"response_code"`
	ResponseText string     `xml:"response_text"`
	Success      int        `xml:"is_success"`
}

type Parameters struct {
	Logintype string `xml:"logintype"`
	Username  string `xml:"username"`
	Password  string `xml:"password"`
	Datatype  string `xml:"datatype"`
}
