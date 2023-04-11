package admin

type Form struct {
	Label   string
	Type    string
	Name    string
	Enum    []Enum
	Value   any
	Disable bool
}

type Header struct {
	Label string
	Name  string
}

type Tag struct {
	Label   string
	Type    string
	Name    string
	Admin   bool
	Enum    []Enum
	Disable bool
}

type Enum struct {
	Key   string
	Value string
}
