package schema

// Summary is the minimal structural description extracted from a file.
type Summary struct {
	Format    string
	Separator rune
	HasHeader bool
	Columns   []string
	Rows      int
}

// Plan describes a minimal transformation to execute.
type Plan struct {
	Objective  string   `json:"objective"`
	Source     string   `json:"source"`
	Target     string   `json:"target"`
	InputPath  string   `json:"input_path"`
	OutputPath string   `json:"output_path"`
	Operation  string   `json:"operation"`
	Steps      []string `json:"steps"`
	Summary    Summary  `json:"summary"`
}
