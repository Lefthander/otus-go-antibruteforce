Feature: AntiBruteForce Service
   As gRPC client of antibruteforce service
   In order to understand that user is blocked or not
   I want receive event from antibruteforce service

    Scenario: Verify limiter for login
        Given login "login1"
        And password "random"
        And ipaddress "random"
        And delay between request is "0.6s"
        When send 10 requests
        Then all requests are not blocked

    Scenario: Verify limiter for password
        Given login "random"
        And password "password1"
        And ipaddress "random"
        And delay between request is "0.06s"
        When send 100 requests
        Then All requests are not blocked
    
    Scenario: Verify limiter for ipaddress
        Given login "random"
        And password "random"
        And ipaddress "17.1.2.3"
        And delay between request is "0.006s"
        When send 1000 requests
        Then All requests are not blocked

    Scenario: Verify limiter for login with reset
        Given login "login1"
        And password "random"
        And ipaddress "17.1.2.3"
        And delay between request is "0.3s"
        And reset at 4 request
        When send 20 requests
        Then 15 request are passed