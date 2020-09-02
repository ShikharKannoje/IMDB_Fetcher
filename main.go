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
	Duration           string `json:"duration"`
	Genre              string `json:"genre"`
	Summary            string `json:"summary"`
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
	baseURL := "https://www.imdb.com/india/top-rated-indian-movies"
	resp, err := http.Get(baseURL)
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

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					path := strings.Split(a.Val, "/")
					//fmt.Println(path)
					if len(path) < 2 {
						break
					}
					if path[1] != "title" {
						break
					} else {
						if count < movie_count_limit {
							newURL := "https://www.imdb.com" + a.Val
							resp2, err := http.Get(newURL)
							//fmt.Println(newURL, "success")
							if err != nil {
								log.Println(err)
							}
							defer resp.Body.Close()

							doc2, err := html.Parse(resp2.Body)
							if err != nil {
								log.Println(err)
							}
							var f2 func(*html.Node)
							f2 = func(n *html.Node) {
								if n.Type == html.ElementNode && n.Data == "div" {
									text := &bytes.Buffer{}
									for _, a := range n.Attr {
										var DurationAndGenre string
										if a.Key == "class" && a.Val == "subtext" {
											collectText(n, text)
											if text.String() != " " {
												str := text.String()

												DurationAndGenre = strings.TrimFunc(str, func(r rune) bool {
													return !unicode.IsNumber(r)
												})

												DurationAndGenre = strings.ReplaceAll(DurationAndGenre, string(' '), "")
												DurationAndGenre = strings.ReplaceAll(DurationAndGenre, "|", "")
												DurationAndGenre = strings.ReplaceAll(DurationAndGenre, "\t", "")
												t1 := strings.Split(DurationAndGenre, "\n")
												var time string
												var genre string
												time = t1[0]
												genre = t1[3]
												i := 4
												for i < len(t1) {
													if genre[len(genre)-1] != byte(',') {
														break
													} else {
														genre = genre + t1[i]
														i++
													}
												}
												//fmt.Println(time, genre)
												jsonForm.Duration = time
												jsonForm.Genre = genre
											}

										}
										if a.Key == "class" && a.Val == "summary_text" {
											collectText(n, text)
											if text.String() != " " {
												str := text.String()

												summary := strings.TrimFunc(str, func(r rune) bool {
													return !unicode.IsLetter(r)
												})

												jsonForm.Summary = summary
												//fmt.Println(jsonForm)
											}

										}

									}

								}

								for c := n.FirstChild; c != nil; c = c.NextSibling {
									f2(c)
								}

							}

							f2(doc2)

						} else {
							break
						}

					}

				}
			}
		}
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
							//fmt.Println(jsonForm)
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
