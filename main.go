package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	mq "github.com/aravinth2094/GoProxy/socks5/messagequeue"

	"github.com/aravinth2094/GoProxy/socks5"
)

var blacklist map[string]bool = make(map[string]bool)
var initChannel chan interface{} = make(chan interface{}, 1)

const (
	BLACK_LIST_FILE = "blacklist.txt"
)

// DownloadFile downloads a file from the url and saves to local file
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Parse function will parse the host file downloaded from the internet
func Parse(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#") {
			continue
		}
		split := strings.Split(line, " ")
		if len(split) > 1 {
			block(strings.TrimSpace(split[1]))
			count++
		}
	}
	log.Println("Added", count, "blacklist host(s)")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func createBlacklist() {
	log.Println("Downloading latest blacklist...")
	err := DownloadFile(BLACK_LIST_FILE, "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn/hosts")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Parsing the blacklist file...")
	Parse(BLACK_LIST_FILE)
	os.Remove(BLACK_LIST_FILE)
}

func isBlocked(host string) bool {
	return blacklist[host]
}

func block(host string) {
	blacklist[host] = true
}

// GetInitChannel returns a channel that notifies
// if the proxy server has been initialized
func GetInitChannel() <-chan interface{} {
	return initChannel
}

func initialized() {
	initChannel <- struct{}{}
}

func main() {
	bootstrapServers := flag.String("kafka-server", "localhost:9092", "Kafka Bootstrap Server IP:PORT")
	proxyServerPort := flag.String("port", "1080", "Proxy Server Bind Port")
	proxyServerHost := flag.String("host", "0.0.0.0", "Proxy Server Bind Host")
	kafkaEnable := flag.Bool("kafka", true, "Enable/Disable kafka streaming")
	kafkaTopic := flag.String("kafka-topic", "proxyMonitor", "Kafka Streaming Topic")
	help := flag.Bool("h", false, "Print this help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	log.Println("GoProxy Starting...")
	var messageQueue *mq.Producer
	if *kafkaEnable {
		var err error
		messageQueue, err = mq.CreateProducer([]string{*bootstrapServers}, *kafkaTopic)
		if err != nil {
			log.Fatalln(err)
		}
	}
	createBlacklist()
	go func() {
		log.Println("Blacklist auto updater started...")
		ticker := time.NewTicker(24 * time.Hour)
		for {
			<-ticker.C
			createBlacklist()
		}
	}()
	srv := socks5.New()
	srv.AuthNoAuthenticationRequiredCallback = func(c *socks5.Conn) error {
		c.Data = "anonymous"
		return nil
	}
	srv.HandleConnectFunc(func(c *socks5.Conn, host string) (newHost string, err error) {
		domain, _, _ := net.SplitHostPort(host)
		blocked := isBlocked(domain)
		remoteAddress, _, _ := net.SplitHostPort(c.RemoteAddr())
		if *kafkaEnable {
			messageQueue.Send(&mq.Message{
				Domain:        domain,
				RemoteAddress: remoteAddress,
				Timestamp:     time.Now(),
				Blocked:       blocked,
			})
		}
		if blocked {
			log.Println("Blocked", domain)
			return host, socks5.ErrConnectionNotAllowedByRuleset
		}
		return host, nil
	})

	log.Println("GoProxy Server Started")
	initialized()
	log.Fatalln(srv.ListenAndServe(net.JoinHostPort(*proxyServerHost, *proxyServerPort)))
}
