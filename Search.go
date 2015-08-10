package googlesearch
import (
	"net/http/cookiejar"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

type Result struct {
	position int
	link string
	name string
}

var (
	url_search = "https://www.google.fr/search"

)

func Search(query string) []Result {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Jar : jar}
	resp, err := client.Get(url_search+"?q="+query)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)

	fmt.Println("nb result ", doc.Find("h3").Size())

	var results []Result

	for index, res_dom := range doc.Find("h3 a").Nodes {
		fmt.Println(res_dom.Attr)
		res := Result{position : index, link:"www.google.fr"}
		results = append(results, res)
	}


	//fmt.Println(doc.Find("p").First().Text())


	return results
}
