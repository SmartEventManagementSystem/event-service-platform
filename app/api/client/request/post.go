package request

type CreatePostRequest struct {
    Title   string   `json:"title" validate:"required,min=1,max=200"`
    Content string   `json:"content" validate:"required"`
    Image   string   `json:"image"`
    Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
    Title   *string  `json:"title" validate:"omitempty,min=1,max=200"`
    Content *string  `json:"content" validate:"omitempty"`
    Image   *string  `json:"image"`
    Tags    []string `json:"tags"`
}
