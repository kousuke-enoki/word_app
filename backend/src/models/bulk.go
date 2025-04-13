package models

type BulkTokenizeRequest struct {
	Text string `json:"text" binding:"required,max=3000"`
}

type BulkTokenizeResponse struct {
	Candidates   []string `json:"candidates"`
	Registered   []string `json:"registered"`
	NotExistWord []string `json:"not_exists"`
}

type BulkRegisterRequest struct {
	Words []string `json:"words" binding:"required,min=1,max=200,dive,bulk"`
}

type FailedWord struct {
	Word   string `json:"word"`
	Reason string `json:"reason"`
}

type BulkRegisterResponse struct {
	Success []string     `json:"success"`
	Failed  []FailedWord `json:"failed"`
}
