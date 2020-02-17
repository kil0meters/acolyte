package links

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Article struct {
	Title   string    `db:"title" json:"title"`
	Link    string    `db:"link" json:"link"`
	Icon    string    `db:"icon" json:"icon"`
	Content string    `db:"content" json:"content"`
	Date    time.Time `db:"created_at" json:"published_date"`
}

func GetArticleInfo(linkStr string) Article {
	linkFull, err := url.Parse(linkStr)
	if err != nil || linkFull.Scheme == "" {
		return Article{}
	}

	log.Println(linkFull)

	article := Article{
		Icon: linkFull.Scheme + "://" + linkFull.Host + "/icon.ico",
		Link: linkFull.Scheme + "://" + linkFull.Host + linkFull.Path,
	}

	err = database.DB.QueryRowx("SELECT * FROM link_cache WHERE link = $1", article.Link).StructScan(&article)
	if err != nil { // if it's not in the database, fetch it
		res, err := http.Get(linkFull.String())
		if err != nil {
			log.Printf("Error fetching article \"%s\"\n", linkFull)
			return Article{}
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

func LinkSearch(w http.ResponseWriter, r *http.Request) {
	account := authorization.GetAccount(w, r)

	if account.Permissions.AtLeast(authorization.LoggedOut) {
		params := mux.Vars(r)

		json.NewEncoder(w).Encode(GetArticleInfo(params["link"]))
	} else {
		w.Write([]byte("{\"error\":\"unauthorized\"}"))
	}
}
