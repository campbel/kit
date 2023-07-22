package main

import (
	"reflect"
	"testing"
)

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
