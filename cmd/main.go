package main


func init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	/*provider, err := getLDAPService(
		viper.GetString("ldap.url"),
		viper.GetString("ldap.servername"),
		viper.GetString("ldap.user"),
		viper.GetString("ldap.password"),
	)*/

	/*provider, err := admin.NewGoogleService(viper.GetString("gapps.servicekeyfile"), viper.GetString("gapps.adminaccount"))
	if err != nil {
		panic(err)
	}*/

	/*consumer, err := getGoogleService()
	if err != nil {

	}*/

	/*g, err := provider.Groups()
	if err != nil {
		panic(err)
	}

	if g != nil {
		fmt.Print(g)
	}*/
	/*
		err = consumer.UpdateGroups(g)
		if err != nil {

		}
	*/
}
