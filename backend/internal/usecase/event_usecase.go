package usecase

import "github.com/aswinsreeraj/evntx/internal/repository"

type EventUsecase struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) *EventUsecase {
	return &EventUsecase{repo: repo}
}

func (u *EventUsecase) ListEvents(city string, page int, limit int) (interface{}, int64, error) {
	return u.repo.ListLiveEvents(city, page, limit)
}

func (u *EventUsecase) GetEvent(slug string) (interface{}, interface{}, interface{}, error) {

	event, err := u.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, nil, nil, err
	}

	details, err := u.repo.GetEventDetails(event.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	personnels, err := u.repo.GetEventPersonnels(event.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	return event, details, personnels, nil
}
