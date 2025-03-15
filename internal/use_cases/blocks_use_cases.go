package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
)

type BlocksUseCase struct {
	storage *storage.Storage
}

func NewBlocksUseCase(st *storage.Storage) *BlocksUseCase {
	return &BlocksUseCase{storage: st}
}

func (euc *BlocksUseCase) List(req *requests.FilterBlocksRequestBody) ([]*models.Block, error) {
	var blocks []*models.Block
	result := euc.storage.DB.Preload("Exercises").Where("draft = ?", req.Draft).Order("updated_at DESC").Find(&blocks)
	return blocks, result.Error
}

func (euc *BlocksUseCase) updateBlock(block models.Block, req requests.BlockRequestBody) (*models.Block, error) {
	if len(req.ExercisesIds) != 0 {
		var existingExercises []models.Exercise
		result := euc.storage.DB.Where("id IN ?", req.ExercisesIds).Find(&existingExercises)
		if result.Error != nil {
			return nil, result.Error
		}

		block.Exercises = existingExercises
	}

	if req.TitleRu != "" {
		block.TitleRu = req.TitleRu
	}
	if req.TitleEn != "" {
		block.TitleEn = req.TitleEn
	}
	if req.TotalDuration != 0 {
		block.TotalDuration = req.TotalDuration
	}

	if req.OnTime != 0 {
		block.OnTime = req.OnTime
	}
	block.RelaxTime = req.RelaxTime
	if req.Draft != block.Draft {
		block.Draft = req.Draft
	}
	return &block, nil
}

func (buc *BlocksUseCase) Create(req *requests.BlockRequestBody) (*models.Block, error) {
	var block models.Block
	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return nil, err
	}

	result := buc.storage.DB.Create(&updatedBlock)
	return updatedBlock, result.Error
}

func (buc *BlocksUseCase) Find(id int) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("Exercises").First(&block, id)
	return &block, result.Error
}

func (buc *BlocksUseCase) Update(id int, req *requests.BlockRequestBody) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("Exercises").First(&block, id)
	if result.Error != nil {
		return nil, result.Error
	}

	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return nil, err
	}

	result = buc.storage.DB.Save(&updatedBlock)
	return updatedBlock, result.Error
}

func (buc *BlocksUseCase) Delete(id int) error {
	var block *models.Block
	result := buc.storage.DB.First(&block, id)

	result = buc.storage.DB.Delete(&models.Block{}, id)
	return result.Error
}
