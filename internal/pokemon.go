package internal

type PokemonInfo struct {
	Name    string `json:"name"`
	BaseEXP int    `json:"base_experience"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Stats   []Stat `json:"stats"`
	Types   []Type `json:"types"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	statName     `json:"stat"`
}

type statName struct {
	Name string `json:"name"`
}

type Type struct {
	typeName    `json:"type"`
}

type typeName struct {
	Name string `json:"name"`
}
