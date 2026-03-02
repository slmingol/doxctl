/*
Package cmdhelp - ...

Copyright © 2021 Sam Mingolelli <github@lamolabs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmdhelp

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPipeline_SingleCommand(t *testing.T) {
	cmd := exec.Command("echo", "hello world")
	output, stderr, err := Pipeline(cmd)

	if err != nil {
		t.Fatalf("Pipeline failed with error: %v", err)
	}

	if len(stderr) > 0 {
		t.Errorf("Expected no stderr, got: %s", string(stderr))
	}

	result := strings.TrimSpace(string(output))
	if result != "hello world" {
		t.Errorf("Expected 'hello world', got: %s", result)
	}
}

func TestPipeline_MultipleCommands(t *testing.T) {
	cmd1 := exec.Command("echo", "hello\nworld\ntest")
	cmd2 := exec.Command("grep", "world")
	output, stderr, err := Pipeline(cmd1, cmd2)

	if err != nil {
		t.Fatalf("Pipeline failed with error: %v", err)
	}

	if len(stderr) > 0 {
		t.Errorf("Expected no stderr, got: %s", string(stderr))
	}

	result := strings.TrimSpace(string(output))
	if result != "world" {
		t.Errorf("Expected 'world', got: %s", result)
	}
}

func TestPipeline_ThreeCommands(t *testing.T) {
	cmd1 := exec.Command("echo", "one\ntwo\nthree\nfour")
	cmd2 := exec.Command("grep", "-v", "two")
	cmd3 := exec.Command("grep", "three")
	output, stderr, err := Pipeline(cmd1, cmd2, cmd3)

	if err != nil {
		t.Fatalf("Pipeline failed with error: %v", err)
	}

	if len(stderr) > 0 {
		t.Errorf("Expected no stderr, got: %s", string(stderr))
	}

	result := strings.TrimSpace(string(output))
	if result != "three" {
		t.Errorf("Expected 'three', got: %s", result)
	}
}

func TestPipeline_NoCommands(t *testing.T) {
	output, stderr, err := Pipeline()

	if err != nil {
		t.Errorf("Expected no error for empty command list, got: %v", err)
	}

	if output != nil {
		t.Errorf("Expected nil output, got: %v", output)
	}

	if stderr != nil {
		t.Errorf("Expected nil stderr, got: %v", stderr)
	}
}

func TestPipeline_CommandWithError(t *testing.T) {
	cmd := exec.Command("ls", "/nonexistent/path/that/does/not/exist")
	_, _, err := Pipeline(cmd)

	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
}

func TestPipeline_PipelineWithGrepNoMatch(t *testing.T) {
	cmd1 := exec.Command("echo", "hello world")
	cmd2 := exec.Command("grep", "nomatch")
	_, _, err := Pipeline(cmd1, cmd2)

	// grep returns exit code 1 when no match is found
	if err == nil {
		t.Error("Expected error when grep finds no match, got nil")
	}
}
