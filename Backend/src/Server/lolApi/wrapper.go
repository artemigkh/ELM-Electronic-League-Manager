package lolApi

import (
	"Server/dataModel"
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
		responseBody <- nil
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Non 200 Status Code: %v %v", res.StatusCode, res.Status))
		responseBody <- nil
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read body. Error:", err.Error())
		responseBody <- nil
		return
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

func (l *lolApi) CompletePlayerStubs(stub *dataModel.LoLTeamStub) (*dataModel.LoLTeamWithRosters, error) {
	team := dataModel.LoLTeamWithRosters{
		TeamId:           stub.TeamId,
		Name:             stub.Name,
		Description:      stub.Description,
		Tag:              stub.Tag,
		IconSmall:        stub.IconSmall,
		IconLarge:        stub.IconLarge,
		Wins:             stub.Wins,
		Losses:           stub.Losses,
		MainRoster:       make([]*dataModel.LoLPlayer, 0),
		SubstituteRoster: make([]*dataModel.LoLPlayer, 0),
	}

	allPlayerIds := make([]string, 0)
	for _, player := range append(stub.MainRoster, stub.SubstituteRoster...) {
		allPlayerIds = append(allPlayerIds, player.ExternalId)
	}

	summonerMap := l.GetSummonerInformation(allPlayerIds)
	for _, player := range stub.MainRoster {
		info := summonerMap[player.ExternalId]
		team.MainRoster = append(team.MainRoster, &dataModel.LoLPlayer{
			PlayerId:       player.PlayerId,
			GameIdentifier: info.GameIdentifier,
			MainRoster:     player.MainRoster,
			Position:       player.Position,
			Rank:           info.Rank,
			Tier:           info.Tier,
		})
	}
	for _, player := range stub.SubstituteRoster {
		info := summonerMap[player.ExternalId]
		team.SubstituteRoster = append(team.SubstituteRoster, &dataModel.LoLPlayer{
			PlayerId:       player.PlayerId,
			GameIdentifier: info.GameIdentifier,
			MainRoster:     player.MainRoster,
			Position:       player.Position,
			Rank:           info.Rank,
			Tier:           info.Tier,
		})
	}
	//TODO: proper error
	return &team, nil
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

func (l *lolApi) GetMatchStats(id string) (*dataModel.LoLMatchInformation, error) {
	bodyBytes, err := l.callWrapperServer(fmt.Sprintf("gameStats?id=%v", id))
	if err != nil {
		println(err.Error())
		return nil, err
	}
	var matchInfo dataModel.LoLMatchInformation
	if err := json.Unmarshal(bodyBytes, &matchInfo); err != nil {
		println(err.Error())
		return nil, err
	}

	return &matchInfo, nil
}
