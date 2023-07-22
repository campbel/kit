package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetTaskFile(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		fileExistsFn func(path string) bool
		shouldExist  bool
		want         string
	}{
		{
			name: "returns taskfile from flag",
			args: []string{"cmd", "-t", "Taskfile.yml"},
			fileExistsFn: func(path string) bool {
				return path == "Taskfile.yml"
			},
			shouldExist: true,
			want:        "Taskfile.yml",
		},
		{
			name: "returns taskfile from equals flag",
			args: []string{"cmd", "--taskfile=Taskfile.yml"},
			fileExistsFn: func(path string) bool {
				return path == "Taskfile.yml"
			},
			shouldExist: true,
			want:        "Taskfile.yml",
		},
		{
			name: "returns default taskfile",
			args: []string{"cmd"},
			fileExistsFn: func(path string) bool {
				return path == "Taskfile.yml"
			},
			shouldExist: true,
			want:        "Taskfile.yml",
		},
		{
			name: "returns empty string if no taskfile found",
			args: []string{"cmd"},
			fileExistsFn: func(path string) bool {
				return false
			},
			shouldExist: false,
			want:        "",
		},
		{
			name: "returns another default if that exists",
			args: []string{"cmd"},
			fileExistsFn: func(path string) bool {
				return path == "taskfile.yml"
			},
			shouldExist: true,
			want:        "taskfile.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			got, exists := getTaskFile(tt.fileExistsFn)
			if exists != tt.shouldExist {
				t.Errorf("getTaskFile() exists = %v, want %v", exists, tt.want != "")
			}
			if got != tt.want {
				t.Errorf("getTaskFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "no flags",
			args: []string{"cmd", "arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
		{
			name: "taskfile flag",
			args: []string{"cmd", "--taskfile", "Taskfile.yml", "arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
		{
			name: "taskfile flag with equals sign",
			args: []string{"cmd", "-t=Taskfile.yml", "arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
		{
			name: "missing arg after taskfile flag",
			args: []string{"cmd", "--taskfile", "Taskfile.yml"},
			want: []string{},
		},
		{
			name: "multiple flags",
			args: []string{"cmd", "--taskfile", "Taskfile.yml", "-t=Taskfile2.yml", "arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
		{
			name: "multiple flags with missing arg",
			args: []string{"cmd", "--taskfile", "Taskfile.yml", "-t=Taskfile2.yml"},
			want: []string{},
		},
		{
			name: "multiple flags with missing param",
			args: []string{"cmd", "--taskfile"},
			want: []string{"--taskfile"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterArgs(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
