package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

func isError(err error) bool {
	if err != nil {
		log.Print("[error] ", err)
		return true
	}
	return false
}

type Game struct {
	ImageURL    string
	Name        string
	Description string
	Results     []string
	BonusBall   string
	Date        string
	TimeOfDay   string
}

func (g *Game) ToString() string {
	return strings.Join(g.Results, ",")
}

func FetchGames() []Game {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("GET", "https://supremeventures.com/game-results/", nil)
	if isError(err) {
		return nil
	}
	request.Header.Set("User-Agent", "Not a browser")

	res, err := httpClient.Do(request)
	if isError(err) {
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var games []Game
	doc.Find(".game-container").Each(func(i int, s *goquery.Selection) {
		var game Game

		//Find game name
		name, exists := s.Find(".game").Attr("class")
		imgURI, _ := s.Find(".game-logo a img").Attr("data-lazy-src")

		if exists {
			var results []string

			game.Name = strings.ToUpper(strings.Replace(name, "game ", " ", 1))
			game.ImageURL = imgURI

			s.Find(".game-content .game-result").Each(func(idx int, r *goquery.Selection) {
				game.Date = r.Find("h4:nth-of-type(1)").Text()
				game.TimeOfDay = strings.ToUpper(r.Find("h5").Text())
				if strings.Trim(game.Name, " ") == "CASHPOT" {
					game.Description = r.Find("p:nth-of-type(2)").Text()
				} else if strings.Trim(game.Name, " ") == "LOTTO" || strings.Trim(game.Name, " ") == "SUPERLOTTO" {
					game.Description = r.Find("span.jackpot").Text()
				}
			})

			// megaBalls.Find(".game-content .megaball").Text()

			//Find game results
			s.Find(".game-content .game-result .result-number").Each(func(idx int, r *goquery.Selection) {
				if strings.Trim(game.Name, " ") == "LOTTO" && idx == int(6) {
					game.BonusBall = r.Text()
				} else if strings.Trim(game.Name, " ") == "SUPERLOTTO" && idx == int(5) {
					game.BonusBall = r.Text()
				} else {
					results = append(results, r.Text())
				}
			})

			game.Results = results
		}

		if len(game.Results) > 0 {
			games = append(games, game)
		}
	})

	return games
}

func DisplayGames(games []Game) {
	for _, g := range games {
		fmt.Printf("Game: %s \n", g.Name)
		fmt.Printf("Image: %s \n", g.ImageURL)
		fmt.Printf("Date: %s \n", g.Date)
		fmt.Printf("Time of Day: %s \n", g.TimeOfDay)
		fmt.Printf("Results: %s \n\n", g.ToString())
	}
}
