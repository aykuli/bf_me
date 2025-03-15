package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"errors"
)

type BlocksUseCase struct {
	storage *storage.Storage
}

func NewBlocksUseCase(st *storage.Storage) *BlocksUseCase {
	return &BlocksUseCase{storage: st}
}

func (buc *BlocksUseCase) List(req *requests.FilterBlocksRequestBody) ([]*models.Block, error) {
	var blocks []*models.Block
	result := buc.storage.DB.Where("draft = ?", req.Draft).Order("updated_at DESC").Find(&blocks)
	return blocks, result.Error
}

// todo database transactions
func (buc *BlocksUseCase) AddBlockExercise(blockID, exerciseID uint) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.First(&block, blockID)
	if result.Error != nil {
		return nil, result.Error
	}
	if !block.Draft {
		return nil, errors.New("block is not draft")
	}

	var ebs []models.ExerciseBlock
	result = buc.storage.DB.Where("block_id = ?", blockID).Find(&ebs)
	if result.Error != nil {
		return nil, result.Error
	}

	nextOrder := buc.findNextOrder(ebs)
	eb := models.ExerciseBlock{
		ExerciseID:    exerciseID,
		BlockID:       blockID,
		ExerciseOrder: nextOrder,
	}
	result = buc.storage.DB.Create(&eb)
	if result.Error != nil {
		return nil, result.Error
	}

	result = buc.storage.DB.First(&block, blockID)
	return &block, result.Error
}

func (buc *BlocksUseCase) findNextOrder(ebs []models.ExerciseBlock) uint {
	var order uint = 0
	for _, e := range ebs {
		if e.ExerciseOrder >= order {
			order = e.ExerciseOrder + 1
		}
	}

	return order
}

func (buc *BlocksUseCase) updateBlock(block models.Block, req requests.BlockRequestBody) (*models.Block, error) {
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
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, id)
	return &block, result.Error
}

func (buc *BlocksUseCase) Update(id int, req *requests.BlockRequestBody) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.First(&block, id)
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
