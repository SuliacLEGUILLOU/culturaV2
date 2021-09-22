package main

import (
    "context"
    "io/ioutil"
    "net/http"
    "fmt"
    "encoding/json"
    "os"

    "github.com/mattn/go-mastodon"
)

type content map[string]string

type WikiResp struct {
    Type string
    Title string
    Extract string
    Lang string
    ContentUrls map[string]content `json:"content_urls"`
}

/**
 *  Get a page from wikipedia
 *  Given url should link to the random page for more fun
 */
func getWikiPage(url string) (WikiResp, error) {
    resp, err := http.Get(url)
    if err != nil { return WikiResp{}, err }
    if resp.StatusCode != 200 { return WikiResp{}, fmt.Errorf(resp.Status)}

    raw_body, err := ioutil.ReadAll(resp.Body)
    if err != nil { return WikiResp{}, err }

    var page WikiResp

    err = json.Unmarshal(raw_body, &page)
    if err != nil { return WikiResp{}, err }

    return page, nil
}

/**
 *  FormatMessage into a Toot
 */
func formatMessage(page WikiResp) mastodon.Toot {
    return mastodon.Toot{
        Status: fmt.Sprintf("%s\n\n%s", page.Extract, page.ContentUrls["desktop"]["page"]),
        Visibility: "unlisted",
    }
}

func main() {
    wiki_url := os.Getenv("WIKI_URL")
    config := mastodon.Config{
        Server: os.Getenv("INSTANCE_URL"),
        ClientID: os.Getenv("CLIENT_ID"),
        ClientSecret: os.Getenv("CLIENT_SECRET"),
        AccessToken: os.Getenv("ACCESS_TOKEN"),
    }

    client := mastodon.NewClient(&config)
    
    page, err := getWikiPage(wiki_url)
    if err != nil {
        fmt.Println("Error while getting wikipedia page:", err)
        return
    }

    message := formatMessage(page)
    
    _, err = client.PostStatus(context.Background(), &message)
    if err != nil {
        fmt.Println("Error while posting on mastodon:", err)
    }
}