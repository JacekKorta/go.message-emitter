package message


type Body struct {
	Password string `json:"password"`
}

type Header struct {
	Data string `json:"data"`
}

type Message struct {
	Header Header `json:"header"`	
	Body Body `json:"body"`
}