package requests

type CreateExerciseRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
	Tag_ids []int  `json:"tag_ids"`
}

type CreateTagRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
}
