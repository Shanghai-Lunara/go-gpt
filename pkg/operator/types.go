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

// Commands
type Command struct {
	ProjectName string `json:"project_name"`
	BranchName  string `json:"branch_name"`
	Command     string `json:"command"`
	Message     string `json:"message"`
	ZipType     string `json:"zip_type"`
	ZipFlags    string `json:"zip_flags"`
}

// Tasks
