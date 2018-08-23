package inventory

type Entry struct {
	Subject string `json:"subject"`
	Change  int    `json:"change"`
}
