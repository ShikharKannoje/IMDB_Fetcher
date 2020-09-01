package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/golang/go/src/pkg/fmt"
	"github.com/golang/go/src/pkg/log"
	"golang.org/x/net/html"
)

// type MovieDetailinput struct {
// 	Rank  string
// 	Movie string
// 	Year  string
// }

type MovieDetail struct {
	Rank               string `json:"rank"`
	Title              string `json:"title"`
	Movie_release_year string `json:"movie_release_year"`
	IMDB_rating        string `json:"imdb_rating"`
}

func sendForMarshal(jsonForm MovieDetail) {
	js, _ := json.Marshal(jsonForm)
	fmt.Println(string(js))
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
func main() {
	resp, err := http.Get("https://www.imdb.com/india/top-rated-indian-movies")
	var movie_count_limit = 2
	var count = 0
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var jsonForm MovieDetail
	// Recursively visit nodes in the parse tree
	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "td" {
			text := &bytes.Buffer{}
			for _, a := range n.Attr {
				var titleAndYear, imdbRatings string
				if a.Key == "class" && a.Val == "titleColumn" {
					collectText(n, text)
					if text.String() != " " {
						str := text.String()

						titleAndYear = strings.TrimFunc(str, func(r rune) bool {
							return !unicode.IsLetter(r) && !unicode.IsNumber(r)
						})
						titleAndYearSlice := strings.Split(titleAndYear, "\n")

						jsonForm.Rank = strings.Replace(titleAndYearSlice[0], ".", "", -1)
						jsonForm.Title = strings.Trim(titleAndYearSlice[1], " ")
						jsonForm.Movie_release_year = strings.Replace(strings.Trim(titleAndYearSlice[2], " "), "(", "", -1)
						//jsonForm = jsonBuilder(titleAndYear)
						//fmt.Println(jsonForm)
					}

				}

				if a.Key == "class" && a.Val == "ratingColumn imdbRating" {
					collectText(n, text)
					if text.String() != " " {
						str := text.String()

						imdbRatings = strings.TrimFunc(str, func(r rune) bool {
							return !unicode.IsLetter(r) && !unicode.IsNumber(r)
						})
						//jsonForm.IMDB_rating = imdbRatings
						//jsonBuilder(str2)
						//fmt.Println("Moive Ratings : ", imdbRatings)
						jsonForm.IMDB_rating = imdbRatings
						if count < movie_count_limit {
							count++
							sendForMarshal(jsonForm)
						} else {
							break
						}

					}
				}

				//jsonbuilt = titleAndYear //+ "\n" + imdbRatings + "\n"
				//fmt.Println(jsonbuilt)
				//jsonBuilder(jsonbuilt)

			}

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

	}

	f(doc)
}

// func main() {
// 	//if the caller didn't provide a URL to fetch...
// 	if len(os.Args) < 2 {
// 		//print the usage and exit with an error
// 		fmt.Printf("usage:\n  pagetitle <url>\n")
// 		os.Exit(1)
// 	}

// 	URL := os.Args[1]
// 	//GET the URL
// 	resp, err := http.Get(URL)

// 	//if there was an error, report it and exit
// 	if err != nil {
// 		//.Fatalf() prints the error and exits the process
// 		log.Fatalf("error fetching URL: %v\n", err)
// 	}

// 	//make sure the response body gets closed
// 	defer resp.Body.Close()

// 	//check response status code
// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatalf("response status code was %d\n", resp.StatusCode)
// 	}

// 	//check response content type
// 	ctype := resp.Header.Get("Content-Type")
// 	if !strings.HasPrefix(ctype, "text/html") {
// 		log.Fatalf("response content type was %s not text/html\n", ctype)
// 	}

// 	//create a new tokenizer over the response body
// 	tokenizer := html.NewTokenizer(resp.Body)

// 	//loop until we find the title element and its content
// 	//or encounter an error (which includes the end of the stream)
// 	for {
// 		//get the next token type
// 		tokenType := tokenizer.Next()

// 		//if it's an error token, we either reached
// 		//the end of the file, or the HTML was malformed
// 		if tokenType == html.ErrorToken {
// 			err := tokenizer.Err()
// 			if err == io.EOF {
// 				//end of the file, break out of the loop
// 				break
// 			}
// 			//otherwise, there was an error tokenizing,
// 			//which likely means the HTML was malformed.
// 			//since this is a simple command-line utility,
// 			//we can just use log.Fatalf() to report the error
// 			//and exit the process with a non-zero status code
// 			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
// 		}

// 		//...existing looping and
// 		//error-checking code from above...

// 		//if this is a start tag token...
// 		if tokenType == html.StartTagToken {
// 			//get the token
// 			token := tokenizer.Token()
// 			//if the name of the element is "title"
// 			if "head" == token.Data {
// 				//the next token should be the page title
// 				tokenType = tokenizer.Next()
// 				//just make sure it's actually a text token
// 				fmt.Println(tokenizer.Token().Data)
// 				break
// 				// if tokenType == html.TextToken {
// 				// 	//report the page title and break out of the loop
// 				// 	fmt.Println(tokenizer.Token().Data)
// 				// 	break
// 				// }
// 			}
// 		}
// 	}

// }

// package main

// import (
// 	"bytes"
// 	"fmt"
// )

// func main() {

// 	var n int
// 	fmt.Scanln(&n)
// 	fmt.Println(n)
// 	abc := []int{1, 2, 3}
// 	fmt.Println(change(abc))
// }

// func change(values []int) string {
// 	var buf bytes.Buffer
// 	buf.WriteByte('[')
// 	for i, v := range values {
// 		if i > 0 {
// 			buf.WriteString(", ")
// 		}
// 		fmt.Fprintf(&buf, "%d", v)
// 		//buf.WriteByte(string(v))
// 	}
// 	buf.WriteByte(']')
// 	return buf.String()
// }
