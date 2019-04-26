package lolApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func startWrapperServer() {
	// start python api wrapper server
	cmd := exec.Command("python", "Backend/src/Server/lolApi/api_wrapper_server.py")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// catch exit signal to kill process cleanly
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigchan
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill lol api wrapper server: ", err)
		}
	}()
}

func (l *lolApi) callWrapperServer(url string) ([]byte, error) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:8090/"+url),
		nil)
	res, err := l.client.Do(httpReq)
	if err != nil {
		fmt.Println("Failed to execute HTTP Request. Error:", err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Non 200 Status Code: %v %v", res.StatusCode, res.Status))
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read body. Error:", err.Error())
		return nil, err
	}

	return bodyBytes, nil
}

func (l *lolApi) callWrapperServerAsync(url string, responseBody chan []byte) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:8090/"+url),
		nil)
	res, err := l.client.Do(httpReq)
	if err != nil {
		fmt.Println("Failed to execute HTTP Request. Error:", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Non 200 Status Code: %v %v", res.StatusCode, res.Status))
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read body. Error:", err.Error())
	}

	responseBody <- bodyBytes
}

func (l *lolApi) GetSummonerInformation(ids []string) map[string]*SummonerInformation {
	waitChannels := make(map[string]chan []byte)
	summonerInformation := make(map[string]*SummonerInformation)

	for _, id := range ids {
		waitChannels[id] = make(chan []byte, 1)
		go l.callWrapperServerAsync(fmt.Sprintf("summonerInformation?id=%v", id), waitChannels[id])
	}

	for _, id := range ids {
		bodyBytes := <-waitChannels[id]
		si := &SummonerInformation{
			GameIdentifier: "",
			Rank:           "",
			Tier:           "",
		}
		if err := json.Unmarshal(bodyBytes, &si); err != nil {
			fmt.Println("Failed to unpack response. Error:", err.Error())
		}

		summonerInformation[id] = si
	}
	return summonerInformation
}

type idContainer struct {
	Id string `json:"id"`
}

func (l *lolApi) GetSummonerId(name string) (string, error) {
	bodyBytes, err := l.callWrapperServer(fmt.Sprintf("summonerId?name=%v", name))
	if err != nil {
		return "", err
	}
	var summonerId idContainer
	if err := json.Unmarshal(bodyBytes, &summonerId); err != nil {
		return "", err
	}

	println("got id: ", summonerId.Id)
	return summonerId.Id, nil
}

func (l *lolApi) GetMatchStats(id string) (*MatchInformation, error) {
	bodyBytes, err := l.callWrapperServer(fmt.Sprintf("gameStats?id=%v", id))
	if err != nil {
		println(err.Error())
		return nil, err
	}
	var matchInfo MatchInformation
	if err := json.Unmarshal(bodyBytes, &matchInfo); err != nil {
		println(err.Error())
		return nil, err
	}

	return &matchInfo, nil
}
