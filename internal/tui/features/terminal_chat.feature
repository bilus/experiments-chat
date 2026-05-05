Feature: Minimal terminal chat prototype
  Scenario: Show the quit shortcut at startup
    Given the terminal chat application is running
    Then the terminal chat should show "Press Ctrl+D to quit."

  # REQ(#1)
  Scenario: Respond to any user input with the PoC response
    Given the terminal chat application is running
    When the user types "hello"
    And the user submits the prompt
    Then the terminal chat should render "You: hello"
    And the terminal chat should render "I don't understand."

  Scenario: Quit when the user presses Ctrl+D
    Given the terminal chat application is running
    When the user presses Ctrl+D
    Then the terminal chat should exit

  # REQ(#25)
  Scenario: Edit the pending prompt before submission
    Given the terminal chat application is running
    When the user types "heXlo"
    And the user presses Left Arrow
    And the user presses Left Arrow
    And the user presses Backspace
    And the user types "l"
    And the user submits the prompt
    Then the terminal chat should render "You: hello"
    And the terminal chat should render "I don't understand."
