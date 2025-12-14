package ai

type Request struct {
	Task        string // "commit" or "tag"
	Diff        string // For commit: the diff. For tag: the commit log
	StagedFiles string // Only for commit
	Model       string
	IsBrief     bool
	UseEmoji    bool
}

type Provider interface {
	Name() string
	IsAvailable() bool
	ListModels() ([]string, error)
	Generate(req Request) (string, error)
}
