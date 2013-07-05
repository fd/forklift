package apps

type Env struct {
	CurrentUser Account   // the current heroku account
	OwnerPool   []Account // list of heroku accounts to use when creating a nuw app
}
