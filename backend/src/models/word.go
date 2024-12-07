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
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	WordInfos    []WordInfo `json:"wordInfos"`
	IsRegistered bool       `json:"isRegistered"`
	TestCount    int        `json:"testCount"`
	CheckCount   int        `json:"checkCount"`
	Memo         string     `json:"memo"`
}

type RegisterWordRequest struct {
	WordID       int  `json:"wordId" binding:"required"`
	UserID       int  `json:"userId"`
	IsRegistered bool `json:"isRegistered"`
}

type RegisterWordResponse struct {
	Name         string `json:"name"`
	IsRegistered bool   `json:"isRegistered"`
	Message      string `json:"message"`
}

type SaveMemoRequest struct {
	WordID int    `json:"wordId" binding:"required"`
	UserID int    `json:"userId"`
	Memo   string `json:"memo"`
}

type SaveMemoResponse struct {
	Name    string `json:"name"`
	Memo    string `json:"memo"`
	Message string `json:"message"`
}
