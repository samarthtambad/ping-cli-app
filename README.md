# Ping CLI Application
The CLI app accepts a hostname or an IP address as its argument, then sends ICMP "echo requests" 
in a loop to the target while receiving "echo reply" messages. It also reports loss and RTT times 
for each sent message.

## Usage
To run the application, use the following command: \
```sudo go run ping.go [-ttl=value] host```

Administrator privileges required for running this command. Please use ```sudo```.

Note that ```ttl``` is an optional argument. If not specified, it shall default to 255.
