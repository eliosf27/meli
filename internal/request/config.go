package request

type Config struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type KeyConfig struct {
	Key string `json:"key" validate:"required"`
}

type FindConfig struct {
	Config []KeyConfig `json:"config"`
}

type SaveConfig struct {
	Config []Config `json:"config"`
}
