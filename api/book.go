package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gitlab.odds.team/plus1/backend-go/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"golang.org/x/net/html/charset"
)

// GetBookDetail ...
func (db *MongoDB) GetBookDetail(e echo.Context) error {
	url := &model.BookUrl{}
	if err := e.Bind(url); err != nil {
		return err
	}

	book := &model.Book{}
	doc, err := Init(url.BookUrl)
	if err == nil {
		cover, err := CheckCover(doc)
		if err {
			book.Format = cover
			fmt.Println("Is Hardcover or Paperback?: ", cover)
			imageUrl, err := GetUrlImgage(doc)
			if err != nil {
				return err
			}
			book.Imgage = imageUrl
			bookName, err := GetBookName(doc)
			if err != nil {
				return err
			}
			book.BookName = bookName
			book_author, err := GetBookAuthor(doc)
			if err != nil {
				return err
			}
			book.Author = book_author
			price, err := GetPrice(doc)
			if err != nil {
				return err
			}
			book.Price = price
		} else {
			err := &model.FindError{
				Message: "Please select Hardcover or Paperback",
			}

			return e.JSON(404, err)
		}

	}
	return e.JSON(http.StatusOK, book)

}
func GetUrlImgage(doc *goquery.Document) (reulst string, err error) {
	//type aa []string
	// Find the review items
	var img_url string
	var firstTime bool = true
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band, ok := s.Find("img").Attr("data-a-dynamic-image")
		if ok && firstTime {
			firstTime = false
			img_url = band[2:strings.Index(band, "\":")]
			// fmt.Println(i, img_url)
		} else {
			return
		}
	})
	return img_url, nil
}

func GetBookName(doc *goquery.Document) (result string, err error) {
	//type aa []string
	// Find the review items
	var book_name string
	// var firstTime bool = true
	doc.Find("h1#title").Each(func(i int, s *goquery.Selection) {
		result = s.Find("span#ebooksProductTitle").Text()
		if result != "" {
			book_name = result
		} else {
			book_name = s.Find("span#productTitle").Text()
		}
		// fmt.Println(i, band)
		// if i == 0 {
		return
		// }

	})
	return book_name, nil

}

func GetBookAuthor(doc *goquery.Document) (result string, err error) {
	var book_author string = ""
	doc.Find("span.author.notFaded").Each(func(j int, se *goquery.Selection) {
		// fmt.Println(se.Html())
		html, err := se.Html()
		if err != nil {
			return
		}
		if strings.Contains(html[0:100], "<a") {
			if strings.Contains(se.Find("span.contribution").Find("span").Text(), "(Author)") {
				result = se.Find("a").Text()
				// fmt.Println(result)
				if book_author != "" && result != "" {
					book_author = book_author + ", "
				}
				// fmt.Println("In first Loop")
				book_author = book_author + result
				return
			}

		} else {
			// fmt.Println("In second Loop")
			se.Find("span.a-declarative").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(se.Find("span.contribution").Find("span").Text(), "(Author)") {
					result = s.Find("a").Text()
					// fmt.Println(result)
					if book_author != "" && result != "" {
						book_author = book_author + ", "
					}
					// fmt.Println(i, result)
					book_author = book_author + result
					return
				}
			})
		}

	})

	return book_author, nil
}

func GetPrice(doc *goquery.Document) (result string, err error) {
	var price string
	var temp string
	// var isFirst bool = true
	// fmt.Println(doc.Find(".a-color-price").
	re := regexp.MustCompile(`\d[\d,]*[\.]?[\d]*`)
	result, err = doc.Find("div#buybox").Html()
	if result != "" {
		temp = doc.Find("div#buybox").Find("span.a-color-price").Text()
		price = re.FindString(temp)
	} else {
		temp = doc.Find("li.a-tab-heading.a-active.mediaTab_heading").Find("span.a-size-base.mediaTab_subtitle").Text()
		price = re.FindString(temp)
	}
	// doc.Find("div#buy-box").Each(func(i int, s *goquery.Selection) {
	// 	if isFirst {
	// 		isFirst = false
	// 		price = strings.TrimSpace(s.Text())
	// 		return
	// 	}
	// })
	// fmt.Println(temp)
	return price, nil
}

func CheckCover(doc *goquery.Document) (result string, err bool) {
	var cover string
	doc.Find("table#productDetailsTable .bucket .content ").Each(func(i int, s *goquery.Selection) {
		s.Find("b").Each(func(j int, se *goquery.Selection) {
			fmt.Println(se.Html())
			if strings.Contains(se.Text(), "Hardcover") || strings.Contains(se.Text(), "Paperback") {
				cover = se.Text()
			} else {
				return
			}
		})
		return
	})
	if cover != "" {
		fmt.Println("return true")
		return cover[:len(cover)-1], true
	}
	return "", false
}

func Init(url string) (*goquery.Document, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := req.Header.Get("Content-Type") // Optional, better guessing
	utf8reader, err := charset.NewReader(resp.Body, contentType)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utf8reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
