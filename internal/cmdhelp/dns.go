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
	"bytes"
	"os/exec"
)

// Pipeline chain exec.Commands together in a piped seq. of commands
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	// Use separate stderr buffers for each command to avoid race conditions
	stderrBuffers := make([]*bytes.Buffer, len(cmds))
	for i := range stderrBuffers {
		stderrBuffers[i] = &bytes.Buffer{}
	}

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to its own buffer
		cmd.Stderr = stderrBuffers[i]
	}

	// Connect the output and error for the last command
	cmds[last].Stdout = &output
	cmds[last].Stderr = stderrBuffers[last]

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			// Wait for any commands that were already started
			for j, c := range cmds {
				if c.Process != nil {
					c.Wait() // Ignore error, we're already handling one
				}
				if j >= len(cmds) {
					break
				}
			}
			// Combine all stderr buffers collected so far
			var combinedStderr bytes.Buffer
			for _, buf := range stderrBuffers {
				combinedStderr.Write(buf.Bytes())
			}
			return output.Bytes(), combinedStderr.Bytes(), err
		}
	}

	// Wait for each command to complete - collect all errors but ensure all complete
	var firstErr error
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	// Now all commands have completed and all I/O is done
	// Combine all stderr buffers
	var combinedStderr bytes.Buffer
	for _, buf := range stderrBuffers {
		combinedStderr.Write(buf.Bytes())
	}

	// Return the pipeline output and the collected standard error
	if firstErr != nil {
		return output.Bytes(), combinedStderr.Bytes(), firstErr
	}
	return output.Bytes(), combinedStderr.Bytes(), nil
}
