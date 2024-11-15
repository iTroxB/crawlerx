package main

import (
        "crypto/sha256"
        "flag"
        "fmt"
        "github.com/gocolly/colly"
        "log"
        "math/rand"
        "net/url"
        "os"
        "strings"
        "sync"
        "time"
)

var (
        domain     string
        outputFile string
        crawlDepth int
        maxDepth   = 10
        useHTTP    bool
        userAgents  = []string{
                "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; Maxthon 2.0)",
                "Mozilla/5.0 (Android 11; Mobile; rv:115.0) Gecko/115.0 Firefox/115.0",
                "Mozilla/5.0 (iPad; CPU OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
                "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
                "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
                "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
                "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 EdgA/114.0.0.0",
                "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 OPR/63.3.3216.58675",
                "Mozilla/5.0 (Linux; Android 11; Redmi Note 9 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
                "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/16.0 Chrome/114.0.0.0 Mobile Safari/537.36",
                "Mozilla/5.0 (Linux; Android 12; SM-G991U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
                "Mozilla/5.0 (Linux; U; Android 10; en-us; SM-G970U Build/QP1A.190711.020) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:115.0) Gecko/20100101 Firefox/115.0",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_6_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 10.0; ARM64; rv:115.0) Gecko/20100101 Firefox/115.0",
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.0.0",
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 OPR/80.0.4170.72",
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
                "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.67 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.29 Safari/537.36 OPR/15.0.1147.24 (Edition Next)",
                "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.71 (KHTML like Gecko) WebVideo/1.0.1.10 Version/7.0 Safari/537.71",
                "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:35.0) Gecko/20100101 Firefox/35.0",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML like Gecko) Maxthon/4.0.0.2000 Chrome/22.0.1229.79 Safari/537.1",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.12 Safari/537.36 OPR/14.0.1116.4",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/28.0.1469.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.76 Safari/537.36 OPR/19.0.1326.56",
                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.154 Safari/537.36 OPR/20.0.1387.91",
                "Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1667.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/28.0.1469.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0",
                "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2049.0 Safari/537.36",
                "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36 OPR/18.0.1284.49",
                "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.19.4 (KHTML, like Gecko) Version/5.0.2 Safari/533.18.5",
                "Mozilla/5.0 (Windows; U; Windows NT 6.2; es-US ) AppleWebKit/540.0 (KHTML like Gecko) Version/6.0 Safari/8900.00",
                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
                "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
                "Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
                "Opera/9.80 (Windows NT 6.1; U; en) Presto/2.7.62 Version/11.01",
                "Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.16",
        }
)

func showBanner() {
	fmt.Println("\033[1;36m")
	fmt.Println("  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ ██╗  ██╗ ")
	fmt.Println(" ██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗╚██╗██╔╝ ")
	fmt.Println(" ██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝ ╚███╔╝ ")
	fmt.Println(" ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗ ██╔██╗ ")
	fmt.Println(" ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║██╔╝ ██╗ ")
	fmt.Println("  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ \n")
	fmt.Println("  Web crawling tool for domains")
	fmt.Println("  Version 1.0")
	fmt.Println("  Made by iTrox")
        fmt.Println("  Website: https://www.itrox.site")
        fmt.Println("  Blog: https://itrox.gitbook.io/itrox\n")
	fmt.Println("  crawlerx [-h] to view help menu")
	fmt.Println("\033[0m")
}

func showUsage() {
	showBanner()
	fmt.Println("Usage of Crawlerx:\n")
	fmt.Println("  -d <domain>              Domain to crawl (e.g., example.com)")
	fmt.Println("  -o <output file>         Output file")
	fmt.Println("  -D <depth crawling>      Crawling depth (default: 5, maximum: 10)")
	fmt.Println("  -H                       Forcing HTTP connection instead of HTTPS\n")
	// flag.PrintDefaults()
}

func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func writeOutput(output *os.File, text string) {
	fmt.Println(text)
	if output != nil {
		_, err := output.WriteString(text + "\n")
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}
	}
}

func normalizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	parsedURL.Fragment = ""

	if parsedURL.RawQuery != "" {
		values, _ := url.ParseQuery(parsedURL.RawQuery)
		parsedURL.RawQuery = values.Encode()
	}

	if parsedURL.Path == "" {
		parsedURL.Path = "/"
	}

	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	return parsedURL.String()
}

func setupCrawler(domain string, outputFile string, depth int) {
	c := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.Async(true),
		colly.AllowedDomains(domain),
	)

	c.UserAgent = randomUserAgent()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	var output *os.File
	var err error
	if outputFile != "" {
		output, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening the output file: %v", err)
		}
		defer output.Close()
	}

	visitedURLs := make(map[[32]byte]struct{})
	var mu sync.RWMutex

	c.OnRequest(func(r *colly.Request) {
		if r.Headers.Get("Referer") == "" {
			r.Headers.Set("Referer", fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host))
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = e.Request.AbsoluteURL(link)
		if strings.HasPrefix(link, "http") {
			parsedLink, err := url.Parse(link)
			if err != nil {
				return
			}

			linkHost := strings.ToLower(parsedLink.Hostname())
			domainHost := strings.ToLower(domain)

			if linkHost != domainHost {
				return
			}

			normalizedLink := normalizeURL(link)
			hash := sha256.Sum256([]byte(normalizedLink))

			mu.RLock()
			_, visited := visitedURLs[hash]
			mu.RUnlock()

			if !visited {
				mu.Lock()
				visitedURLs[hash] = struct{}{}
				mu.Unlock()

				writeOutput(output, normalizedLink)
				e.Request.Visit(normalizedLink)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if r != nil && (r.StatusCode == 403 || r.StatusCode == 502) {
			return
		}
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			return
		}
	})

	protocol := "https"
	if useHTTP {
		protocol = "http"
	}

	startURL := fmt.Sprintf("%s://%s", protocol, domain)
	normalizedStartURL := normalizeURL(startURL)
	hash := sha256.Sum256([]byte(normalizedStartURL))

	mu.Lock()
	visitedURLs[hash] = struct{}{}
	mu.Unlock()

	err = c.Visit(startURL)
	if err != nil {
		log.Fatalf("Error starting crawling: %v", err)
	}

	c.Wait()
}

func main() {
	flag.StringVar(&domain, "d", "", "Domain to crawl")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.IntVar(&crawlDepth, "D", 5, "Crawling depth")
	flag.BoolVar(&useHTTP, "H", false, "Forcing HTTP connection instead of HTTPS")

	flag.Usage = showUsage

	flag.Parse()

	if domain == "" {
		flag.Usage()
		return
	}

	if crawlDepth > maxDepth {
		crawlDepth = maxDepth
	}

	rand.Seed(time.Now().UnixNano())

	showBanner()
	setupCrawler(domain, outputFile, crawlDepth)
}
