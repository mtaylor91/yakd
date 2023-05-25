package executor

type ChrootExecutor struct {
	// contains filtered or unexported fields
	root string
}

// NewChrootExecutor returns a new ChrootExecutor.
func NewChrootExecutor(root string) *ChrootExecutor {
	return &ChrootExecutor{root: root}
}
