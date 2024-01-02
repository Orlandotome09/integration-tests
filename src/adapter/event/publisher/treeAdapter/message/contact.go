package message

type Contact struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phones Phones `json:"phones"`
}

type Contacts []Contact
