package status

type ImportState string

const (
	StateNotImported        ImportState = "not_imported"
	StateImporting          ImportState = "importing"
	StateImported           ImportState = "imported"
	StateClonePending       ImportState = "clone_pending"
	StateCloneFailed        ImportState = "clone_failed"
	StateRevalidationFailed ImportState = "revalidation_failed"
	StateNeedsRefresh       ImportState = "needs_refresh"
	StateCloned             ImportState = "cloned"
)

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Resolve(current string, checks []CheckResult) ImportState {
	hasValidClone := false
	hasMissing := false

	for _, c := range checks {
		if c.Name == "clone_exists" && c.Passed {
			hasValidClone = true
		}
		if c.Name == "clone_exists" && !c.Passed {
			hasMissing = true
		}
		if c.Name == "git_directory" && !c.Passed {
			hasMissing = true
		}
	}

	switch current {
	case "pending", "":
		if hasMissing {
			return StateClonePending
		}
		return StateImporting
	case "cloned":
		if hasMissing {
			return StateRevalidationFailed
		}
		return StateCloned
	case "imported":
		if hasValidClone {
			return StateCloned
		}
		if hasMissing {
			return StateClonePending
		}
		return StateImported
	case "clone_failed":
		if hasValidClone {
			return StateCloned
		}
		return StateCloneFailed
	default:
		if hasMissing {
			return StateNeedsRefresh
		}
		return ImportState(current)
	}
}

type CheckResult struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Error  string `json:"error,omitempty"`
}

func (e *Engine) ValidTransitions() map[ImportState][]ImportState {
	return map[ImportState][]ImportState{
		StateNotImported:        {StateImporting},
		StateImporting:          {StateClonePending, StateCloneFailed, StateCloned},
		StateClonePending:       {StateCloned, StateCloneFailed, StateRevalidationFailed},
		StateCloneFailed:        {StateClonePending, StateCloned},
		StateCloned:             {StateImported, StateRevalidationFailed},
		StateImported:           {StateNeedsRefresh, StateRevalidationFailed},
		StateRevalidationFailed: {StateCloned, StateNeedsRefresh, StateClonePending},
		StateNeedsRefresh:       {StateImported, StateRevalidationFailed},
	}
}
