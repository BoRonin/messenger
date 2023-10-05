package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"messanger/internal/entity"
	"messanger/internal/service"
	"net/http"
	"sync"
)

type MClient struct {
	URL      string
	Token    string
	Mservice *service.MessageService
}

func NewMClient(u string, t string, ms *service.MessageService) *MClient {
	return &MClient{
		URL:      u,
		Token:    t,
		Mservice: ms,
	}
}

func (MC *MClient) SendMessages(ctx context.Context, messages []entity.Message) error {
	mchan := make(chan uint)
	wg := &sync.WaitGroup{}
	for _, v := range messages {
		wg.Add(1)
		go func(c entity.Message) {
			client := &http.Client{}
			message2go := entity.MessageToGo{
				ID:    c.ID,
				Phone: int(c.Client.Tel),
				Text:  c.Messanger.MessageText,
			}
			log.Println("sending now: ", message2go)

			url := fmt.Sprintf("%s%d", MC.URL, message2go.ID)
			println(url)
			auth := fmt.Sprintf("Bearer %s", MC.Token)
			jsonMes, _ := json.Marshal(message2go)

			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonMes))
			req.Header.Add("Authorization", auth)
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			response := struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{}
			json.NewDecoder(resp.Body).Decode(&response)
			defer resp.Body.Close()
			if err := MC.Mservice.DB.SetStatusToTrue(context.Background(), message2go.ID); err != nil {
				println(err)
			}
			if response.Message == "OK" {
				mchan <- message2go.ID
			}
			defer wg.Done()
		}(v)
	}
	go func() {
		wg.Wait()
		log.Println("closing uint message chan")
		close(mchan)
	}()
	for {
		i, ok := <-mchan
		if !ok {
			break
		}
		fmt.Printf("Message with id %d was sent\n", i)
	}
	println("finished sending messages")
	return nil
}
