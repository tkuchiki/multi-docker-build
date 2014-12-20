package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"
)

func DockerBuild(dockerfile, tagOption, options string, quiet bool) (exitCode int, err error) {
	var cmd *exec.Cmd

	if len(options) == 0 {
		fmt.Println(fmt.Sprintf(`Start "docker build %s %s"`, tagOption, dockerfile))
		fmt.Println()
		cmd = exec.Command("docker", "build", "-t", tagOption, ".")
	} else {
		fmt.Println(fmt.Sprintf(`Start "docker build %s %s %s"`, tagOption, options, dockerfile))
		fmt.Println()
		cmd = exec.Command("docker", "build", "-t", tagOption, options, ".")
	}

	exitCode, err = RunCommand(cmd, quiet)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	return
}

func RunCommand(cmd *exec.Cmd, quiet bool) (exitCode int, err error) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	go PrintOutput(outReader, quiet)
	go PrintOutput(errReader, quiet)

	err = cmd.Wait()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}

	return
}

func PrintOutput(r io.Reader, quiet bool) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if !quiet {
			fmt.Println(scanner.Text())
		}
	}
}
