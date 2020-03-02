Feature: AntiBruteForce Service
   As gRPC client of antibruteforce service
   In order to understand that user is blocked or not
   I want receive event from antibruteforce service

Scenario: Check IP Address which is not in the B/W Lists 
    Given ipaddress "30.30.0.1"
    When check address
    Then request is not blocked

Scenario: Add Network to the Black List
    Given network "10.10.0.0/24"
    When add network to black List

Scenario: Add the same network to the Black list
    Given network "10.10.0.0/24"
    When add network to black list
    Then error network already exists

Scenario: Add Network to the White List
    Given network "20.20.0.0/24"
    When add network to white List

Scenario: Add the same Network to the White List
    Given network "20.20.0.0/24"
    When add network to white List
    Then error network already exists

Scenario: Check Address from White List
    Given ipaddress "20.20.0.1"
    When check ipaddress
    Then request is not blocked

Scenario: Check Address from Black List
    Given ipaddress "10.10.0.1"
    When check ipaddress
    Then request is blocked

Scenario: Delete Network from Black List
    Given network "10.10.0.0/24"
    When delete network from black list

Scenario: Check IP from the black list
    Given ipaddress "10.10.0.1"
    When check address
    Then request is not blocked

Scenario: Delete the same Network from Black List 
    Given network "10.10.0.0/24"
    When delete network from black list
    Then error network not found

Scenario: Delete Network from White list
    Given network "20.20.0.0/24"
    When delete network from white list

Scenario: Delete the same Network from White list
    Given network "20.20.0.0/24"
    When delete network from white list
    Then error network not found

Scenario: Check Login rate
    Given login "login1"
    And password "random"
    And IP "random"
    And delay between request is "0s"
    When send 50 requests
    Then 10 requests are not blocked

Scenario: Check Password rate
    Given login "random"
    And password "password1"
    And IP "random"
    And delay between request is "0s"
    When send 500 requests
    Then 100 requests are not blocked
    
Scenario: Check Password rate
    Given login "random"
    And password "random"
    And IP "10.10.0.1"
    And delay between request is "0s"
    When send 1000 requests
    Then 1 requests are not blocked


  