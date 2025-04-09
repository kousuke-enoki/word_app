// backend/internal/dictimport/jmdict_model.go
package dictimport

// JMdictJSON は jmdict‑simplified（JSON 版）のトップレベル
type JMdictJSON struct {
	Version string    `json:"version"`
	Words   []JMEntry `json:"words"`
}

// JMEntry は 1 つの見出し語
type JMEntry struct {
	ID    string  `json:"id"`
	Kana  []Kana  `json:"kana"`
	Sense []Sense `json:"sense"`
}

// Kana (かな形) ― kanji は使わず英語→日本語に利用
type Kana struct {
	Text string `json:"text"`
}

// Sense ― 品詞・訳語のまとまり
type Sense struct {
	PartOfSpeech []string `json:"partOfSpeech"` // 例: ["n","vs"]
	Gloss        []Gloss  `json:"gloss"`
}

// Gloss ― 各言語の訳語
type Gloss struct {
	Lang string `json:"lang"` // "eng", "jpn" など
	Text string `json:"text"`
}
