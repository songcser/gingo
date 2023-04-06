package admin

type Form struct {
	Label string
	Type  string
	Name  string
	Enum  []Enum
	Value any
}

type Header struct {
	Label string
	Name  string
}

type Tag struct {
	Label string
	Type  string
	Name  string
	Admin bool
	Enum  []Enum
}

type Enum struct {
	Key   string
	Value string
}

type SsoUser struct {
	Error     string `json:"error"`
	Retval    string `json:"retval"`
	Username  string `json:"username"`
	GivenName string `json:"given_name"`
	Email     string `json:"email"`
}
