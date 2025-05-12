package domain

type Choice struct {
	ID      string         `yaml:"id"`
	Text    string         `yaml:"text"`
	Next    string         `yaml:"next"`
	Effects map[string]int `yaml:"effects"` // rage, honor, karma и т.д.
}

type Scene struct {
	ID      string   `yaml:"id"`
	Text    string   `yaml:"text"`
	Choices []Choice `yaml:"choices"`
}
