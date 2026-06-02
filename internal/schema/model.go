package schema

// Summary is the minimal structural description extracted from a file.
type Summary struct {
	Format    string
	Separator rune
	HasHeader bool
	Columns   []string
	Rows      int
}
