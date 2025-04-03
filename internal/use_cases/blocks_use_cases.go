package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"slices"
)

var (
	ErrBlockNotReady        = errors.New("block is not ready to be published\nadd more exercises")
	ErrBlockCannotBeDeleted = errors.New("block cannot be deleted becase it is a part of workout")
	ErrBlockFullOfExercises = errors.New("block full of exercises\n check it and be ready to publish it")
	ErrExerciseDeleted      = errors.New("exercise was deleted\nchoose another one")
)

type BlocksUseCase struct {
	storage *storage.Storage
}

func NewBlocksUseCase(st *storage.Storage) *BlocksUseCase {
	return &BlocksUseCase{storage: st}
}

func (buc *BlocksUseCase) List(req *requests.FilterRequestBody) ([]models.Block, error) {
	var blocks []models.Block
	updatedAtSql := fmt.Sprintf("updated_at %s", req.UpdatedAt)

	if req.Suggestion != "" {
		result := buc.storage.DB.Where("title_en ILIKE ? OR title_ru ILIKE ?", "%"+req.Suggestion+"%", "%"+req.Suggestion+"%").Find(&blocks)
		return blocks, result.Error
	}

	if req.BlockType != "" {
		whereClause := ""
		if req.BlockType == "draft" {
			whereClause = "draft = true"
		}

		if req.BlockType == "ready" {
			whereClause = "draft = false"
		}

		result := buc.storage.DB.Where(whereClause).Order(updatedAtSql).Preload("ExerciseBlocks").Preload("Exercises").Find(&blocks)
		return blocks, result.Error
	}

	result := buc.storage.DB.Order(updatedAtSql).Preload("ExerciseBlocks").Preload("Exercises").Find(&blocks)
	return blocks, result.Error
}

func (buc *BlocksUseCase) AddBlockExercise(blockID, exerciseID uint, req *requests.AddBlockExerciseRequestBody) (models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").First(&block, blockID)
	if result.Error != nil {
		return block, result.Error
	}
	if !block.Draft {
		return block, errors.New("block is not draft\nyou cannot add exercise")
	}
	//check if exercises count is not reached its highest level
	full := buc.checkBlockFullOfExercises(&block)
	if full {
		return block, ErrBlockFullOfExercises
	}

	var exercise models.Exercise
	result = buc.storage.DB.First(&exercise, exerciseID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return block, ErrExerciseDeleted
	}
	if result.Error != nil {
		return block, result.Error
	}

	var ebs []models.ExerciseBlock
	result = buc.storage.DB.Where("block_id = ?", blockID).Find(&ebs)
	if result.Error != nil {
		return block, result.Error
	}

	nextOrder := buc.findNextOrder(block.ExerciseBlocks)
	side := ""
	if slices.Contains([]string{"right", "left", ""}, req.Side) {
		side = req.Side
	}
	eb := models.ExerciseBlock{
		ExerciseID:    exerciseID,
		BlockID:       blockID,
		ExerciseOrder: nextOrder,
		Side:          side,
	}
	result = buc.storage.DB.Create(&eb)
	if result.Error != nil {
		return block, result.Error
	}

	result = buc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").First(&block, blockID)
	return block, result.Error
}

func (buc *BlocksUseCase) RemoveBlockExercise(blockID, exerciseID uint) (models.Block, error) {
	var exerciseBlockRelation models.ExerciseBlock
	var block models.Block

	result := buc.storage.DB.First(&exerciseBlockRelation, "block_id = ? AND exercise_id = ?", blockID, exerciseID)
	if result.Error != nil {
		return block, result.Error
	}

	result = buc.storage.DB.Unscoped().Delete(&exerciseBlockRelation)
	if result.Error != nil {
		return block, result.Error
	}

	result = buc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").First(&block, blockID)
	return block, result.Error
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

func (buc *BlocksUseCase) updateBlock(block models.Block, req requests.BlockRequestBody) (models.Block, error) {
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

	return block, nil
}

func (buc *BlocksUseCase) Create(req *requests.BlockRequestBody) (models.Block, error) {
	var block models.Block
	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return block, err
	}
	buc.fitTiming(&updatedBlock)

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

func (buc *BlocksUseCase) Find(id int) (models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").First(&block, id)
	return block, result.Error
}

func (buc *BlocksUseCase) Update(id int, req *requests.BlockRequestBody) (models.Block, error) {
	var block models.Block
	result := buc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").First(&block, id)
	if result.Error != nil {
		return block, result.Error
	}

	buc.fitTiming(&block)
	updatedBlock, err := buc.updateBlock(block, *req)
	if err != nil {
		return block, err
	}

	result = buc.storage.DB.Save(&updatedBlock)
	return updatedBlock, result.Error
}

func (buc *BlocksUseCase) ToggleDraft(id int) (models.Block, error) {
	var block models.Block
	result := buc.storage.DB.First(&block, id)
	if result.Error != nil {
		return block, result.Error
	}

	result = buc.storage.DB.Model(&block).Update("draft", !block.Draft)
	if result.Error != nil {
		return block, result.Error
	}

	result = buc.storage.DB.Preload("ExerciseBlocks").Preload("Exercises").First(&block, id)
	return block, result.Error
}

func (buc *BlocksUseCase) Delete(id int) error {
	//check if block has related workout
	var trainingBlocks []models.TrainingBlock
	result := buc.storage.DB.Where("block_id = ?", id).Find(&trainingBlocks)
	if result.Error != nil {
		return result.Error
	}

	for _, tr := range trainingBlocks {
		var training *models.Training
		result = buc.storage.DB.Find(&training, tr.TrainingID)
		if result.Error != nil {
			return result.Error
		}

		//check if training was deleted
		//		if not deleted -> throw an error
		//		else continue
		deletedValue, err := training.DeletedAt.Value()
		if err != nil {
			return err
		}
		if deletedValue == nil {
			return errors.New(fmt.Sprintf("block cannot be deleted because it is a part of the workout with id=%d", training.ID))
		}
	}

	//we should delete all training and exercise relations
	var exerciseBlockRelations []models.ExerciseBlock
	result = buc.storage.DB.Find(&exerciseBlockRelations, "block_id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if len(exerciseBlockRelations) != 0 {
		result = buc.storage.DB.Unscoped().Delete(&exerciseBlockRelations)
		if result.Error != nil {
			return result.Error
		}
	}
	var block *models.Block
	result = buc.storage.DB.First(&block, id)
	if result.Error != nil {
		return result.Error
	}

	result = buc.storage.DB.Delete(&models.Block{}, id)
	return result.Error
}

func (buc *BlocksUseCase) checkBlockFullOfExercises(block *models.Block) bool {
	return int(block.TotalDuration)*60 == len(block.ExerciseBlocks)*int(block.OnTime+block.RelaxTime)
}
