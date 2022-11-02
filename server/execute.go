package server

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
)

type ArgocdAutopilot struct {
	GitRepo      string   `json:"git-repo"`
	GitTokenPath string   `json:"git-token-path"`
	RootCommand  string   `json:"root-command"`
	Args         []string `json:"args"`
}

func CommandHelper(rootcmd string, args ...string) (string, error) {

	cmd := exec.Command(rootcmd, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}
	if err := cmd.Start(); err != nil {
		log.Println("Failed to start", err)
		return "", err
	}

	c := make(chan string)

	go func(o io.ReadCloser) {
		scanner := bufio.NewScanner(o)
		for scanner.Scan() {
			result := scanner.Text()
			log.Println(result)
			c <- result
		}
		close(c)
	}(stdout)

	go func(o io.ReadCloser) {
		scanner := bufio.NewScanner(o)
		for scanner.Scan() {
			result := scanner.Text()
			log.Println(result)
		}
	}(stderr)

	err = cmd.Wait()
	if err != nil {
		log.Printf("%s", cmd.String())
		log.Println(err)
		return "", err
	}
	return <-c, nil
}

func (r *ArgocdAutopilot) RunCommand() error {
	if err := os.Setenv("GIT_REPO", r.GitRepo); err != nil {
		log.Println(err)
		return err
	}

	//Set git token
	resp, err := CommandHelper("cat", r.GitTokenPath)
	if err != nil {
		return err
	}
	if err := os.Setenv("GIT_TOKEN", fmt.Sprintf("%s", resp)); err != nil {
		log.Println(err)
		return err
	}

	//Run ArgoCD Command
	output, err := CommandHelper(r.RootCommand, r.Args...)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func ExecuteCommands(c *gin.Context) {
	var newCommand ArgocdAutopilot

	// Call BindJSON to bind the received JSON to newCommand
	if err := c.BindJSON(&newCommand); err != nil {
		log.Printf("error %s", err)
		return
	}

	//Run commands
	if err := newCommand.RunCommand(); err != nil {
		log.Printf("error %s", err)
		myResp := Response{
			Status:  500,
			Message: "Internal Server Error",
		}
		SendResponse(c, myResp)
		return
	}

	myResp := Response{
		Status:  201,
		Message: "All good",
	}
	SendResponse(c, myResp)
	return
}
