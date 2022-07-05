package steps

import "github.com/cucumber/godog"

// AnswerSteps Struct
type AnswerSteps struct {
}

func NewSteps(s *godog.ScenarioContext) *AnswerSteps {
	steps := &AnswerSteps{}

	// Given
	//suite.Step(`^the authorization "([^"]*)"$`, steps.givenTheAuthorization)

	// When
	//suite.Step(`^i send authorization$`, steps.whenISendAuthorization)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	s.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, theResponseShouldMatchJson)

	// Then
	//suite.Step(`^the authorization status is "([^"]*)"$`, steps.thenTheAuthorizationResultStatusIs)

	return steps
}

func iSendRequestTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
}

func theResponseShouldMatchJson(arg1 *godog.DocString) error {
	return godog.ErrPending
}
