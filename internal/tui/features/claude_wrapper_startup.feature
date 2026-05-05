Feature: Claude ACP wrapper startup
  Scenario: Show readiness after connecting to Claude
    Given Claude is the selected provider
    When the terminal chat application starts
    And connection to Claude is successful
    Then the system should indicate readiness
