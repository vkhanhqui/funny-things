package types

type Issue struct {
	Nonce      string `json:"nonce"`
	Checksum   string `json:"checksum"`
	Difficulty int    `json:"difficulty"`
}

type VerifyIssueReq struct {
	Nonce      string `json:"nonce"`
	Checksum   string `json:"checksum"`
	Difficulty int    `json:"difficulty"`
	Counter    int    `json:"counter"`
	Hash       string `json:"hash"`
}

type ParsedIssue struct {
	Nonce, Checksum, Counter, Hash []byte
	Difficulty                     int
}
