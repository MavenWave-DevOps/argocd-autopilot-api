package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ArgocdAutopilot struct {
	GitRepo      string   `json:"git-repo"`
	GitTokenPath string   `json:"git-token-path"`
	RootCommand  string   `json:"root-command"`
	Args         []string `json:"args"`
}

type Log struct {
	Message string `json:"logMessage"`
}

type ArgoCommand bool
type TokenCommand bool

func (r *ArgoCommand) StreamOutput(c *gin.Context, cChan chan Log) (error, string) {
	c.Stream(func(w io.Writer) bool {
		output, ok := <-cChan
		if !ok {
			return false
		}
		jsonOut, err := json.Marshal(output)
		if err != nil {
			log.Println("Error marshalling json: ", err)
		}
		outputBytes := bytes.NewBuffer([]byte("\n"))
		c.Writer.Write(append(outputBytes.Bytes(), jsonOut...))
		return true
	})
	return nil, ""
}

func (r *TokenCommand) StreamOutput(c *gin.Context, cChan chan Log) (error, string) {
	finalOutput := <-cChan
	return nil, finalOutput.Message
}

func CommandHelper(c *gin.Context, rootcmd string, args ...string) (string, error) {
	cmd := exec.Command(rootcmd, args...)
	cChan := make(chan Log)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		log.Println("Failed to start", err)
		return "", err
	}

	go func(o io.ReadCloser) {
		scanner := bufio.NewScanner(o)
		for scanner.Scan() {
			result := scanner.Text()
			cChan <- Log{
				Message: result,
			}
		}
		if rootcmd != "cat" {
			cChan <- Log{
				Message: fmt.Sprintf("%s %s: Command successful", rootcmd, strings.Trim(fmt.Sprint(args), "[]")),
			}
		}
		close(cChan)
	}(stdout)

	go func(o io.ReadCloser) {
		scanner := bufio.NewScanner(o)
		for scanner.Scan() {
			result := scanner.Text()
			log.Println(result)
		}
	}(stderr)
	var data string
	if rootcmd == "cat" {
		var cmdType TokenCommand = true
		err, data = StreamByCommand(&cmdType, c, cChan)
		if err != nil {
			return "", err
		}
	} else {
		var cmdType ArgoCommand = true
		err, data = StreamByCommand(&cmdType, c, cChan)
		if err != nil {
			return "", err
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("%s", cmd.String())
		log.Println(err)
		return "", err
	}
	return data, nil
}

func (r *ArgocdAutopilot) RunCommand(c *gin.Context) error {
	if err := os.Setenv("GIT_REPO", r.GitRepo); err != nil {
		log.Println(err)
		return err
	}
	//Set git token
	resp, err := CommandHelper(c, "cat", r.GitTokenPath)
	if err != nil {
		return err
	}
	if err := os.Setenv("GIT_TOKEN", fmt.Sprintf("%s", resp)); err != nil {
		log.Println(err)
		return err
	}
	log.Println("Git token set")

	//Run ArgoCD Command
	_, err = CommandHelper(c, r.RootCommand, r.Args...)
	if err != nil {
		return err
	}
	return nil
}

func ExecuteCommands(c *gin.Context) {
	var newCommand ArgocdAutopilot
	var myResponse = Response{}
	var wg = sync.WaitGroup{}

	wg.Add(1)

	// Call BindJSON to bind the received JSON to newCommand
	if err := c.BindJSON(&newCommand); err != nil {
		log.Printf("error %s", err)
		myResponse = Response{
			Status:  500,
			Message: "Internal Server Error",
		}
		SendResponse(c, myResponse)
		return
	}

	//Run commands
	go func() {
		if err := newCommand.RunCommand(c); err != nil {
			log.Printf("error %s", err)
			myResponse = Response{
				Status:  500,
				Message: "Internal Server Error",
			}
			SendResponse(c, myResponse)
			wg.Done()
			return
		}
		wg.Done()
	}()

	myResponse = Response{
		Status:  201,
		Message: "API Called Successfully",
	}
	SendResponse(c, myResponse)
	wg.Wait()
	return
}
