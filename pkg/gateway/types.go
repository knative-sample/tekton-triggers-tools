package gateway

type CommonResponse struct {
}

type AddPipelineRunArgs struct {
	GitRepo          string
	GitBranch        string
	GitUserName      string
	GitPassword      string
	NewImageFullName string
	NewImageTag      string
	CRUser           string
	CRPassword       string

	Name      string
	Namespace string
}
