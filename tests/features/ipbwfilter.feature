Feature: Black / White IP Filter

   Verifies the IP Filter with Black & White Lists

    Scenario: Verify Non defined IP address
        Given ipaddress "192.168.0.1"
        When verify ipaddress 
        Then request is not blocked

    Scenario: Add IP address to Black List
        Given ipaddress "10.0.0.0/24"
        When add ipaddress to black list

    Scenario: Verify IP address from Black List
        Given ipaddress "10.0.0.1"
        When verify ipaddress
        Then request is blocked

    Scenario: Add another IP address to Black List
        Given ipaddress "15.0.0.0/24"
        When add ipaddress to black list

    Scenario: Add IP address to Black List with error
        Given ipaddress "15.0.0.0/24"
        When add ipaddress to black list
        Then error reported - network already exists

    Scenario: Add IP address to White List
        Given ipaddress "20.0.0.0/24"
        When add ipaddress to white list

    Scenario: Verify IP address from White List
        Given ipaddress "20.0.0.1"
        When verify ipaddress
        Then request is not blocked

    Scenario: Add another IP address to White List
        Given ipaddress "25.0.0.0/24"
        When add ipaddress to white list

    Scenario: Add IP address to White List with error
        Given ipaddress "25.0.0.0/24"
        When add ipaddress to white list
        Then error reported - network already exists

# To check the priority of White list against Black 
    Scenario: Add Blacklisted IP Address to White List
        Given ipaddress "15.0.0.0/24"
        When add ipaddress to white list
        Then request is not blocked

    Scenario: Get the White List Contents
        When get white list contents
        Then received ipaddresses "20.0.0.0/24" And "25.0.0.0/24" And "15.0.0.0/24"

    Scenario: Get the Black List Contents
        When get black list contents
        Then received ipaddresses "10.0.0.0/24" And "15.0.0.0/24"

    Scenario: Delete unknown IP address from White List with error
        Given ipaddress "172.0.0.0/29"
        When delete ipaddress from white list
        Then received error - ipaddress not found
    
    Scenario: Delete IP address from White List
        Given ipaddress "25.0.0.0/24"
        When delete ipaddress from white list
        Then received status Ok

    Scenario: Delete black listed IP address from White List
        Given ipaddress "15.0.0.0/24"
        When delete ipaddress from white list
        Then received status Ok
    
    Scenario: Delete another IP address from White List
        Given ipaddress "20.0.0.0/24"
        When delete ipaddress from white list
        Then received status Ok

    Scenario: Get White list Contents
        When get white list contents
        Then received empty list

    Scenario: Delete unkown IP address from Black List with error
        Given ipaddress "172.0.0.0/29"
        When delete ipaddress from black list
        Then received error - ipaddress not found

    Scenario: Delete IP address from Black List
        Given ipaddress "15.0.0.0/24"
        When delete ipaddress from black list
        Then received status Ok

    Scenario: Delete another IP address from Black List
        Given ipaddress "10.0.0.0/24"
        When delete ipaddress from black list
        Then received status Ok

    Scenario: Get Black List Contents
        When get black list contents
        Then received empty list
