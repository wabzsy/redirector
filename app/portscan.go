package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	IPAddr chan string
	Thread chan bool

	wg      sync.WaitGroup
	timeout time.Duration

	total   = 0
	found   = 0
	waiting = 0
)

func (a *App) PortScan() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "portscan",
		Short:   "tcp port scanner",
		Aliases: []string{"ps", "scan"},
		Args:    cobra.ExactArgs(1),
		Example: "./red scan 192.168.1-3.1-128\n" +
			"./red scan -t 30 -p 1-1000,8080 -n 1 10.2.20.1",
	}

	cmd.Flags().IntP("thread", "t", 30, "number of threads")
	cmd.Flags().StringP("port", "p", "21-23,80,88,111,135,139,443,445,1080,1086,1433,1521,2049,3128,3306,3389,6379,8080,8888", "port range")
	cmd.Flags().IntP("timeout", "n", 1, "connection timeout (seconds)")

	cmd.RunE = a.PortScanHandler
	return cmd
}

func (a *App) PortScanHandler(cmd *cobra.Command, args []string) error {
	hostRange := args[0]

	thread, err := cmd.Flags().GetInt("thread")
	if err != nil {
		return err
	}

	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}

	n, err := cmd.Flags().GetInt("timeout")
	if err != nil {
		return err
	}

	Thread = make(chan bool, thread)
	IPAddr = make(chan string, thread)

	makeTasks(hostRange, port)
	timeout = time.Duration(n) * time.Second

	go printStatus()

	for i := 0; i < total; i++ {
		waiting = total - i
		Thread <- true
		wg.Add(1)
		go doScan()
	}

	wg.Wait()
	log.Println("All Done.")
	log.Printf("Scan: %d / Found: %d\n", total, found)
	return nil
}

func doScan() {
	addr := <-IPAddr
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		_ = conn.Close()
		log.Println(addr)
		fmt.Printf("\a")
		found++
	}
	<-Thread
	wg.Done()
}

func printStatus() {
	for {
		fmt.Printf("  [%d/%d/%d]            \r", waiting, found, total)
		time.Sleep(10 * time.Millisecond)
	}
}

func makeTasks(host, port string) int {
	IPList := ParseHosts(host)
	PortList := RemoveRepByMap(ParsePorts(port))
	total = len(IPList) * len(PortList)
	log.Printf("[Host: %d | Port: %d | Total: %d]\n", len(IPList), len(PortList), total)

	go func() {
		for i := range IPList {
			for j := range PortList {
				IPAddr <- fmt.Sprintf("%s:%s", IPList[i], PortList[j])
			}
		}
	}()
	return total
}

func ParsePorts(str string) []string {
	PortList := make([]string, 0)
	s := strings.Split(str, ",")

	for _, v := range s {
		if strings.Contains(v, "-") {
			limit := strings.Split(v, "-")
			if len(limit) != 2 {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			start, err := strconv.Atoi(limit[0])
			if err != nil {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			end, err := strconv.Atoi(limit[1])
			if err != nil {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			if start >= end {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			if start > 65535 || end > 65535 {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			for i := start; i <= end; i++ {
				str := strconv.Itoa(i)
				PortList = append(PortList, str)
			}

		} else {
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Printf("Error: Invaild Port Range [%s]\n", v)
				os.Exit(3)
			}

			if num > 65535 {
				fmt.Printf("Error: Invaild Port Range [%d]\n", num)
				os.Exit(3)
			}

			PortList = append(PortList, v)
		}
	}

	return PortList
}

func ParseHosts(str string) []string {
	IPList := make([]string, 0)
	s := strings.Split(str, ".")
	RangeList := make([][]string, 4)

	if len(s) != 4 {
		fmt.Printf("Error: Invaild IP Range [%s]\n", str)
		os.Exit(2)

	}

	for k, v := range s {
		if strings.Contains(v, "-") {
			limit := strings.Split(v, "-")
			if len(limit) != 2 {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			start, err := strconv.Atoi(limit[0])
			if err != nil {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			end, err := strconv.Atoi(limit[1])
			if err != nil {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			if start >= end {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			if start > 255 || end > 255 {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			for i := start; i <= end; i++ {
				str := strconv.Itoa(i)
				RangeList[k] = append(RangeList[k], str)
			}
		} else {
			num, err := strconv.Atoi(v)
			if err != nil {
				fmt.Printf("Error: Invaild IP Range [%s]\n", v)
				os.Exit(2)
			}

			if num > 255 {
				fmt.Printf("Error: Invaild IP Range [%d]\n", num)
				os.Exit(2)
			}

			RangeList[k] = append(RangeList[k], v)
		}
	}

	for _, a := range RangeList[0] {
		for _, b := range RangeList[1] {
			for _, c := range RangeList[2] {
				for _, d := range RangeList[3] {
					IPList = append(IPList, fmt.Sprintf("%s.%s.%s.%s", a, b, c, d))
				}
			}
		}
	}

	return IPList
}

func RemoveRepByMap(list []string) []string {
	result := make([]string, 0)
	tempMap := map[string]byte{}
	for _, v := range list {
		l := len(tempMap)
		tempMap[v] = 0
		if len(tempMap) != l {
			result = append(result, v)
		}
	}
	return result
}
