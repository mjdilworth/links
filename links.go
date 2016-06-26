package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/asaskevich/govalidator"
)

//url link structure for work queue - type of is what tupe is it
type urlLink struct {
	url     string
	typeof  string
	isAsset bool // ot going to follow
}

//Links - structure to hold my links in. with mutex to allow safe writes
type Links struct {
	entries map[string]string
	mux     *sync.Mutex
}

// NewLinkmap Constructs a new NewLinkMap.
func NewLinks() *Links {
	return &Links{
		entries: make(map[string]string),
		mux:     &sync.Mutex{},
	}
}
func main() {

	//input default values
	URLPtr := flag.String("url", "http://wiprodigital.com", "a valide URL")
	//parse input
	flag.Parse()

	//using govalidator - does it add enough value??
	validURL := govalidator.IsURL(*URLPtr)
	if validURL == false {
		fmt.Printf("%s is NOT a valid URL\n", *URLPtr)
		return
	}

	//print something out to say we are working
	fmt.Println("Starting....")
	fmt.Println("")
	//get my links
	linkmap, err := GetLinks(*URLPtr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//print out results
	linkmap.PrintLinks()

}

// GetLinks crawls a start URL for all links and assets and builds
// a links struct with pages and assets per crawled link.
func GetLinks(url string) (*Links, error) {
	links := NewLinks()

	wg := &sync.WaitGroup{}

	//the i have completed channel
	done := make(chan bool)
	//map to keep theones i have already seen
	seen := make(map[string]bool)
	//the queue of links ot be processed
	queue := make(chan *urlLink)

	//get my domain
	parentDomain, err := GetDomain(url)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	//my first link
	seen[url] = true
	go getURL(url, queue, wg)

	// Waits for all goroutines to finish and signals the fact to
	// the `done` channel in order to terminate the select loop.
	go func() {
		wg.Wait()
		done <- true
	}()

	for {
		select {
		case link := <-queue: //pop off queue
			linkDomain, err := GetDomain(link.url)
			if err != nil {
				continue
			}

			// i will record off site links and binary files but i wont follow
			// We only allow assets to come from a different domain.
			if !link.isAsset && linkDomain != parentDomain {
				links.AddEntry(link.url, "External")
				continue
			}

			// No need to fetch assets.
			if link.isAsset {
				links.AddEntry(link.url, "Asset")
				continue
			}

			// Ensures we don't visit URLs twice.
			if seen[link.url] {
				continue
			}
			//defult is a page that we want to see
			links.AddEntry(link.url, "Page")
			seen[link.url] = true

			wg.Add(1)
			go getURL(link.url, queue, wg)
		case <-done:
			return links, nil
		}
	}
}

// AddEntry - Adds a `Link` to the sitemap in a thread-safe manner.
func (s *Links) AddEntry(url, typeof string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.entries[url] = typeof
}

// PrintLinks - to print out our links is some form of oerder.
func (s *Links) PrintLinks() {
	//all External
	fmt.Println("All External Links")
	for key, value := range s.entries {
		if value == "External" {
			fmt.Printf("%s\n", key)
		}
	}
	fmt.Println("")
	//all assets
	fmt.Println("All Assets")
	for key, value := range s.entries {
		if value == "Asset" {
			fmt.Printf("%s\n", key)
		}
	}
	fmt.Println("")
	//all pages
	fmt.Println("All Pages")
	for key, value := range s.entries {
		if value == "Page" {
			fmt.Printf("%s\n", key)
		}
	}
}

// getURL a URL and enqueues extracted links for further processing.
func getURL(url string, queue chan *urlLink, wg *sync.WaitGroup) {
	defer wg.Done()

	// i will remove this from final code
	//fmt.Println("Visiting", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	links := ExtractLinks(url, resp.Body)

	for _, link := range links {
		queue <- link
	}
}

//AbsoluteURL Transforms a URL to an absolute URL given its parent.
func AbsoluteURL(href, parent string) (string, error) {
	url, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	parentURL, err := url.Parse(parent)
	if err != nil {
		return "", err
	}
	resolved := parentURL.ResolveReference(url)
	return resolved.String(), nil
}

// ExtractLinks - Extracts and returns a list of absolute URLs (links and assets)
// from an HTML document.
func ExtractLinks(url string, body io.Reader) []*urlLink {
	links := make([]*urlLink, 0)
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()
			isAsset := false
			href := ""
			if t.Data == "a" {
				href = extractAttr(t, "href")
			} else if t.Data == "script" {
				href = extractAttr(t, "src")
				isAsset = true
			} else if t.Data == "link" {
				// Extract link tags but limit the set to just stylesheets.
				rel := extractAttr(t, "rel")
				if rel != "stylesheet" {
					continue
				}
				href = extractAttr(t, "href")
				isAsset = true
			}
			if href == "" {
				continue
			}
			//here i need to cut out any # links
			if strings.HasPrefix(strings.ToLower(href), "#") {
				continue
			}
			//here i need to cut out any mailto links
			if strings.HasPrefix(strings.ToLower(href), "mailto") {
				continue
			}
			href, err := AbsoluteURL(href, url)
			if err != nil {
				continue
			}
			//typeof
			//no exe, zip,
			if strings.HasSuffix(strings.ToLower(href), ".exe") {

				isAsset = true
			}
			if strings.HasSuffix(strings.ToLower(href), ".zip") {

				isAsset = true
			}

			if strings.HasSuffix(strings.ToLower(href), ".pdf") {

				isAsset = true
			}
			if strings.HasSuffix(strings.ToLower(href), ".png") {

				isAsset = true
			}
			if strings.HasSuffix(strings.ToLower(href), ".jpg") {

				isAsset = true
			}

			//add to be processed..
			links = append(links, &urlLink{href, url, isAsset})
		}
	}
}

// Given a URL it returns its domain.
func GetDomain(href string) (string, error) {
	url, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	tokens := strings.Split(url.Host, ":")
	return tokens[0], nil
}

// Extracts an attribute from an HTML token. Returns an empty
// string if the attribute is not found.
func extractAttr(t html.Token, attr string) string {
	for _, a := range t.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
