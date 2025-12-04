package domain

type Service struct {
	repo  Repo
	event EventPublisher
}

type Repo interface {
}

type EventPublisher interface {
}
