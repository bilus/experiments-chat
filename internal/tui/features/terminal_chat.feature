Feature: Minimal terminal chat prototype
  Scenario: Respond to any user input with the PoC response
    Given the terminal chat application is running
    When the user enters "hello"
    Then the terminal chat should render "I don't understand."
