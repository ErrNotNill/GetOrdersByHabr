package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/gocolly/colly/v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var SecretKey = "AA1379657aa"
var Token = "vk1.a.uuGYW0HLLXGopIeCS0mwuHG6GnuBJwDMa8GaYgIwy07nKgjNRrE-gTUUg868oQvep7pjozOiTixAD9j_CpZXgxbCY37NFWoIm392Mxp41-4XfX6U_fhOXH2fg_o50dbGIJtRCNsH7J8YdHNxwcVYOhHC24X2lCqf8YghKs8FWJa4r7sW6u7qK4W_3CR2uSnT3hiXR9azxqjp8ZW89VyMSQ"
var ConfirmationToken = "1ee66a1c"

var newUserId int
var newPostId int

// KeyWord for searching in titles from orders.
var KeyWord string

type Fields struct {
	Fields struct {
		TITLE              string `json:"TITLE"`
		NAME               string `json:"NAME"`
		COMMENTS           string `json:"COMMENTS"`
		SOURCE_DESCRIPTION string `json:"SOURCE_DESCRIPTION"`
	} `json:"fields"`
}

const UrlDb = `postgres://postgres:postgres@localhost:5432/postgres`
const Url = `https://onviz.bitrix24.site/`

const TokenUnknown = "1953280162:AAFMVzq63WHhr_KkNjwGgObHbI4PbQcmQqg"

var Input string
var Href string

var ArrayWithoutDuplicates = make([]string, 0)

/*func ifPrefix(absoluteURL string) (b bool) {

	switch false {
	case strings.HasPrefix(absoluteURL, "https://www"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false
	case strings.HasPrefix(absoluteURL, "mailto"):
		b = false
	case strings.HasPrefix(absoluteURL, "tel"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false
	case strings.HasPrefix(absoluteURL, "https://b24"):
		b = false

	}
	return true

}*/

/*func checkForValue(userValue int, students map[string]int) string {
	//traverse through the map
	for key, value := range students {
		//check if present value is equals to userValue
		if value == userValue {
			//if same return true
			return key
		}
	}
	//if value not found return false
	return ""
}*/

/*func SearchEverywhere(sites map[string]int) {
	c := colly.NewCollector()
	for k, _ := range sites {

	}
}*/

func main() {

	fmt.Println("Server started")
	cb := callback.NewCallback()
	fmt.Println("Callback service started")
	cb.ConfirmationKey = ConfirmationToken
	cb.SecretKey = SecretKey

	//var UrlOnUser = "https://vk.com/id"
	//var UrlOnPost = "https://vk.com/onviz?w=wall-165775952_"

	fmt.Println("Confirmation accepted")

	cb.WallReplyNew(func(ctx context.Context, obj events.WallReplyNewObject) {

		newUserId = obj.FromID
		newPostId = obj.PostID

		convPostToStr := strconv.Itoa(newPostId)
		convUserIdToStr := strconv.Itoa(newUserId)

		UrlOnUser := fmt.Sprintf("https://vk.com/id%v", convUserIdToStr)
		UrlOnPost := fmt.Sprintf("https://vk.com/onviz?w=wall-165775952_%v", convPostToStr)

		tn := &Fields{struct {
			TITLE              string `json:"TITLE"`
			NAME               string `json:"NAME"`
			COMMENTS           string `json:"COMMENTS"`
			SOURCE_DESCRIPTION string `json:"SOURCE_DESCRIPTION"`
		}{TITLE: "Комментарий из ВК", NAME: UrlOnUser, COMMENTS: obj.Text, SOURCE_DESCRIPTION: UrlOnPost}}

		jsnm, err := json.Marshal(tn)
		if err != nil {
			log.Println("Error to convert json fields from struct")
		}
		r := bytes.NewReader(jsnm)

		if obj.FromID != 628998745 {
			_, err = http.Post("https://onviz.bitrix24.ru/rest/13938/pqq6j4ohvutvzfmi/crm.lead.add", "application/json", r)
			if err != nil {
				log.Println("Error http:post request to Bitrix24")
			}
			log.Println("Lead was send to Bitrix24")
			log.Printf("User Url: %v / Post Url: %v / Comment Text: %v", UrlOnUser, UrlOnPost, obj.Text)
		}
	})

	http.HandleFunc("/callback", cb.HandleFunc)
	http.HandleFunc("/parse", testHandleFunc)
	http.HandleFunc("/test", newTestHandleFunc)

	//cb.HandleFunc()
	fmt.Println("Server started on port")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Server started with error")
		return
	}
}

func Par(input string) (text []string, hrefs []string, rmvdDupls []string, data *FullData) {
	newFullArr := make([]string, 0)
	hrefsArray := make([]string, 0)

	textArr := make([]string, 0)
	hrefsArr := make([]string, 0)
	//fullNewArr := make([]string,0)

	for k, _ := range IndexSites {

		var fullData string

		c := colly.NewCollector()
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {

			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			var re = regexp.MustCompile(`[[:punct:]]`)
			str45 := re.ReplaceAllString(e.Text, "")
			//fmt.Println(str45)

			fullData = str45

			founded := searchString(fullData, input)
			if founded == true {

				data = &FullData{Title: fullData, Hrefs: absoluteURL}

				textArr = append(textArr, fullData)
				hrefsArr = append(hrefsArr, absoluteURL)
				//fmt.Println(textArr)
				//fmt.Println(hrefsArr)

				fullDataPlus := fmt.Sprintf("%v %v", fullData, absoluteURL)

				if strings.Contains(fullDataPlus, "http") {
					newFullArr = append(newFullArr, fullDataPlus)
					hrefsArray = append(hrefsArray, k)
				}

				//fmt.Println(newFullArr)
				//fmt.Println(hrefsArray)
			} else {
				//to do something if keyword not founded
			}
		})

		err := c.Visit(k)
		if err != nil {
			log.Println(err.Error())
		}
	}

	rmvdDupls = removeDuplicateStr(newFullArr)

	//here we get massive with title and hrefs

	return text, hrefs, rmvdDupls, data
}

func newTestHandleFunc(w http.ResponseWriter, r *http.Request) {
	url := "https://onviz.bitrix24.site/"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	fmt.Println(responseString)

}

type ViewData struct {
	Title     []string
	Link      []string
	Available bool
}

func testHandleFunc(w http.ResponseWriter, r *http.Request) {
	//Linked := "WorldCup"

	UrlName := r.URL.Query().Get("q")
	ts, err := template.ParseFiles("cmd/search.html")

	if err != nil {
		fmt.Println("Error parsefile form.html")
	}

	//TextArray, hrefsArray := Par(UrlName)

	_, _, rmvDupls, data := Par(UrlName)

	var Linked string
	var Title string

	var ArrTitle = make([]string, 0)
	var ArrLink = make([]string, 0)

	//var CreatedLink string

	for _, text := range rmvDupls {

		//index := strings.Index(text,"https")

		strFields := strings.Fields(text)
		//fmt.Println("STRFIELDS", strFields)

		Linked = strFields[len(strFields)-1]
		//fmt.Println("TITLE", Linked)
		//Text = strFields[len(strFields)-]

		ArrTitle = append(ArrTitle, Linked)
		if strings.Contains(text, "https") {
			Title = text
			//fmt.Println("LINK", Title)
			ArrLink = append(ArrLink, Title)
		}

		//CreatedLink = fmt.Sprintf(`<a href="%s"> %s</a>`, Title, Linked)

		//CreateLink = strFields[1]
		//ShowText = strFields[0]

	}

	fmt.Println(rmvDupls)

	//SomeArr := make([]string, 0)

	//SomeArr = append(SomeArr, data.Title, data.Hrefs)

	if data != nil {
		data = &FullData{Hrefs: Href, SomeArr: rmvDupls, Available: true}
	}

	err = ts.Execute(w, data)
	if err != nil {
		fmt.Println("Execute html error")
		return
	}

	//for _, v := range TextArray {

	//_, err := fmt.Fprintf(w, "%v", v)
	//	if err != nil {
	//	fmt.Println("Error Fprintf")
	//		return
	//	}

	//w.Write([]byte(v))
	//http.Redirect(w, r, "https://google.com", http.StatusSeeOther)
	//w.Write([]byte("Wwrite Byte"))
	//fmt.Println("Hello")
	//fmt.Println(r.Body)
	//fmt.Println(v)
	//}

	//age := r.URL.Query().Get("age")

	//fmt.Println(UrlName, "Url.Query.Get is work")

	//fmt.Println("Parse func started")

	//if r != nil {
	//	fmt.Println("query was send")
	//}

	//fmt.Println(r.Body)
	//UrlName := r.FormValue("query")
	//fmt.Println(UrlName)

	/*if r.Method == "POST" {

		for _, v := range hrefsArray {
			w.Write([]byte(v))
			//http.Redirect(w, r, "https://google.com", http.StatusSeeOther)
			//w.Write([]byte("Wwrite Byte"))
			fmt.Println("Hello")
			fmt.Println(r.Body)
			//fmt.Println(v)
		}

	}*/

	/*jsnm, _ := json.Marshal(r.Body)

	rs := bytes.NewReader(jsnm)

	fmt.Println(string(jsnm))
	resp, err := http.Post("https://onviz.bitrix24.site/", "Hello", rs)
	if err != nil {
		fmt.Println("ERROR POST")
	}
	fmt.Println(resp.Body)*/

}

func searchInArray(docs []string, term string) []string {
	var r []string
	for _, doc := range docs {
		if strings.Contains(doc, term) {
			r = append(r, doc)
			break
		}
	}
	return r
}

func searchString(fullText string, keyword string) bool {
	makeArrayFromText := strings.Fields(fullText)
	for _, v := range makeArrayFromText {
		if strings.EqualFold(v, keyword) {
			return true
		}
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

var IndexSites = map[string]int{
	"https://onviz.bitrix24.site/":                                                      1,
	"https://onviz.bitrix24.site/rulonnye/instruktsii/":                                 2,
	"https://onviz.bitrix24.site/rulonnye/rimskie/instruktsii/":                         3,
	"https://onviz.bitrix24.site/lift/instruktsii/":                                     4,
	"https://onviz.bitrix24.site/pergola/instruktsiy/":                                  5,
	"https://onviz.bitrix24.site/pergola/uni2/uni2instruktsiy/":                         6,
	"https://onviz.bitrix24.site/vertikalka/instruktsiivertikalnyy/":                    7,
	"https://onviz.bitrix24.site/instruktsii/":                                          8,
	"https://onviz.bitrix24.site/rulonnye/sborkakarniza38/":                             9,
	"https://onviz.bitrix24.site/rulonnye/sborkakarniza50mm/":                           10,
	"https://onviz.bitrix24.site/rulonnye/zebra/":                                       11,
	"https://onviz.bitrix24.site/rulonnye/nastroykiprivoda/":                            12,
	"https://onviz.bitrix24.site/rulonnye/nastroykimr160/":                              13,
	"https://onviz.bitrix24.site/rulonnye/pustayastranitsa/":                            14,
	"https://onviz.bitrix24.site/rulonnye/smenanapravleniy/":                            15,
	"https://onviz.bitrix24.site/rulonnye/podklyucheniekmihome/":                        16,
	"https://onviz.bitrix24.site/rulonnye/sbrosnastroek/":                               17,
	"https://onviz.bitrix24.site/rulonnye/350obshchienastroyki/":                        18,
	"https://onviz.bitrix24.site/rulonnye/privyazkapulta/":                              19,
	"https://onviz.bitrix24.site/rulonnye/krayniepolozheniya/":                          20,
	"https://onviz.bitrix24.site/rulonnye/smenanapravleniy_asgd/":                       21,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/privyazkapultadlyars/":           22,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/sbrosnastroekdlyars/":            23,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/smenanapravleniydlyars/":         24,
	"https://onviz.bitrix24.site/rulonnye/nastroykaprivoda50/":                          25,
	"https://onviz.bitrix24.site/rulonnye/wf50dobavitmikhom/":                           26,
	"https://onviz.bitrix24.site/rulonnye/smenanapravlenii/":                            27,
	"https://onviz.bitrix24.site/rulonnye/sbros/":                                       28,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/nastroykaprivodov/":              29,
	"https://onviz.bitrix24.site/rulonnye/smenanapravleniy_baew/":                       30,
	"https://onviz.bitrix24.site/rulonnye/sbrosakb50/":                                  31,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/faznayaklavisha_jzmy/":           32,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/sensornayaknopka/":               33,
	"https://onviz.bitrix24.site/rulonnye/privyazkapulta350fc/":                         34,
	"https://onviz.bitrix24.site/rulonnye/sbrosnastroek350fc/":                          35,
	"https://onviz.bitrix24.site/rulonnye/smenanapravleniymr350fc/":                     36,
	"https://onviz.bitrix24.site/rulonnye/klavisha350fc/":                               37,
	"https://onviz.bitrix24.site/rulonnye/nastroykapultov/":                             38,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/poshagovyy/":                     39,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/broadlink/":                      40,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/radiorele/":                      41,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/gruppovoerele/":                  42,
	"https://onviz.bitrix24.site/rulonnye/rezhimzatemneniya/":                           43,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/krayniepolozheniya/":             44,
	"https://onviz.bitrix24.site/rulonnye/rimskie/sborkarim38/":                         45,
	"https://onviz.bitrix24.site/rulonnye/rimskie/sborka50/":                            46,
	"https://onviz.bitrix24.site/rulonnye/nastro":                                       47,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/alisaimikhom/":                   48,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/alisaibroadlink/":                49,
	"https://onviz.bitrix24.site/rulonnye/rimskie/smenastoronyprivoda/":                 50,
	"https://onviz.bitrix24.site/rulonnye/zaputalasnit/":                                51,
	"https://onviz.bitrix24.site/pergola/uni2/skripitrimskiy/":                          52,
	"https://onviz.bitrix24.site/lift/sborkaliftsistemy/":                               53,
	"https://onviz.bitrix24.site/neisp":                                                 54,
	"https://onviz.bitrix24.site/pergola/sborkapergoly/":                                55,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/pustayastranitsa_rssw/":          56,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/pustayastranitsa_ayzl/":          57,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/smenanapravleniy/":               58,
	"https://onviz.bitrix24.site/neispravnostirazdvizh/faznayaklavisha/":                59,
	"https://onviz.bitrix24.site/pergola/uni2/sborkauni2/":                              60,
	"https://onviz.bitrix24.site/privyazkapulta/":                                       61,
	"https://onviz.bitrix24.site/pergola/uni2/smenanapravleniyauni/":                    62,
	"https://onviz.bitrix24.site/pergola/uni2/sbrosnastroek/":                           63,
	"https://onviz.bitrix24.site/pergola/uni2/nastroykapultovuni/":                      64,
	"https://onviz.bitrix24.site/vertikalka/sborkaobshchiy/":                            65,
	"https://onviz.bitrix24.site/vertikalka/sborkavertikalnykhzhalyuzi/":                66,
	"https://onviz.bitrix24.site/vertikalka/montazhelektrokarniza/":                     67,
	"https://onviz.bitrix24.site/vertikalka/obshchienastroykidlyavertikalnykhzhalyuzi/": 68,
	"https://onviz.bitrix24.site/vertikalka/privyazkapultak2234/":                       69,
	"https://onviz.bitrix24.site/vertikalka/sbosnastroekvertikalnyezhalyuzi/":           70,
	"https://onviz.bitrix24.site/vertikalka/smenanapravleniivertik/":                    71,
	"https://onviz.bitrix24.site/vertikalka/mihomevertikalnyezhalyuzi/":                 72,
	"https://onviz.bitrix24.site/vertikalka/regulirovkapovorotalameley/":                73,
}

/*func VisitLastHrefs(lA []string) (text string, urls []string, err error) {
	c := colly.NewCollector()

	for _, vUrl := range lA {

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			urls = append(urls, absoluteURL)
			//fmt.Println(urls) //array with all HREFS
		})

		err = c.Visit(vUrl)
		if err != nil {
			log.Println(err.Error())
			return "", nil, err
		}
		return text, urls, err
	}
	urls = make([]string, 0)
	return "", urls, err
}*/

/*func VisitNextHrefs() {
	c := colly.NewCollector()

	for _, v := range ArrayWithoutDuplicates {

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
			//		fmt.Println("error visit URL")
			continue
		}

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {

			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			ArrayWithoutDuplicates = append(ArrayWithoutDuplicates, absoluteURL)
		})

	}

}*/

/*func LastMain(){
	var fullData string
	var count = 1
	firstSite := checkForValue(count, urls.IndexSites)
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//	fmt.Println(e.Text)
		//fmt.Println(e.Text, "VisitMainPAGE")
		absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

		fullData = e.Text
		//fmt.Println(fullData)

		founded := searchString(fullData, "Сборка")
		if founded == true {
			fmt.Println(fullData, "FULLDATA")
			fmt.Println(absoluteURL)
		} else {
			count = count + 1
			//firstSite
			//fmt.Println("no keyword here")
		}

	})

	err := c.Visit(firstSite)
	if err != nil {
		log.Println(err.Error())
	}

	//HERE is Scan of all hrefs
	for v, _ := range urls.IndexSites {
		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			//fmt.Println(e.Text, "VisitMainPAGE")
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			fullData = e.Text
			//fmt.Println(fullData)

			founded := searchString(fullData, "Сборка")
			if founded == true {
				fmt.Println(fullData, "FULLDATA")
				fmt.Println(absoluteURL)
			} else {
				fmt.Println("no keyword here")
			}

		})

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
	}

	//HERE all hrefs in Massive
	newArr := make([]string, 0)
	for v, _ := range urls.IndexSites {

		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//	fmt.Println(e.Text)
			//fmt.Println(e.Text, "VisitMainPAGE")
			absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))

			if strings.HasPrefix(absoluteURL, `https://onviz`) {
				newArr = append(newArr, absoluteURL)
			}

			//fmt.Println(absoluteURL)
			//fmt.Println(urls) //array with all HREFS
		})

		//withoutDupls := removeDuplicateStr(newArr)

		founded := search(newArr, "alisa")

		//foundedWithoutDuplicates := removeDuplicateStr(founded)

		//newF := removeDuplicateStr(founded)

		fmt.Println(founded)

		/*
			for _, v := range withoutDupls {
				fmt.Println(v)
			}

		err := c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
	}

	fmt.Println(len(newArr))
}
*/

/*func OldMain(){
	Input = "Название заголовка" //пишем своё
	Href = "ID ссылки задания"   //добавляем в нужное место
	p := new(repository.Postgres)
	p.GetOrderByHref(Href) //получем инфу по вставленной ссылке (title,descr,value...)

	repository.PostgresInit(UrlDb)
	//p := new(repository.Postgres)

	page := utils.GenPages(Url)
	checkDuplicate := make([]string, 0)
	//hrefs := scraper.UpdateHrefs(Url)
	//firstHref := hrefs[0]
	scraper.GenHrefs(Url)
	for _, v := range page {
		hrefs := scraper.GenHrefs(v)
		for _, v := range hrefs {
			checkDuplicate = append(checkDuplicate, v)
		}
	}

	var ordersFromHabr string
	var orders = make([]string, 0) //there our orders in hybrid json format todo edit correct formats, example to share (1 - array, 2 - json)
	fmt.Println(orders)
	//ords := make([]scraper.Habr, 0)

	//fmt.Println(ords)

	pages := utils.GenPages(Url)
	var count int
	for _, v := range pages {
		count++
		habr, err := scraper.HabrScraper(v, Input)
		if err != nil {
			log.Println("err scrap on page: ", count)
		}
		var m scraper.Habr
		ja, _ := json.Marshal(habr)
		json.Unmarshal(ja, &m)
		ordersFromHabr = string(ja)
		orders = append(orders, ordersFromHabr)
		//ords, err = parseHabr(ja)
		//fmt.Println(orders)
	}

	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}
	newbot := telegram.NewBot(bot)

	if err = newbot.StartBot(orders); err != nil {
		log.Fatal(err)
	}
}*/

/*var golang = []string{
	"go", "Golang", "golang",
}
*/
//fmt.Println(scraper.NewMassive)
//convert,_ := strconv.Atoi("2609257031")

/*nums := []string{
"	2609257031",
	2609257031
}*/

/*habr2, err := scraper.HabrScraper(Url2, "Go")
if err != nil {
	log.Prefix()
	return
}
var mm scraper.Habr
ja2, _ := json.Marshal(habr2)
json.Unmarshal(ja2, &mm)
fmt.Println(string(ja2))*/

/*func oldMain(){
	err := GetOrdersByB24.VisitHrefs(Url)

	if err != nil {
		fmt.Println(err.Error())
		//fmt.Println("whats going on?")
		return
	}
	//fmt.Println("FULL ARRAY", GetOrdersByB24.FullArray)
	withoutDuplicates := removeDuplicateStr(GetOrdersByB24.FullArray)
	for _, v := range withoutDuplicates {
		if strings.HasPrefix(v, "https://onviz.bitrix24") {
			ArrayWithoutDuplicates = append(ArrayWithoutDuplicates, v)
			bigFullArray = append(bigFullArray, v)
		}
	}

	//	for _, v := range ArrayWithoutDuplicates {
	//	fmt.Println(v)
	//}

	//fmt.Println(ArrayWithoutDuplicates)

	VisitNextHrefs()

	rmvDupls := removeDuplicateStr(ArrayWithoutDuplicates)
	//urlsNew := make([]string, 0)
	for _, v := range rmvDupls {

		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			//absoluteURL := e.Request.AbsoluteURL(e.Attr("href"))
			//urlsNew = append(urlsNew, absoluteURL)
			if strings.HasPrefix(v, "https://onviz.bitrix24") {
				//urlsNew = append(urlsNew, v)
				bigFullArray = append(bigFullArray, v)

			}

		})

		err = c.Visit(v)
		if err != nil {
			log.Println(err.Error())
		}
		//fmt.Println(urlsNew)
	}

	//newRemovedDupls := removeDuplicateStr(urlsNew)

	//for _, v := range newRemovedDupls {
	//fmt.Println(v)
	//}

	removed := removeDuplicateStr(bigFullArray)
	for _, v := range removed {

		lastestArray = append(lastestArray, v)

		//fmt.Println(v)
	}
	for _, v := range lastestArray {
		fmt.Println(v)
	}

	fmt.Println(len(lastestArray))

	fmt.Println("end")
}*/
