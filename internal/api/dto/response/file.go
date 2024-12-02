package response

type FileLineContent struct {
	Content string   `json:"content"`
	End     bool     `json:"end"`
	Path    string   `json:"path"`
	Total   int      `json:"total"`
	Lines   []string `json:"lines"`
}
