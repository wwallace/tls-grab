package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {

	sc := bufio.NewScanner(os.Stdin)

	jobs := make(chan string)
	var wg sync.WaitGroup

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	for i := 0; i < 20; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range jobs {
				if !strings.Contains(host, ":") {
					host = host + ":443"
				}
				conn, err := tls.Dial("tcp", host, conf)
				if err != nil {
					// log.Println("Error in Dial", err)
					continue
				}
				defer conn.Close()
				certs := conn.ConnectionState().PeerCertificates
				for _, cert := range certs {
					var dnsNames = cert.DNSNames
					for _, name := range dnsNames {
						name = strings.Replace(name, "*.", "", 0)
						fmt.Println(name)
					}
				}
			}

		}()

	}

	for sc.Scan() {
		host := sc.Text()
		jobs <- host
	}

	close(jobs)
	wg.Wait()
}
