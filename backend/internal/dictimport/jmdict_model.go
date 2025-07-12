package dictimport

type JMdictJSON struct {
	Version string    `json:"version"`
	Words   []JMEntry `json:"words"`
}

type JMEntry struct {
	ID    string  `json:"id"`
	Kana  []Kana  `json:"kana"`
	Kanji []Kanji `json:"kanji"`
	Sense []Sense `json:"sense"`
}

type Kana struct {
	Text string `json:"text"`
}

type Kanji struct {
	Text string `json:"text"`
}

type Sense struct {
	PartOfSpeech []string `json:"partOfSpeech"`
	Gloss        []Gloss  `json:"gloss"`
}

type Gloss struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}
type Options struct {
	Workers   int
	BatchSize int
}

type ImportErr struct {
	ID      string `json:"id"`
	Message string `json:"error_message"`
}
