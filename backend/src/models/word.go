package models

type CreateWordRequest struct {
	Name      string     `json:"name"`
	WordInfos []WordInfo `json:"wordInfos"`
	UserID    int        `json:"userId"`
}

type CreateWordResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Word struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	RegistrationCount int        `json:"registrationCount"`
	WordInfos         []WordInfo `json:"wordInfos"`
	IsRegistered      bool       `json:"isRegistered"`
	AttentionLevel    int        `json:"attentionLevel"`
	TestCount         int        `json:"testCount"`
	CheckCount        int        `json:"checkCount"`
}

type WordInfo struct {
	ID             int            `json:"id"`
	PartOfSpeech   PartOfSpeech   `json:"partOfSpeech"`
	PartOfSpeechID int            `json:"partOfSpeechId"`
	JapaneseMeans  []JapaneseMean `json:"japaneseMeans"`
}

type PartOfSpeech struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JapaneseMean struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AllWordListRequest struct {
	UserID int    `json:"userId"`
	Search string `json:"search"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type AllWordListResponse struct {
	Words      []Word `json:"words"`
	TotalPages int    `json:"totalPages"`
}

type WordShowRequest struct {
	WordID int `json:"id" binding:"required"`
	UserID int `json:"userId"`
}

type WordShowResponse struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	RegistrationCount int        `json:"registrationCount"`
	WordInfos         []WordInfo `json:"wordInfos"`
	IsRegistered      bool       `json:"isRegistered"`
	AttentionLevel    int        `json:"attentionLevel"`
	TestCount         int        `json:"testCount"`
	CheckCount        int        `json:"checkCount"`
	Memo              string     `json:"memo"`
}

type DeleteWordRequest struct {
	WordID int `json:"id" binding:"required"`
	UserID int `json:"userId"`
}

type DeleteWordResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type RegisterWordRequest struct {
	WordID       int  `json:"wordId" binding:"required"`
	UserID       int  `json:"userId"`
	IsRegistered bool `json:"isRegistered"`
}

type RegisterWordResponse struct {
	Name              string `json:"name"`
	IsRegistered      bool   `json:"isRegistered"`
	RegistrationCount int    `json:"registrationCount"`
	Message           string `json:"message"`
}

type RegisteredWordCountRequest struct {
	WordID       int  `json:"wordId" binding:"required"`
	UserID       int  `json:"userId"`
	IsRegistered bool `json:"isRegistered"`
}

type RegisteredWordCountResponse struct {
	RegistrationCount int `json:"registrationCount"`
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
