check\_pa\_sessions
============

Nagios check for Palo Alto session statistics

## Usage

```
Usage of check_pa_sessions:
  -H="127.0.0.1": Target host
  -community="public": SNMP community string
  -critical=10000: Critical threshold value
  -max=250000: Maximum value of control
  -min=0: Minimum value of control
  -mode="tcp-sessions": Specify session mode. tcp, udp, icmp, or total
  -timeout=10: SNMP connection timeout
  -warning=10000: Warning threshold value
```


## Examples

```
$> check_pa_sessions -H 1.1.1.1 -community="public" -mode tcp-sessions -timeout 5 --warning 60000 --critical 80000
OK: TCP Sessions - 18522 | tcpsessions=18522;;;60000;80000
```

## References
[Useful SNMP OIDs for Monitoring Palo Alto Networks Devices](https://live.paloaltonetworks.com/docs/DOC-1744)
