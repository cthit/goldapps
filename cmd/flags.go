package main

import "flag"

type flagStruct struct {
	from          string
	to            string
	interactive   bool
	noInteraction bool
	dryRun        bool
}

var flags = flagStruct{}

func loadFlags() {
	flag.StringVar(&flags.from, "from", "ldap", "Set the group source to 'ldap', 'gapps' or '*.json'. In case of gapps config value 'gappsProvider' will be used")
	flag.StringVar(&flags.to, "to", "gapps", "Set the group consumer to 'gapps' or '*.json'")
	flag.BoolVar(&flags.dryRun, "dry", false, "Setting this flag will cause the application to only print information and not update any groups")
	flag.BoolVar(&flags.noInteraction, "y", false, "Setting this flag will cause the application to not ask for any user confirmation")
	flag.BoolVar(&flags.interactive, "i", false, "Setting this flag will cause the application to ask the user for input in every stage ")
	flag.Parse()
}
