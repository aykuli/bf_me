package presenters

import (
	"bf_me/internal/models"
	"time"
)

type Exercise struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	TitleEn   string    `json:"titleEn"`
	TitleRu   string    `json:"titleRu"`
	Filename  string    `json:"filename"`
	Tags      []string  `json:"tagIds,omitempty;"`
}

type Presenter struct {
}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) Exercise(e *models.Exercise) *Exercise {
	//todo tags
	return &Exercise{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		TitleEn:   e.TitleEn,
		TitleRu:   e.TitleRu,
		Filename:  e.Filename,
	}
}

func (p *Presenter) Exercises(es []*models.Exercise) []*Exercise {
	exercises := make([]*Exercise, len(es))
	for i, e := range es {
		exercises[i] = &Exercise{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			TitleEn:   e.TitleEn,
			TitleRu:   e.TitleRu,
			Filename:  e.Filename,
		}
	}
	return exercises
}
