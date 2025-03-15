package presenters

import (
	"bf_me/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Exercise struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	TitleEn   string    `json:"titleEn"`
	TitleRu   string    `json:"titleRu"`
	Filename  string    `json:"filename"`
	Tips      []string  `json:"tips"`
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
		Tips:      e.Tips,
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
			Tips:      e.Tips,
		}
	}
	return exercises
}

type Session struct {
	Token pgtype.UUID `json:"token"`
}

func (p *Presenter) Session(s *models.Session) *Session {
	return &Session{Token: s.ID}
}

type Block struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	TitleEn       string    `json:"titleEn"`
	TitleRu       string    `json:"titleRu"`
	TotalDuration uint8     `json:"totalDuration"` // minutes
	OnTime        uint8     `json:"onTime"`        // seconds
	RelaxTime     uint8     `json:"relaxTime"`     // seconds
	Draft         bool      `json:"draft"`
	ExercisesIds  []uint    `json:"exercisesIds,omitempty;"`
}

func (p *Presenter) buildBlock(block *models.Block) *Block {
	exercisesIds := make([]uint, len(block.Exercises))
	for i, e := range block.Exercises {
		exercisesIds[i] = e.ID
	}
	return &Block{
		ID:            block.ID,
		CreatedAt:     block.CreatedAt,
		TitleEn:       block.TitleEn,
		TitleRu:       block.TitleRu,
		TotalDuration: block.TotalDuration,
		OnTime:        block.OnTime,
		RelaxTime:     block.RelaxTime,
		Draft:         block.Draft,
		ExercisesIds:  exercisesIds,
	}
}

func (p *Presenter) Block(block *models.Block) *Block {
	return p.buildBlock(block)
}

func (p *Presenter) Blocks(bs []*models.Block) []*Block {
	exercises := make([]*Block, len(bs))
	for i, b := range bs {
		exercises[i] = p.buildBlock(b)
	}
	return exercises
}
