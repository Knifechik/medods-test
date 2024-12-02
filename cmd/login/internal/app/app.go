package app

type App struct {
	repo   Repo
	auth   Auth
	notify Notification
}

func New(repo Repo, auth Auth, notify Notification) *App {
	return &App{
		repo:   repo,
		auth:   auth,
		notify: notify,
	}
}
