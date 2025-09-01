package bringauto_ssh

import (
	"io"
	"strings"
)

// ShellEvaluator run commands by the bash thru the SSH.
//
type ShellEvaluator struct {
	// Commands to execute, each command must exit with 0 return code else error is returned.
	Commands []string
	// Preparing commands are executed before main commands. Their exit values are not captured and checked.
	PreparingCommands []string
	// Environments variables to set
	Env map[string]string
	// StdOut writer to capture stdout of the process
	StdOut io.Writer
}

// getEnvStr
// Returns string for setting environment variables from shell.Env.
// If shell.Env is empty, it returns empty string. If shell.Env is not empty, the returned string
// exports all variables. The string contains trailing " && " string for connecting with following
// commands.
func (shell *ShellEvaluator) getEnvStr() string {
	// We cannot use SSHSession/SSHConnection setenv function
	// for Env. setting because SetEnv must be configured at the Server side
	envStr := ""
	for envName, envValue := range shell.Env {
		envStr = envStr + "export " + envName + "=" + escapeVariableValue(envValue) + " && "
	}
	return envStr
}

// getCommandStr
// Returns string for running commands. The string contains trailing " && " string for connecting
// with following commands.
func (shell *ShellEvaluator) getCommandStr() string {
	commandStr := ""
	for _, value := range shell.Commands {
		commandStr = commandStr + value  + " && "
	}
	commandStr = commandStr + "echo done"

	return commandStr
}

// getPreparingCommandStr
// Returns string for running preparing commands.
func (shell *ShellEvaluator) getPreparingCommandStr() string {
	if len(shell.PreparingCommands) == 0 {
		return ""
	}
	prepCommandStr := "{ "
	for _, value := range shell.PreparingCommands {
		prepCommandStr = prepCommandStr + value  + "; "
	}
	prepCommandStr = prepCommandStr + "} || true && "
	return prepCommandStr
}

// RunOverSSH
// Runs command over SSH.
//
// All commands specified in Commands are run be Bash by a one Bash session
//
// All environment variables are preserved across command run and can be used by other
// subsequent commands.
func (shell *ShellEvaluator) RunOverSSH(credentials SSHCredentials) error {
	var err error
	pipeReader, _ := io.Pipe()
	session := SSHSession{
		StdOut: shell.StdOut,
		StdErr: shell.StdOut,
		StdIn:  pipeReader,
	}

	err = session.LoginMultipleAttempts(credentials)
	if err != nil {
		return err
	}

	raw := shell.getEnvStr() + shell.getPreparingCommandStr() + shell.getCommandStr()
	safe := strings.ReplaceAll(raw, "'", "'\\''") // Escaping all single quotes
	cmdStr := "bash -li -c '" + safe + "'"

	err = session.Run(cmdStr)
	if err != nil {
		return err
	}

	return nil
}

func escapeVariableValue(varValue string) string {
	return "\"" + varValue + "\""
}
