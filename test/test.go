package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"golang.org/x/net/html"
)

type Book struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Genres struct {
	Name string `json:"genre"`
	Url  string `json:"url"`
	Book []Book `json:"book"`
}

func main() {

	fmt.Println("Connecting to the website!...")
	
	resp, err := http.Get("https://www.webtoons.com/en/genres/")
	if err != nil {
		fmt.Println("Error fetching WEBTOON homepage: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML: ", err)
		return
	}
	fmt.Println("Downloading... ")
	genres := getGenres(doc)
	var allGenre []Genres
	var allBook []Book

	

	for _, genre := range genres {
		books := getTiltle(genre.Url)
		allBook = append(allBook, books...)
		genre.Book = allBook
		allGenre = append(allGenre, genre)
		allBook = []Book{}
	}
	writeToJson("manga_info.json", allGenre)
}

func getGenres(n *html.Node) []Genres {
	var genres []Genres
	var currentGenre Genres
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "data-genre" {
					currentGenre.Name = attr.Val
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					currentGenre.Url = attr.Val
					if currentGenre.Name != "" && currentGenre.Url != "" && currentGenre.Name != "OTHERS" {
						genres = append(genres, currentGenre)
						currentGenre = Genres{}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return genres
}

func getTiltle(url string) []Book {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching WEBTOON homepage: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		os.Exit(1)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML: ", err)
		os.Exit(1)
	}
	var book []Book
	var currentBook Book
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					currentBook.Url = attr.Val
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "alt" {
					currentBook.Title = attr.Val
					if currentBook.Title != "" {
						book = append(book, currentBook)
						currentBook = Book{}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return book
}

func writeToJson(filename string, key interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating JSON file: ", err)
		return
	}
	defer file.Close()

	data, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		fmt.Println("Error encoding books to JSON: ", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing JSON to file: ", err)
		return
	}
	fmt.Println("Books have been written to books.json")
}
