package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
)

var (
	ErrBlockNotReady        = errors.New("block is not ready to be published\nadd more exercises")
	ErrBlockFullOfExercises = errors.New("block full of exercises\n check it and be ready to publish it")
	ErrExerciseDeleted      = errors.New("exercise was deleted\nchoose another one")
)

type BlocksUseCase struct {
	storage *storage.Storage
}

func NewBlocksUseCase(st *storage.Storage) *BlocksUseCase {
	return &BlocksUseCase{storage: st}
}

func (buc *BlocksUseCase) List(req *requests.FilterBlocksRequestBody) ([]*models.Block, error) {
	var blocks []*models.Block
	updatedAtSql := fmt.Sprintf("updated_at %s", req.UpdatedAt)

	if req.Suggestion != "" {
		result := buc.storage.DB.Where("title_en ILIKE ? OR title_ru ILIKE ?", "%"+req.Suggestion+"%", "%"+req.Suggestion+"%").Find(&blocks)
		return blocks, result.Error
	}

	if req.BlockType == "draft" {
		result := buc.storage.DB.Where("draft = ?", true).Order(updatedAtSql).Preload("ExerciseBlocks").Find(&blocks)
		return blocks, result.Error
	}

	if req.BlockType == "ready" {
		result := buc.storage.DB.Where("draft = ?", false).Order(updatedAtSql).Preload("ExerciseBlocks").Find(&blocks)
		return blocks, result.Error
	}

	result := buc.storage.DB.Order(updatedAtSql).Preload("ExerciseBlocks").Find(&blocks)
	return blocks, result.Error
}

func (buc *BlocksUseCase) AddBlockExercise(blockID, exerciseID uint) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, blockID)
	if result.Error != nil {
		return nil, result.Error
	}
	if !block.Draft {
		return nil, errors.New("block is not draft\ncannot add exercise")
	}
	//check if exercises count is not reached its highest level
	full := buc.checkBlockFullOfExercises(&block)
	if full {
		return nil, ErrBlockFullOfExercises
	}

	var exercise models.Exercise
	result = buc.storage.DB.First(&exercise, exerciseID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrExerciseDeleted
	}
	if result.Error != nil {
		return nil, result.Error
	}

	var ebs []models.ExerciseBlock
	result = buc.storage.DB.Where("block_id = ?", blockID).Find(&ebs)
	if result.Error != nil {
		return nil, result.Error
	}

	nextOrder := buc.findNextOrder(block.ExerciseBlocks)
	eb := models.ExerciseBlock{
		ExerciseID:    exerciseID,
		BlockID:       blockID,
		ExerciseOrder: nextOrder,
	}
	result = buc.storage.DB.Create(&eb)
	if result.Error != nil {
		return nil, result.Error
	}

	result = buc.storage.DB.Preload("ExerciseBlocks").First(&block, blockID)
	return &block, result.Error
}

func (buc *BlocksUseCase) RemoveBlockExercise(blockID, exerciseID uint) (*models.Block, error) {
	var exerciseBlockRelation models.ExerciseBlock
	result := buc.storage.DB.First(&exerciseBlockRelation, "block_id = ? AND exercise_id = ?", blockID, exerciseID)
	if result.Error != nil {
		return nil, result.Error
	}

	result = buc.storage.DB.Unscoped().Delete(&exerciseBlockRelation)
	if result.Error != nil {
		return nil, result.Error
	}

	var block models.Block
	result = buc.storage.DB.Preload("ExerciseBlocks").First(&block, blockID)
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

	return &block, nil
}

func (buc *BlocksUseCase) Create(req *requests.BlockRequestBody) (*models.Block, error) {
	var block models.Block
	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return nil, err
	}
	buc.fitTiming(updatedBlock)

	result := buc.storage.DB.Create(&updatedBlock)
	return updatedBlock, result.Error
}

func (buc *BlocksUseCase) fitTiming(block *models.Block) {
	//TotalDuration
	if block.TotalDuration > 60 {
		block.TotalDuration = 60
	}
	if block.TotalDuration < 10 {
		block.TotalDuration = 10
	}

	//relaxTime
	if block.RelaxTime > 30 {
		block.RelaxTime = 30
	}

	//onTime
	if block.OnTime > 60 {
		block.OnTime = 60
	}
	if block.OnTime < 20 {
		block.OnTime = 20
	}

	exercisesCount := (int(block.TotalDuration) * 60) / int(block.OnTime+block.RelaxTime)
	//if all counting fits with each other
	if exercisesCount*int(block.OnTime+block.RelaxTime) == int(block.TotalDuration)*60 {
		return
	}

	block.OnTime = uint8(math.Ceil(float64(block.OnTime)/10) * 10)
	block.RelaxTime = 60 - block.OnTime
}

func (buc *BlocksUseCase) Find(id int) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, id)
	return &block, result.Error
}

func (buc *BlocksUseCase) Update(id int, req *requests.BlockRequestBody) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, id)
	if result.Error != nil {
		return nil, result.Error
	}

	//fmt.Printf("\nreq: %+v\n\n", req)

	//if req.Draft == false && req.TitleEn == "" && req.TitleRu == "" {
	//	full := buc.checkBlockFullOfExercises(&block)
	//	if full == false {
	//		return nil, ErrBlockNotReady
	//	}
	//}

	buc.fitTiming(&block)

	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return nil, err
	}

	result = buc.storage.DB.Save(&updatedBlock)
	return updatedBlock, result.Error
}

func (buc *BlocksUseCase) ToggleDraft(id int) (*models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if block.Draft {
		block.Draft = false
	} else {
		block.Draft = true
	}

	result = buc.storage.DB.Save(&block)
	return &block, result.Error
}

func (buc *BlocksUseCase) Delete(id int) error {
	var block *models.Block
	result := buc.storage.DB.First(&block, id)

	result = buc.storage.DB.Delete(&models.Block{}, id)
	return result.Error
}

func (buc *BlocksUseCase) checkBlockFullOfExercises(block *models.Block) bool {
	return int(block.TotalDuration)*60 == len(block.ExerciseBlocks)*int(block.OnTime+block.RelaxTime)
}
