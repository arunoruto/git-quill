package ai

type Request struct {
	Diff        string
	StagedFiles string
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
