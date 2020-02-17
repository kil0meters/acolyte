package links

import (
	"testing"
)

func TestLinks(t *testing.T) {
	article := GetArticleInfo("https://www.cnn.com/2020/02/14/politics/evelyn-yang-doctor-sexual-assault-accusers-invs/index.html")
	t.Logf("CNN: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.nytimes.com/2020/02/14/us/politics/andrew-mccabe-william-barr-flynn.html")
	t.Logf("The New York Times: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.washingtonpost.com/politics/bernie-sanders-is-powered-by-a-loyal-base-but-results-in-iowa-and-new-hampshire-show-the-movement-has-limits/2020/02/14/e10c570a-4296-11ea-b5fc-eefa848cde99_story.html")
	t.Logf("The Washington Post: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.foxnews.com/politics/bernie-sanders-russian-pranksters-greta-thunberg")
	t.Logf("Fox News: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.wired.com/story/marathon-investigation-cheaters-racing-data/")
	t.Logf("Wired: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.theatlantic.com/family/archive/2020/02/valentines-day-everyone-relationships/606580/")
	t.Logf("The Atlantic: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)

	article = GetArticleInfo("https://www.theverge.com/2020/2/14/21137852/reddit-female-dating-advice-strategy-women-rulebook-memes")
	t.Logf("The Verge: title: \"%s\", date: \"%s\", description: \"%s\"", article.Title, article.Date, article.Content)
}

func TestTwitter(t *testing.T) {

}

func TestReddit(t *testing.T) {

}
