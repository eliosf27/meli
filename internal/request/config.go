package request

type UpdateStorage struct {
	Storage string `json:"storage" validate:"required"`
}
