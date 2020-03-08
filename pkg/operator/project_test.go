package operator

import (
	"context"
	"reflect"
	"testing"
)

func Test_projects_Add(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
		p           *project
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			ph.Add(tt.args.projectName, tt.args.p)
		})
	}
}

func Test_projects_GetProject(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantP   *project
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			gotP, err := ph.GetProject(tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("projects.GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("projects.GetProject() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func Test_projects_GetGitInfo(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantGi  GitInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			gotGi, err := ph.GetGitInfo(tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("projects.GetGitInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGi, tt.wantGi) {
				t.Errorf("projects.GetGitInfo() = %v, want %v", gotGi, tt.wantGi)
			}
		})
	}
}

func Test_projects_GitGenerate(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
		branchName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			if err := ph.GitGenerate(tt.args.projectName, tt.args.branchName); (err != nil) != tt.wantErr {
				t.Errorf("projects.GitGenerate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projects_SetGitBranchSvnTag(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
		branchName  string
		svnTag      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			if err := ph.SetGitBranchSvnTag(tt.args.projectName, tt.args.branchName, tt.args.svnTag); (err != nil) != tt.wantErr {
				t.Errorf("projects.SetGitBranchSvnTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projects_SvnCommit(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
		branchName  string
		svnMessage  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			if err := ph.SvnCommit(tt.args.projectName, tt.args.branchName, tt.args.svnMessage); (err != nil) != tt.wantErr {
				t.Errorf("projects.SvnCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projects_SvnLog(t *testing.T) {
	type fields struct {
		projects map[string]*project
	}
	type args struct {
		projectName string
		number      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []Logentry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := &projects{
				projects: tt.fields.projects,
			}
			gotRes, err := ph.SvnLog(tt.args.projectName, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("projects.SvnLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("projects.SvnLog() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestNewProject(t *testing.T) {
	type args struct {
		conf []ProjectConfig
		ctx  context.Context
	}
	tests := []struct {
		name string
		args args
		want *Project
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProject(tt.args.conf, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProject() = %v, want %v", got, tt.want)
			}
		})
	}
}
