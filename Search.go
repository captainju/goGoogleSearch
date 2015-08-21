/*
Google search from go
*/
package googlesearch
import (
	"net/http/cookiejar"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"net/url"
	"time"
)

type Result struct {
	Link string
	Name string
}

type pagedSearch struct {
	current_page int
	query        string
	tld          string
	hl           string
	http_client  http.Client
	results      chan Result
}

var (
	url_search = "https://www.google.%s/search?q=%s&num=%d&start=%d&safe=off&hl=%s&filter=0"
	nb_results_per_page = 100
	searchLock = make(chan bool)
)

func Search(query string, tld string, hl string) chan Result {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Jar : jar}
	chanResult := make(chan Result, nb_results_per_page)

	pagedSearch := pagedSearch{current_page: 0, query: query, tld: tld, hl:hl, http_client: client, results: chanResult}

	go doPagedSearch(&pagedSearch)
	searchLock <- false

	return chanResult
}

func doPagedSearch(paged_search *pagedSearch) {
	for {
		<-searchLock
		go releaseSearchLock(searchLock)

		url_str := fmt.Sprintf(url_search,
			paged_search.tld,
			url.QueryEscape(paged_search.query),
			nb_results_per_page,
			paged_search.current_page * nb_results_per_page,
			paged_search.hl)

		req, _ := http.NewRequest("GET", url_str, nil)
		req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0)")
		resp, err := paged_search.http_client.Do(req)
		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			panic(err)
		}

		resp.Body.Close()

		entries := doc.Find("h3 a")

		for _, res_dom := range entries.Nodes {
			rawUrl := res_dom.Attr[0].Val
			parsedUrl, parseUrlErr := url.Parse(rawUrl)
			if parseUrlErr != nil {
				continue
			}
			parsed_query, parse_query_err := url.ParseQuery(parsedUrl.RawQuery)
			if parse_query_err != nil {
				continue
			}
			if parsed_query["q"] == nil {
				continue
			}
			res := Result{Link:parsed_query["q"][0]}
			paged_search.results <- res
		}


		if doc.Find("#nav").Size() == 0 {
			//no navigation bar - one page result
			break
		}
		if doc.Find("#nav .b a").Size() == 1 {
			//either previous or next
			if paged_search.current_page != 0 {
				//no next button
				break;
			}
		}

		paged_search.current_page = paged_search.current_page + 1;
	}
	close(paged_search.results)
}

func releaseSearchLock(t chan bool) {
	//only one search every n seconds
	time.Sleep(time.Second * 3)
	t <- false
}
