package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var (
	ErrTrainingDeleted = errors.New("exercise was deleted\nchoose another one")
)

type TrainingsUseCase struct {
	storage *storage.Storage
}

func NewTrainingsUseCase(st *storage.Storage) *TrainingsUseCase {
	return &TrainingsUseCase{storage: st}
}

func (tuc *TrainingsUseCase) List(req *requests.FilterRequestBody) ([]*models.Training, error) {
	var trainings []*models.Training
	updatedAtSql := fmt.Sprintf("updated_at %s", req.UpdatedAt)

	if req.Suggestion != "" {
		result := tuc.storage.DB.Where("title_en ILIKE ? OR title_ru ILIKE ?", "%"+req.Suggestion+"%", "%"+req.Suggestion+"%").Find(&trainings)
		return trainings, result.Error
	}

	if req.BlockType == "draft" {
		result := tuc.storage.DB.Where("draft = ?", true).Order(updatedAtSql).Preload("TrainingBlocks").Find(&trainings)
		return trainings, result.Error
	}

	if req.BlockType == "ready" {
		result := tuc.storage.DB.Where("draft = ?", false).Order(updatedAtSql).Preload("TrainingBlocks").Find(&trainings)
		return trainings, result.Error
	}

	result := tuc.storage.DB.Order(updatedAtSql).Preload("TrainingBlocks").Preload("Blocks").Find(&trainings)
	return trainings, result.Error
}

func (tuc *TrainingsUseCase) AddTrainingBlock(trainingID, blockID uint) (*models.Training, []models.Block, error) {
	var training models.Training
	result := tuc.storage.DB.Preload("TrainingBlocks").First(&training, trainingID)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}
	if !training.Draft {
		return nil, []models.Block{}, errors.New("block is not draft\ncannot add exercise")
	}

	var block models.Block
	result = tuc.storage.DB.First(&block, blockID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, []models.Block{}, ErrTrainingDeleted
	}
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	var tbs []models.TrainingBlock
	result = tuc.storage.DB.Where("training_id = ?", trainingID).Find(&tbs)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	nextOrder := tuc.findNextOrder(training.TrainingBlocks)
	eb := models.TrainingBlock{
		TrainingID: trainingID,
		BlockID:    blockID,
		BlockOrder: nextOrder,
	}
	result = tuc.storage.DB.Create(&eb)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	result = tuc.storage.DB.Preload("TrainingBlocks").Preload("Blocks").First(&training, trainingID)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	blockIds := make([]uint, len(training.TrainingBlocks))
	for i, t := range training.TrainingBlocks {
		blockIds[i] = t.BlockID
	}

	var blocks []models.Block
	result = tuc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").Where("id IN ?", blockIds).Find(&blocks)

	return &training, blocks, result.Error
}

func (tuc *TrainingsUseCase) RemoveTrainingBlock(trainingID, blockID uint) (*models.Training, []models.Block, error) {
	var trainingBlockRelation models.TrainingBlock
	result := tuc.storage.DB.First(&trainingBlockRelation, "training_id = ? AND block_id = ?", trainingID, blockID)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	result = tuc.storage.DB.Unscoped().Delete(&trainingBlockRelation)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	var training models.Training
	result = tuc.storage.DB.Preload("TrainingBlocks").Preload("Blocks").First(&training, trainingID)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	blockIds := make([]uint, len(training.TrainingBlocks))
	for i, t := range training.TrainingBlocks {
		blockIds[i] = t.BlockID
	}

	var blocks []models.Block
	result = tuc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").Where("id IN ?", blockIds).Find(&blocks)

	return &training, blocks, result.Error
}

func (tuc *TrainingsUseCase) findNextOrder(tbs []models.TrainingBlock) uint {
	var order uint = 0
	for _, t := range tbs {
		if t.BlockOrder >= order {
			order = t.BlockOrder + 1
		}
	}

	return order
}

func (tuc *TrainingsUseCase) updateTraining(tr models.Training, req requests.TrainingRequestBody) (*models.Training, error) {
	if req.TitleRu != "" {
		tr.TitleRu = req.TitleRu
	}
	if req.TitleEn != "" {
		tr.TitleEn = req.TitleEn
	}

	return &tr, nil
}

func (tuc *TrainingsUseCase) Create(req *requests.TrainingRequestBody) (*models.Training, error) {
	var training models.Training
	updatedTr, err := tuc.updateTraining(training, *req)
	if err != nil {
		return nil, err
	}

	result := tuc.storage.DB.Create(&updatedTr)
	return updatedTr, result.Error
}

func (tuc *TrainingsUseCase) Find(id int) (*models.Training, []models.Block, error) {
	var training models.Training
	result := tuc.storage.DB.Preload("TrainingBlocks").First(&training, id)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	blockIds := make([]uint, len(training.TrainingBlocks))
	for i, t := range training.TrainingBlocks {
		blockIds[i] = t.BlockID
	}

	var blocks []models.Block
	result = tuc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").Where("id IN ?", blockIds).Find(&blocks)

	return &training, blocks, result.Error
}

func (tuc *TrainingsUseCase) Update(id int, req *requests.TrainingRequestBody) (*models.Training, []models.Block, error) {
	var training models.Training
	result := tuc.storage.DB.Preload("TrainingBlocks").Preload("Blocks").First(&training, id)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	updatedTr, err := tuc.updateTraining(training, *req)
	if err != nil {
		return nil, []models.Block{}, err
	}

	result = tuc.storage.DB.Save(&updatedTr)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	blockIds := make([]uint, len(training.TrainingBlocks))
	for i, t := range training.TrainingBlocks {
		blockIds[i] = t.BlockID
	}

	var blocks []models.Block
	result = tuc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").Where("id IN ?", blockIds).Find(&blocks)

	return &training, blocks, result.Error
}

func (tuc *TrainingsUseCase) ToggleDraft(id int) (*models.Training, []models.Block, error) {
	var training models.Training
	result := tuc.storage.DB.Preload("TrainingBlocks").Preload("Blocks").First(&training, id)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	if training.Draft {
		training.Draft = false
	} else {
		training.Draft = true
	}

	result = tuc.storage.DB.Save(&training)
	if result.Error != nil {
		return nil, []models.Block{}, result.Error
	}

	blockIds := make([]uint, len(training.TrainingBlocks))
	for i, t := range training.TrainingBlocks {
		blockIds[i] = t.BlockID
	}

	var blocks []models.Block
	result = tuc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").Where("id IN ?", blockIds).Find(&blocks)

	return &training, blocks, result.Error
}

func (tuc *TrainingsUseCase) Delete(id int) error {
	var training *models.Training
	result := tuc.storage.DB.First(&training, id)

	result = tuc.storage.DB.Delete(&models.Training{}, id)
	return result.Error
}
