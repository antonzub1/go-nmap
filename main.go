package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "sync"

    "github.com/sirupsen/logrus"
)

type portsRange struct {
    startPort uint64;
    endPort uint64;
}

func parsePortsRange (portArgs string) (*portsRange) {
    s := strings.Split(portArgs, "-");
    startPort, _ := strconv.ParseUint(s[0], 0, 16);
    endPort, _ := strconv.ParseUint(s[1], 0, 16);
    p := &portsRange {
        startPort: startPort,
        endPort: endPort,
    };
    return p;
}

func scanHost(hostname *string, ports *portsRange) {
    var wg sync.WaitGroup;
    for port := ports.startPort; port <= ports.endPort; port++ {
        address := fmt.Sprintf(*hostname + ":%d", port);
        wg.Add(1);
        go func (addr *string) {
            defer wg.Done()
            connection, err := net.Dial("tcp", *addr);
            if err != nil {
                logrus.Error("Unable to connect to ", *addr);
            } else {
                logrus.Info("Connected to ", *addr);
                connection.Close();
            }
        }(&address);
    }
    wg.Wait()
}

func main () {
    portsRangeString := flag.String("range", "", "Ports range to scan");
    hostname := flag.String("hostname", "", "Hostname to connect");

    flag.Parse();
    if *portsRangeString == "" || *hostname == "" {
        flag.PrintDefaults();
        os.Exit(1);
    }
    ports := parsePortsRange(*portsRangeString);
    scanHost(hostname, ports);
}
