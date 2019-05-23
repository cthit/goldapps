package cli

import "flag"

type flagStruct struct {
	from          string
	to            string
	additions     string
	interactive   bool
	noInteraction bool
	dryRun        bool
	onlyGroups    bool
	onlyUsers     bool
}

var flags = flagStruct{}

func loadFlags() {
	flag.StringVar(&flags.from, "from", "ldap", "Set the source to 'ldap', 'gapps' or '*.json'. In case of gapps config value 'gappsProvider' will be used")
	flag.StringVar(&flags.additions, "additions", "", "Set a json file for additional groups and users")
	flag.StringVar(&flags.to, "to", "gapps", "Set the services to 'gapps' or '*.json'")
	flag.BoolVar(&flags.dryRun, "dry", false, "Setting this flag will cause the application to only print information and not update any groups")
	flag.BoolVar(&flags.noInteraction, "y", false, "Setting this flag will cause the application to not ask for any user confirmation")
	flag.BoolVar(&flags.interactive, "i", false, "Setting this flag will cause the application to ask the user for input in every stage ")
	flag.BoolVar(&flags.onlyGroups, "groups", false, "Setting this flag will cause the application to only collect and update groups")
	flag.BoolVar(&flags.onlyUsers, "users", false, "Setting this flag will cause the application to only collect and update users ")
	flag.Parse()
}
