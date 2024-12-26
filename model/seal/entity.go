package seal

type Seal struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Gender string `json:"gender"`
	Weight uint8  `json:"weight"`
	Age    uint8  `json:"age"`
	Dob    string `json:"dob"`
}
