package models

// WordResponse 構造体でレスポンスを定義
type Word struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	WordInfos []WordInfo `json:"wordInfos"`
}

type WordInfo struct {
	ID            int            `json:"id"`
	PartOfSpeech  PartOfSpeech   `json:"partOfSpeech"`
	JapaneseMeans []JapaneseMean `json:"japaneseMeans"`
}

type PartOfSpeech struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JapaneseMean struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WordResponse struct {
	Name               string     `json:"name"`
	WordInfos          []WordInfo `json:"wordInfos"`
	IsRegistered       bool       `json:"isRegistered"`
	TestCount          int        `json:"testCount"`
	CheckCount         int        `json:"checkCount"`
	RegistrationActive bool       `json:"registrationActive"`
	Memo               string     `json:"memo"`
}
