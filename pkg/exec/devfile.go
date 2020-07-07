package exec

import (
	"fmt"
	"io"
	"strings"
	"sync"

	adaptersCommon "github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/machineoutput"
	"github.com/pkg/errors"
)

// ExecuteDevfileBuildAction executes the devfile build command action
func ExecuteDevfileBuildAction(client ExecClient, exec common.Exec, commandName string, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) error {
	var s *log.Status

	// Change to the workdir and execute the command
	var cmdArr []string
	if exec.WorkingDir != "" {
		// since we are using /bin/sh -c, the command needs to be within a single double quote instance, for example "cd /tmp && pwd"
		cmdArr = []string{adaptersCommon.ShellExecutable, "-c", "cd " + exec.WorkingDir + " && " + exec.CommandLine}
	} else {
		cmdArr = []string{adaptersCommon.ShellExecutable, "-c", exec.CommandLine}
	}

	if show {
		s = log.SpinnerNoSpin("Executing " + commandName + " command " + fmt.Sprintf("%q", exec.CommandLine))
	} else {
		s = log.Spinnerf("Executing %s command %q", commandName, exec.CommandLine)
	}

	defer s.End(false)

	// Emit DevFileCommandExecutionBegin JSON event (if machine output logging is enabled)
	machineEventLogger.DevFileCommandExecutionBegin(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow())

	// Capture container text and log to the screen as JSON events (machine output only)
	stdoutWriter, stdoutChannel, stderrWriter, stderrChannel := machineEventLogger.CreateContainerOutputWriter()

	err := ExecuteCommand(client, compInfo, cmdArr, show, stdoutWriter, stderrWriter)

	// Close the writers and wait for an acknowledgement that the reader loop has exited (to ensure we get ALL container output)
	closeWriterAndWaitForAck(stdoutWriter, stdoutChannel, stderrWriter, stderrChannel)

	// Emit close event
	machineEventLogger.DevFileCommandExecutionComplete(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow(), err)
	if err != nil {
		return errors.Wrapf(err, "unable to execute the build command")
	}

	s.End(true)

	return nil
}

// ExecuteDevfileRunAction executes the devfile run command action using the supervisord devrun program
func ExecuteDevfileRunAction(client ExecClient, exec common.Exec, commandName string, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) error {
	var s *log.Status

	// Exec the supervisord ctl stop and start for the devrun program
	type devRunExecutable struct {
		command []string
	}
	devRunExecs := []devRunExecutable{
		{
			command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "stop", "all"},
		},
		{
			command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "start", string(adaptersCommon.DefaultDevfileRunCommand)},
		},
	}

	s = log.Spinnerf("Executing %s command %q", commandName, exec.CommandLine)
	defer s.End(false)

	for _, devRunExec := range devRunExecs {

		// Emit DevFileCommandExecutionBegin JSON event (if machine output logging is enabled)
		machineEventLogger.DevFileCommandExecutionBegin(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow())

		// Capture container text and log to the screen as JSON events (machine output only)
		stdoutWriter, stdoutChannel, stderrWriter, stderrChannel := machineEventLogger.CreateContainerOutputWriter()

		err := ExecuteCommand(client, compInfo, devRunExec.command, show, stdoutWriter, stderrWriter)

		// Close the writers and wait for an acknowledgement that the reader loop has exited (to ensure we get ALL container output)
		closeWriterAndWaitForAck(stdoutWriter, stdoutChannel, stderrWriter, stderrChannel)

		// Emit close event
		machineEventLogger.DevFileCommandExecutionComplete(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow(), err)
		if err != nil {
			return errors.Wrapf(err, "unable to execute the run command")
		}

	}
	s.End(true)

	return nil
}

// ExecuteDevfileRunActionWithoutRestart executes devfile run command without restarting.
func ExecuteDevfileRunActionWithoutRestart(client ExecClient, exec common.Exec, commandName string, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) error {
	var s *log.Status

	type devRunExecutable struct {
		command []string
	}
	// with restart false, executing only supervisord start command, if the command is already running, supvervisord will not restart it.
	// if the command is failed or not running suprvisord would start it.
	devRunExec := devRunExecutable{
		command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "start", string(adaptersCommon.DefaultDevfileRunCommand)},
	}

	s = log.Spinnerf("Executing %s command %q, if not running", commandName, exec.CommandLine)
	defer s.End(false)

	// Emit DevFileCommandExecutionBegin JSON event (if machine output logging is enabled)
	machineEventLogger.DevFileCommandExecutionBegin(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow())

	// Capture container text and log to the screen as JSON events (machine output only)
	stdoutWriter, stdoutChannel, stderrWriter, stderrChannel := machineEventLogger.CreateContainerOutputWriter()

	err := ExecuteCommand(client, compInfo, devRunExec.command, show, stdoutWriter, stderrWriter)

	// Close the writers and wait for an acknowledgement that the reader loop has exited (to ensure we get ALL container output)
	closeWriterAndWaitForAck(stdoutWriter, stdoutChannel, stderrWriter, stderrChannel)

	// Emit close event
	machineEventLogger.DevFileCommandExecutionComplete(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow(), err)
	if err != nil {
		return errors.Wrapf(err, "unable to execute the run command")
	}

	s.End(true)

	return nil
}

// ExecuteDevfileDebugAction executes the devfile debug command action using the supervisord debugrun program
func ExecuteDevfileDebugAction(client ExecClient, exec common.Exec, commandName string, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) error {
	var s *log.Status

	// Exec the supervisord ctl stop and start for the debugRun program
	type debugRunExecutable struct {
		command []string
	}
	debugRunExecs := []debugRunExecutable{
		{
			command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "stop", "all"},
		},
		{
			command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "start", string(adaptersCommon.DefaultDevfileDebugCommand)},
		},
	}

	s = log.Spinnerf("Executing %s command %q", commandName, exec.CommandLine)
	defer s.End(false)

	for _, debugRunExec := range debugRunExecs {

		// Emit DevFileCommandExecutionBegin JSON event (if machine output logging is enabled)
		machineEventLogger.DevFileCommandExecutionBegin(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow())

		// Capture container text and log to the screen as JSON events (machine output only)
		stdoutWriter, stdoutChannel, stderrWriter, stderrChannel := machineEventLogger.CreateContainerOutputWriter()

		err := ExecuteCommand(client, compInfo, debugRunExec.command, show, stdoutWriter, stderrWriter)

		// Close the writers and wait for an acknowledgement that the reader loop has exited (to ensure we get ALL container output)
		closeWriterAndWaitForAck(stdoutWriter, stdoutChannel, stderrWriter, stderrChannel)

		// Emit close event
		machineEventLogger.DevFileCommandExecutionComplete(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow(), err)

		if err != nil {
			return errors.Wrapf(err, "unable to execute the run command")
		}
	}

	s.End(true)

	return nil
}

// ExecuteDevfileDebugActionWithoutRestart executes devfile run command without restarting.
func ExecuteDevfileDebugActionWithoutRestart(client ExecClient, exec common.Exec, commandName string, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) error {
	var s *log.Status

	type devDebugExecutable struct {
		command []string
	}
	// with restart false, executing only supervisord start command, if the command is already running, supvervisord will not restart it.
	// if the command is failed or not running suprvisord would start it.
	devDebugExec := devDebugExecutable{
		command: []string{adaptersCommon.SupervisordBinaryPath, adaptersCommon.SupervisordCtlSubCommand, "start", string(adaptersCommon.DefaultDevfileDebugCommand)},
	}

	// Emit DevFileCommandExecutionBegin JSON event (if machine output logging is enabled)
	machineEventLogger.DevFileCommandExecutionBegin(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow())

	// Capture container text and log to the screen as JSON events (machine output only)
	stdoutWriter, stdoutChannel, stderrWriter, stderrChannel := machineEventLogger.CreateContainerOutputWriter()

	s = log.Spinnerf("Executing %s command %q, if not running", commandName, exec.CommandLine)
	defer s.End(false)

	err := ExecuteCommand(client, compInfo, devDebugExec.command, show, stdoutWriter, stderrWriter)

	// Close the writers and wait for an acknowledgement that the reader loop has exited (to ensure we get ALL container output)
	closeWriterAndWaitForAck(stdoutWriter, stdoutChannel, stderrWriter, stderrChannel)

	// Emit close event
	machineEventLogger.DevFileCommandExecutionComplete(exec.Id, exec.Component, exec.CommandLine, convertGroupKindToString(exec), machineoutput.TimestampNow(), err)

	if err != nil {
		return errors.Wrapf(err, "unable to execute the run command")
	}

	s.End(true)

	return nil
}

// ExecuteCompositeDevfileAction executes a given composite command in a devfile
// The composite command may reference exec commands, composite commands, or both
func ExecuteCompositeDevfileAction(client ExecClient, composite common.Composite, commandsMap map[string]common.DevfileCommand, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) (err error) {
	if composite.Parallel {
		fmt.Println("ToDo: Parallel")
		var wg sync.WaitGroup
		errChan := make(chan error)

		// Loop over each command and execute it in parallel
		for _, command := range composite.Commands {
			if devfileCommand, ok := commandsMap[strings.ToLower(command)]; ok {
				wg.Add(1)
				go func(client ExecClient, command common.DevfileCommand, commandsMap map[string]common.DevfileCommand, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) {
					defer wg.Done()
					err := execCommandFromComposite(client, devfileCommand, commandsMap, compInfo, show, machineEventLogger)
					if err != nil {
						errChan <- err
					}
				}(client, devfileCommand, commandsMap, compInfo, show, machineEventLogger)
			} else {
				return fmt.Errorf("command %q not found in devfile", command)
			}
		}

		// Wait for all parallel commands to finish
		wg.Wait()
		close(errChan)

		// Check the error channel, if any commands exited with an error
		err = <-errChan
		if err != nil {
			return fmt.Errorf("command execution failed: %v", err)
		}
	} else {
		fmt.Println("Non-parallel")

		for _, command := range composite.Commands {
			if devfileCommand, ok := commandsMap[strings.ToLower(command)]; ok {
				err = execCommandFromComposite(client, devfileCommand, commandsMap, compInfo, show, machineEventLogger)
				if err != nil {
					return fmt.Errorf("command execution failed: %v", err)
				}
			} else {
				// Devfile validation should have caught a missing command earlier, but should include error handling here as well
				return fmt.Errorf("command %q not found in devfile", command)
			}
		}
	}

	return nil
}

// execCommandFromComposite takes a command in a composite command and executes it.
// Any non-long running command (init, build, run, or debug) are treated the same
// Long-running run/debug commands, or run/debug commands that don't restart, need special handling
func execCommandFromComposite(client ExecClient, command common.DevfileCommand, commandsMap map[string]common.DevfileCommand, compInfo adaptersCommon.ComponentInfo, show bool, machineEventLogger machineoutput.MachineEventLoggingClient) (err error) {
	if command.Composite != nil {
		err = ExecuteCompositeDevfileAction(client, *command.Composite, commandsMap, compInfo, show, machineEventLogger)
	} else {
		switch command.Exec.Group.Kind {
		case common.InitCommandGroupType:
		case common.BuildCommandGroupType:
			err = ExecuteDevfileBuildAction(client, *command.Exec, command.Exec.Id, compInfo, show, machineEventLogger)
		case common.RunCommandGroupType:
			// Run commands are special in composite commands
			// Because of the current supervisord integration in odo, only one run command can be long-running (denoted by attribute ...)
			// Otherwise, we treat the run command like an ordinary build command
			err = ExecuteDevfileBuildAction(client, *command.Exec, command.Exec.Id, compInfo, show, machineEventLogger)
		case common.DebugCommandGroupType:
			// Like run commands, debug commands in composite commands are also special
			// Because of the current supervisord integration in odo, only one debug command can be long-running (denoted by attribute ...)
			// Otherwise, we treat the debug command like an ordinary build command
			err = ExecuteDevfileBuildAction(client, *command.Exec, command.Exec.Id, compInfo, show, machineEventLogger)
		}
	}

	return
}

// closeWriterAndWaitForAck closes the PipeWriter and then waits for a channel response from the ContainerOutputWriter (indicating that the reader had closed).
// This ensures that we always get the full stderr/stdout output from the container process BEFORE we output the devfileCommandExecution event.
func closeWriterAndWaitForAck(stdoutWriter *io.PipeWriter, stdoutChannel chan interface{}, stderrWriter *io.PipeWriter, stderrChannel chan interface{}) {
	if stdoutWriter != nil {
		_ = stdoutWriter.Close()
		<-stdoutChannel
	}
	if stderrWriter != nil {
		_ = stderrWriter.Close()
		<-stderrChannel
	}
}

func convertGroupKindToString(exec common.Exec) string {
	if exec.Group == nil {
		return ""
	}
	return string(exec.Group.Kind)
}
