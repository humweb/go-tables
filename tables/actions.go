package tables

type ActionItems struct {
	Label  string   `json:"label"`
	Link   string   `json:"link"`
	Params []string `json:"params"`
}
