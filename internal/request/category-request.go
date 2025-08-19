package request

type CategoryRequest struct {
	Name     string `json:"name" validate:"required,unique"`
	ParentID string `json:"parent_id" validate:"required"`
}
