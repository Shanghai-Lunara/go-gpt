package logic

type GitResponse struct {
	Name         string   `json:"name"`
	ActiveBranch string   `json:"active_branch"`
	ListBranches []string `json:"list_branches"`
	TaskCount    int32    `json:"task_count"`
	CurrentTask  string   `json:"current_task"`
}

type HttpRequest struct {
}

type HttpResponse struct {
}

func (s *Service) Request() {

}

func (s *Service) Response() {

}
