package operator

type ProjectConfig struct {
	ProjectName string    `yaml:"project_name"`
	ScriptsPath string    `yaml:"scripts_path"`
	Git         GitConfig `yaml:"git"`
	Svn         SvnConfig `yaml:"svn"`
	Ftp         FtpConfig `yaml:"ftp"`
}

// git types
type GitConfig struct {
	WorkDir string `yaml:"work_dir"`
}

type GitInfo struct {
	Name         string   `json:"name"`
	ListBranches []Branch `json:"list_branches"`
	TaskCount    int32    `json:"task_count"`
	CurrentTask  string   `json:"current_task"`
}

type Branch struct {
	Name   string `json:"name"`
	Active int    `json:"active"`
	SvnTag string `json:"svn_tag"`
}

// svn types
type SvnConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	WorkDir   string `yaml:"work_dir"`
	Url       string `yaml:"url"`
	Port      int    `yaml:"port"`
	RemoteDir string `yaml:"remote_dir"`
}

// ftp types
type FtpConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	WorkDir string `yaml:"work_dir"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}
