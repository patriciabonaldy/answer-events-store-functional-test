# file: online-create-answer.feature
Feature: create an answer
  In order to create a new answer
  As an API user
  I need to be able to create a new answer

  Scenario: error POST method
    Given the request "empty"
    When I send "POST" request to "/answers"
    Then the response code should be 400
    And the response should match:
      """
      Key: 'CreateRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag
      """

  Scenario: success POST method
    Given the request "default"
    When I send "POST" request to "/answers"
    Then the response code should be 201
    Then I find the answer
    Then the response code should be 200
