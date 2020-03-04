package logic

type ProjectOperator interface {
	GitGenerate()
}

type Project struct {
	Git *Git
	Svn *SvnOperator
}

type ProjectHub struct {
	Projects []Project
}
