package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"gorm.io/gorm"
)

type TagsUseCase struct {
	db *gorm.DB
}

func NewTagsUseCase(db *gorm.DB) *TagsUseCase {
	return &TagsUseCase{db}
}

// todo docs, pagination, filter by fields, fetch tags
func (euc *TagsUseCase) List() ([]*models.Tag, error) {
	var tags []*models.Tag
	result := euc.db.Order("updated_at DESC").Find(&tags)
	return tags, result.Error
}

func (euc *TagsUseCase) Create(req requests.CreateTagRequestBody) (*models.Tag, error) {
	var tag = &models.Tag{TitleEn: req.TitleEn, TitleRu: req.TitleRu}
	result := euc.db.Create(tag)
	return tag, result.Error

}
