# Ping CLI Application
The CLI application accepts a hostname or an IP address as its argument, then sends ICMP "echo requests" 
in a loop to the target while receiving "echo reply" messages. It also reports loss and RTT times 
for each sent message.

### Usage
To run the application, use the following command: \
```sudo go run ping.go [-ttl=value] [-ipv6] host```

##### Arguments
1. [optional] ```ttl``` - specify ttl value. If not specified, it shall default to 255.
2. [optional] ```ipv6``` - set this flag to specify IPv6 protocol. Otherwise, IPv4 is assumed.
3. ```host``` - either a hostname or ip address. Must correspond to the protocol specified.

### Example

##### IPv4
```sudo go run src/ping.go -ttl=100 1.1.1.1```

##### IPv6
```sudo go run src/ping.go -ttl=100 -ipv6 ipv6.google.com```

### Note
Administrator privileges required. Must use ```sudo```.
