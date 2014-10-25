package main

import (
	"flag"
	"fmt"

	"github.com/alouca/gosnmp"
	"github.com/fractalcat/nagiosplugin"
)

func main() {
	var (
		host      string
		community string
		timeout   int64
		warning   int64
		critical  int64
		min       int64
		max       int64
		mode      string
	)

	flag.StringVar(&host, "H", "127.0.0.1", "Target host")
	flag.StringVar(&community, "community", "public", "SNMP community string")
	flag.Int64Var(&timeout, "timeout", 10, "SNMP connection timeout")

	flag.Int64Var(&warning, "warning", 60000, "Warning threshold value")
	flag.Int64Var(&critical, "critical", 80000, "Critical threshold value")

	flag.Int64Var(&min, "min", 0, "Minimum value of control")
	flag.Int64Var(&max, "max", 250000, "Maximum value of control")

	flag.StringVar(&mode, "mode", "tcp-sessions", "Specify session mode. tcp, udp, icmp, or total")

	flag.Parse()

	// Initialize the check - this will return an UNKNOWN result
	// until more results are added.
	check := nagiosplugin.NewCheck()
	// If we exit early or panic() we'll still output a result.
	defer check.Finish()

	// obtain data here
	c, err := gosnmp.NewGoSNMP(host, community, gosnmp.Version2c, timeout)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	var proto string
	var perf string

	if mode == "udp-sessions" {
		proto = "UDP"
		perf = "udpsessions"
	} else if mode == "tcp-sessions" {
		proto = "TCP"
		perf = "tcpsessions"
	} else if mode == "icmp-sessions" {
		proto = "ICMP"
		perf = "icmpsessions"
	} else if mode == "total-sessions" {
		proto = "Total"
		perf = "totalsessions"
	} else {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("Unknown mode: %v", mode))
		return
	}

	sessions, err := getData(c, mode)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	check.AddPerfDatum(perf, "", float64(sessions), float64(warning), float64(critical), float64(min), float64(max))

	crit, _, _, err := parseRange(critical, sessions)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	if crit {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("%s Sessions - %d", proto, sessions))
		return
	}

	warn, _, _, err := parseRange(warning, sessions)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	if warn {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("%s Sessions - %d", proto, sessions))
		return
	}

	check.AddResult(nagiosplugin.OK, fmt.Sprintf("%s Sessions - %d", proto, sessions))
}

func getData(s *gosnmp.GoSNMP, oidType string) (int, error) {
	val := -1

	pkt, err := s.Get(oids[oidType])
	if err != nil {
		return val, err
	}

	for _, v := range pkt.Variables {
		switch v.Type {
		case gosnmp.Integer:
			val = v.Value.(int)
		}
	}

	return val, nil
}

func parseRange(r int64, val int) (bool, float64, float64, error) {
	nr, err := nagiosplugin.ParseRange(fmt.Sprintf("%d", r))
	if err != nil {
		return false, 0, 0, err
	}

	return nr.CheckInt(val), nr.Start, nr.End, nil
}
