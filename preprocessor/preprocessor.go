package preprocessor

// Preprocessor defines the interface for preprocessing command arguments
// before they are passed to translators. Preprocessors can extract hostnames
// from URLs, normalize inputs, or perform other transformations.
type Preprocessor interface {
	// ToolName returns the name of the tool this preprocessor applies to
	ToolName() string

	// Description returns a brief description of what this preprocessor does
	Description() string

	// Preprocess transforms the arguments before translation
	// Returns the preprocessed arguments
	Preprocess(args []string) []string
}
