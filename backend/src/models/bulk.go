package models

// BulkTokenizeRequest is the input payload for the
// “bulk tokenize” endpoint. `Text` may contain up to 3 000 bytes.
type BulkTokenizeRequest struct {
	Text string `json:"text" binding:"required,max=3000"`
}

// BulkTokenizeResponse groups the results of a bulk tokenization.
//   - Candidates    — all tokens detected in the text
//   - Registered    — tokens already saved by the user
//   - NotExistWord  — tokens that don’t exist in the master dictionary
type BulkTokenizeResponse struct {
	Candidates   []string `json:"candidates"`
	Registered   []string `json:"registered"`
	NotExistWord []string `json:"not_exists"`
}

// BulkRegisterRequest is used to register up to 200 words at once.
type BulkRegisterRequest struct {
	Words []string `json:"words" binding:"required,min=1,max=200,dive,bulk"`
}

// FailedWord describes a word that could not be registered and the reason.
type FailedWord struct {
	Word   string `json:"word"`
	Reason string `json:"reason"`
}

// BulkRegisterResponse summarises the outcome of a bulk word-registration call.
type BulkRegisterResponse struct {
	Success []string     `json:"success"` // words that were registered
	Failed  []FailedWord `json:"failed"`  // words that failed & why
}
