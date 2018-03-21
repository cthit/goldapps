package goldapps

type User struct {
	Cid        string `json:"cid"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Nick       string `json:"nick"`
	Mail       string `json:"mail"` // Backup email?  must be investigated
}

func (user User) equals(other User) bool {

	return true
}
