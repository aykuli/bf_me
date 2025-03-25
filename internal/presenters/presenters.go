package presenters

import (
	"bf_me/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"slices"
)

type Exercise struct {
	ID        uint     `json:"id"`
	CreatedAt string   `json:"createdAt"`
	TitleEn   string   `json:"titleEn"`
	TitleRu   string   `json:"titleRu"`
	Filename  string   `json:"filename"`
	Tips      []string `json:"tips"`
	Tags      []string `json:"tagIds,omitempty;"`
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
		CreatedAt: e.CreatedAt.Format("January 2, 2006"),
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
			CreatedAt: e.CreatedAt.Format("January 2, 2006"),
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
	ID            uint            `json:"id"`
	CreatedAt     string          `json:"createdAt"`
	TitleEn       string          `json:"titleEn"`
	TitleRu       string          `json:"titleRu"`
	TotalDuration uint8           `json:"totalDuration"` // minutes
	OnTime        uint8           `json:"onTime"`        // seconds
	RelaxTime     uint8           `json:"relaxTime"`     // seconds
	Draft         bool            `json:"draft"`
	Exercises     []BlockExercise `json:"exercises,omitempty;"`
}

type BlockExercise struct {
	ID       uint   `json:"id"` // exercise id
	Order    uint   `json:"order"`
	Side     string `json:"side"`
	TitleEn  string `json:"titleEn"`
	TitleRu  string `json:"titleRu"`
	Filename string `json:"filename"`
}

func (p *Presenter) Block(block *models.Block) *Block {
	return &Block{
		ID:            block.ID,
		CreatedAt:     block.CreatedAt.Format("January 2, 2006"),
		TitleEn:       block.TitleEn,
		TitleRu:       block.TitleRu,
		TotalDuration: block.TotalDuration,
		OnTime:        block.OnTime,
		RelaxTime:     block.RelaxTime,
		Draft:         block.Draft,
		Exercises:     p.buildBlockExerciseIds(block),
	}
}

func (p *Presenter) buildBlockExerciseIds(block *models.Block) []BlockExercise {
	slices.SortFunc(block.ExerciseBlocks, func(a, b models.ExerciseBlock) int {
		return int(a.ExerciseOrder - b.ExerciseOrder)
	})

	var arr = make([]BlockExercise, len(block.ExerciseBlocks))
	for i, eb := range block.ExerciseBlocks {
		exerciseID := eb.ExerciseID
		exercise := p.takeExerciseByID(block.Exercises, exerciseID)
		arr[i] = BlockExercise{
			ID:       eb.ExerciseID,
			Order:    uint(i),
			Side:     eb.Side,
			TitleEn:  exercise.TitleEn,
			TitleRu:  exercise.TitleRu,
			Filename: exercise.Filename,
		}
	}
	return arr
}

func (p *Presenter) takeExerciseByID(exercises []models.Exercise, exerciseID uint) *models.Exercise {
	for _, e := range exercises {
		if e.ID == exerciseID {
			return &e
		}
	}
	return nil
}

func (p *Presenter) Blocks(bs []*models.Block) []*Block {
	exercises := make([]*Block, len(bs))
	for i, block := range bs {
		exercises[i] = &Block{
			ID:            block.ID,
			CreatedAt:     block.CreatedAt.Format("January 2, 2006"),
			TitleEn:       block.TitleEn,
			TitleRu:       block.TitleRu,
			TotalDuration: block.TotalDuration,
			OnTime:        block.OnTime,
			RelaxTime:     block.RelaxTime,
			Draft:         block.Draft,
		}
	}
	return exercises
}
