package google

type Group struct {
	Name    string
	Email   string
	Members *[]Member
	Alias   *[]string
}
