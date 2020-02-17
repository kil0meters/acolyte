package links

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/kil0meters/acolyte/pkg/database"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Article struct {
	Title   string    `db:"title"`
	Link    string    `db:"link"`
	Icon    string    `db:"icon"`
	Content string    `db:"content"`
	Date    time.Time `db:"created_at"`
}

func GetArticleInfo(linkStr string) Article {
	linkFull, _ := url.Parse(linkStr)

	article := Article{
		Icon: linkFull.Scheme + "://" + linkFull.Host + "/icon.ico",
		Link: linkFull.Scheme + "://" + linkFull.Host + linkFull.Path,
	}

	err := database.DB.QueryRowx("SELECT * FROM link_cache WHERE link = $1", article.Link).StructScan(&article)
	if err != nil { // if it's not in the database, fetch it
		res, err := http.Get(linkStr)
		if err != nil {
			log.Println("Error fetching article")
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Println("Encountered error", res.StatusCode)
		}

		document, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println("Encountered error", err.Error())
		}

		document.Find("meta[itemprop=\"headline\"]," +
			"meta[property=\"og:title\"]").First().Each(func(index int, selection *goquery.Selection) {
			article.Title, _ = selection.Attr("content")
		})

		document.Find("meta[itemprop=\"description\"]," +
			"meta[property=\"og:description\"]").First().Each(func(index int, selection *goquery.Selection) {
			article.Content, _ = selection.Attr("content")
		})

		document.Find("meta[itemprop=\"datePublished\"]," +
			"meta[property=\"article:published\"]," +
			"meta[name=\"last_updated_date\"]," +
			"meta[data-hid=\"dcterms.created\"]," +
			"meta[property=\"article:published_time\"]").First().Each(func(index int, selection *goquery.Selection) {
			timeString, _ := selection.Attr("content")
			article.Date, _ = time.Parse(time.RFC3339, timeString)
		})

		if article.Date.IsZero() {
			article.Date = time.Now()
		}

		_, err = database.DB.Exec("INSERT INTO link_cache (title, link, icon, content, created_at) VALUES ($1, $2, $3, $4, $5)",
			article.Title, article.Link, article.Icon, article.Content, article.Date)
		if err != nil {
			log.Println(err)
		}
	}

	return article
}
