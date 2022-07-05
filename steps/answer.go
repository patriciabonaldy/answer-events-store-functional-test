package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/patriciabonaldy/answer-events-store-functional/domain"
	"github.com/patriciabonaldy/answer-events-store-functional/platform/genericClient"
)

// AnswerSteps Struct
type AnswerSteps struct {
	client   genericClient.Client
	Request  domain.AnswerRequest
	Response domain.Response
	url      string
}

func NewSteps(s *godog.ScenarioContext, client genericClient.Client, routerURL string) *AnswerSteps {
	steps := &AnswerSteps{client: client, url: routerURL}

	// Given
	s.Step(`^the request "([^"]*)"$`, steps.givenRequest)

	// When
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, steps.iSendPostRequestTo)
	s.Step(`^I find the answer$`, steps.iFindTheAnswer)

	// Then
	s.Step(`^the response code should be (\d+)$`, steps.theResponseCodeShouldBe)
	s.Step(`^the response should match:$`, steps.theResponseShouldMatch)

	return steps
}

func (a *AnswerSteps) iSendPostRequestTo() error {
	body, err := json.Marshal(a.Request)
	if err != nil {
		return err
	}

	resp, err := a.client.Post(context.Background(), fmt.Sprintf("%s/answers", a.url), body)
	if err != nil {
		a.Response.Code = strconv.Itoa(resp.StatusCode)
		data, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return err
		}

		a.Response.Content = string(data)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	a.Response.ID = strings.ReplaceAll(string(data), "\"", "")
	a.Response.Code = strconv.Itoa(resp.StatusCode)
	a.Response.Content = string(data)

	return nil
}

func (a *AnswerSteps) iFindTheAnswer() error {
	uri := fmt.Sprintf("%s/answers/%s", a.url, a.Response.ID)
	fmt.Printf("uri: %s\n", uri)

	<-time.After(3 * time.Second)
	resp, err := a.client.Get(context.Background(), uri)
	if err != nil {
		a.Response.Code = strconv.Itoa(resp.StatusCode)
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &a.Response)
	if err != nil {
		return err
	}

	a.Response.Code = strconv.Itoa(resp.StatusCode)
	a.Response.Content = string(data)

	return nil
}

func (a *AnswerSteps) givenRequest(arg string) error {
	switch arg {
	case "default":
		a.Request = givenDefaultRequest()
	case "empty":
		a.Request = domain.AnswerRequest{}
	default:
		return errors.New("invalid type request")
	}

	return nil
}

func givenDefaultRequest() domain.AnswerRequest {
	return domain.AnswerRequest{
		Data: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
	}
}

func (a *AnswerSteps) theResponseCodeShouldBe(code int) error {
	if strconv.Itoa(code) != a.Response.Code {
		errors.New("invalid code") //nolint:errcheck
	}

	return nil
}

func (a *AnswerSteps) theResponseShouldMatch(resp *godog.DocString) error {
	if strings.EqualFold(a.Response.Content, resp.Content) {
		errors.New("invalid response")
	}

	return nil
}
