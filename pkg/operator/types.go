package operator

type ProjectConfig struct {
	ProjectName string          `yaml:"project_name"`
	ScriptsPath string          `yaml:"scripts_path"`
	Git         GitConfig       `yaml:"git"`
	Svn         SvnConfig       `yaml:"svn"`
	Ftp         FtpConfig       `yaml:"ftp"`
	Oss         AliYunOssConfig `yaml:"oss"`
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

// Command
// swagger:response Command
type Command struct {
	ProjectName string `json:"project_name"`
	BranchName  string `json:"branch_name"`
	Command     string `json:"command"`
	Message     string `json:"message"`
	ZipType     string `json:"zip_type"`
	ZipFlags    string `json:"zip_flags"`
}

// AliYunOss types
type AliYunOssConfig struct {
	EndPoint        string         `yaml:"end_point"`
	Bucket          string         `yaml:"bucket"`
	AccessKeyID     string         `yaml:"access_key_id"`
	AccessKeySecret string         `yaml:"access_key_secret"`
	ProxyUrl        string         `yaml:"proxy_url"`
	FileDirectory   string         `yaml:"file_directory"`
	BucketName      string         `yaml:"bucket_name"`
	Envs            []AliYunOssEnv `yaml:"envs"`
}

type AliYunOssEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// NoticeContent
// swagger:response NoticeContent
type NoticeContent struct {
	Title   string `json:"title"`
	Time    string `json:"time"`
	Content string `json:"content"`
}
