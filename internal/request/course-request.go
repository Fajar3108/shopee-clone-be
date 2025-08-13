package request

import "mime/multipart"

type CourseRequest struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Thumbnail   *multipart.FileHeader `json:"thumbnail" validate:"omitempty,image"`
	CategoryID  string                `json:"category_id" validate:"required"`
	Tags        []string              `json:"tags" validate:"dive,required"`
}
