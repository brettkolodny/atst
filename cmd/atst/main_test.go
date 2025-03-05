package main

import (
	"reflect"
	"testing"

	atst "github.com/brettkolodny/atst/pkg/atst"
)

func TestParseCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []atst.Program
	}{
		{
			name: "Basic case - one program with arguments",
			args: []string{"-p", "program1", "-a", "arg1", "-a", "arg2"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{"arg1", "arg2"},
				},
			},
		},
		{
			name: "Multiple programs with arguments",
			args: []string{"-p", "program1", "-a", "arg1", "-a", "arg2", "-p", "program2", "-a", "arg3", "-a", "arg4"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{"arg1", "arg2"},
				},
				{
					Exec: "program2",
					Args: []string{"arg3", "arg4"},
				},
			},
		},
		{
			name: "atst.Program with no arguments",
			args: []string{"-p", "program1"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{},
				},
			},
		},
		{
			name: "Arguments starting with dashes",
			args: []string{"-p", "countdown", "-a", "--countdown", "-a", "10"},
			expected: []atst.Program{
				{
					Exec: "countdown",
					Args: []string{"--countdown", "10"},
				},
			},
		},
		{
			name: "Arguments before program (should be ignored)",
			args: []string{"-a", "ignored", "-p", "program1", "-a", "arg1"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{"arg1"},
				},
			},
		},
		{
			name: "Missing argument value (should be skipped)",
			args: []string{"-p", "program1", "-a"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{},
				},
			},
		},
		{
			name:     "Missing program value (should be skipped)",
			args:     []string{"-p"},
			expected: []atst.Program{},
		},
		{
			name: "Using long flag names",
			args: []string{"--program", "program1", "--arg", "arg1"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{"arg1"},
				},
			},
		},
		{
			name: "Mixed short and long flag names",
			args: []string{"-p", "program1", "--arg", "arg1", "-a", "arg2"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{"arg1", "arg2"},
				},
			},
		},
		{
			name: "Arguments for previous program after defining new one",
			args: []string{"-p", "program1", "-p", "program2", "-a", "arg1"},
			expected: []atst.Program{
				{
					Exec: "program1",
					Args: []string{},
				},
				{
					Exec: "program2",
					Args: []string{"arg1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePrograms(tt.args)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCommands() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsFlag(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"-p", true},
		{"--program", true},
		{"-", true},
		{"p", false},
		{"program", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isFlag(tt.input)
			if result != tt.expected {
				t.Errorf("isFlag(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
