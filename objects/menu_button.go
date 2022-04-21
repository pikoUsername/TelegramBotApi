package objects

type (
	MenuButton interface {
		button() // stub
	}

	MenuButtonWebApp struct {
		Type   string     `json:"type"`
		Text   string     `json:"text"`
		WebApp WebAppInfo `json:"web_app"`
	}

	MenuButtonDefault struct {
		Type string `json:"type"`
	}

	MenuButtonCommands struct {
		Type string `json:"type"`
	}
)

func (MenuButtonCommands) button() {}
func (MenuButtonDefault) button()  {}
func (MenuButtonWebApp) button()   {}
