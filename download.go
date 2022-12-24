package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
)

var (
	title = flag.String("title", "", "title of the epub")
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

type Data struct {
	AllAwardings []struct {
		AwardSubType                     string           `json:"award_sub_type"`
		AwardType                        string           `json:"award_type"`
		AwardingsRequiredToGrantBenefits *int             `json:"awardings_required_to_grant_benefits"`
		CoinPrice                        int              `json:"coin_price"`
		CoinReward                       int              `json:"coin_reward"`
		Count                            int              `json:"count"`
		DaysOfPremium                    *int             `json:"days_of_premium"`
		Description                      string           `json:"description"`
		ID                               string           `json:"id"`
		IconFormat                       *string          `json:"icon_format"`
		IconHeight                       int              `json:"icon_height"`
		IconURL                          string           `json:"icon_url"`
		IconWidth                        int              `json:"icon_width"`
		IsEnabled                        bool             `json:"is_enabled"`
		IsNew                            bool             `json:"is_new"`
		Name                             string           `json:"name"`
		PennyPrice                       *int             `json:"penny_price"`
		ResizedIcons                     []HeightURLWidth `json:"resized_icons"`
		ResizedStaticIcons               []HeightURLWidth `json:"resized_static_icons"`
		StaticIconHeight                 int              `json:"static_icon_height"`
		StaticIconURL                    string           `json:"static_icon_url"`
		StaticIconWidth                  int              `json:"static_icon_width"`
		SubredditCoinReward              int              `json:"subreddit_coin_reward"`
		TiersByRequiredAwardings         *struct {
			Key0 Key `json:"0"`
			Key3 Key `json:"3"`
			Key6 Key `json:"6"`
			Key9 Key `json:"9"`
		} `json:"tiers_by_required_awardings"`
	} `json:"all_awardings"`
	AllowLiveComments          *bool       `json:"allow_live_comments,omitempty"`
	Archived                   bool        `json:"archived"`
	Author                     string      `json:"author"`
	AuthorFlairBackgroundColor *string     `json:"author_flair_background_color"`
	AuthorFlairCSSClass        *string     `json:"author_flair_css_class"`
	AuthorFlairTemplateID      *string     `json:"author_flair_template_id"`
	AuthorFlairText            *string     `json:"author_flair_text"`
	AuthorFlairTextColor       *string     `json:"author_flair_text_color"`
	AuthorFlairType            *string     `json:"author_flair_type,omitempty"`
	AuthorFullname             *string     `json:"author_fullname,omitempty"`
	AuthorIsBlocked            bool        `json:"author_is_blocked"`
	AuthorPatreonFlair         *bool       `json:"author_patreon_flair,omitempty"`
	AuthorPremium              *bool       `json:"author_premium,omitempty"`
	Body                       *string     `json:"body,omitempty"`
	BodyHTML                   *string     `json:"body_html,omitempty"`
	CanGild                    bool        `json:"can_gild"`
	CanModPost                 bool        `json:"can_mod_post"`
	Clicked                    *bool       `json:"clicked,omitempty"`
	Collapsed                  *bool       `json:"collapsed,omitempty"`
	ContentCategories          []string    `json:"content_categories,omitempty"`
	ContestMode                *bool       `json:"contest_mode,omitempty"`
	Controversiality           *int        `json:"controversiality,omitempty"`
	Created                    float64     `json:"created"`
	CreatedUtc                 float64     `json:"created_utc"`
	Depth                      *int        `json:"depth,omitempty"`
	Domain                     *string     `json:"domain,omitempty"`
	Downs                      int         `json:"downs"`
	Edited                     interface{} `json:"edited"`
	Gilded                     int         `json:"gilded"`
	Gildings                   struct {
		Gid1 *int `json:"gid_1,omitempty"`
		Gid2 *int `json:"gid_2,omitempty"`
	} `json:"gildings"`
	Hidden                   *bool   `json:"hidden,omitempty"`
	HideScore                *bool   `json:"hide_score,omitempty"`
	ID                       string  `json:"id"`
	IsCreatedFromAdsUI       *bool   `json:"is_created_from_ads_ui,omitempty"`
	IsCrosspostable          *bool   `json:"is_crosspostable,omitempty"`
	IsMeta                   *bool   `json:"is_meta,omitempty"`
	IsOriginalContent        *bool   `json:"is_original_content,omitempty"`
	IsRedditMediaDomain      *bool   `json:"is_reddit_media_domain,omitempty"`
	IsRobotIndexable         *bool   `json:"is_robot_indexable,omitempty"`
	IsSelf                   *bool   `json:"is_self,omitempty"`
	IsSubmitter              *bool   `json:"is_submitter,omitempty"`
	IsVideo                  *bool   `json:"is_video,omitempty"`
	LinkFlairBackgroundColor *string `json:"link_flair_background_color,omitempty"`
	LinkFlairCSSClass        *string `json:"link_flair_css_class,omitempty"`
	LinkFlairTemplateID      *string `json:"link_flair_template_id,omitempty"`
	LinkFlairText            *string `json:"link_flair_text,omitempty"`
	LinkFlairTextColor       *string `json:"link_flair_text_color,omitempty"`
	LinkFlairType            *string `json:"link_flair_type,omitempty"`
	LinkID                   *string `json:"link_id,omitempty"`
	Locked                   bool    `json:"locked"`
	MediaEmbed               *struct {
	} `json:"media_embed,omitempty"`
	MediaOnly             *bool       `json:"media_only,omitempty"`
	Name                  string      `json:"name"`
	NoFollow              bool        `json:"no_follow"`
	NumComments           *int        `json:"num_comments,omitempty"`
	NumCrossposts         *int        `json:"num_crossposts,omitempty"`
	NumDuplicates         *int        `json:"num_duplicates,omitempty"`
	Over18                *bool       `json:"over_18,omitempty"`
	ParentID              *string     `json:"parent_id,omitempty"`
	ParentWhitelistStatus *string     `json:"parent_whitelist_status,omitempty"`
	Permalink             string      `json:"permalink"`
	Pinned                *bool       `json:"pinned,omitempty"`
	Pwls                  *int        `json:"pwls,omitempty"`
	Quarantine            *bool       `json:"quarantine,omitempty"`
	Replies               interface{} `json:"replies,omitempty"`
	Saved                 bool        `json:"saved"`
	Score                 int         `json:"score"`
	ScoreHidden           *bool       `json:"score_hidden,omitempty"`
	SecureMediaEmbed      *struct {
	} `json:"secure_media_embed,omitempty"`
	Selftext              *string  `json:"selftext,omitempty"`
	SelftextHTML          *string  `json:"selftext_html,omitempty"`
	SendReplies           bool     `json:"send_replies"`
	Spoiler               *bool    `json:"spoiler,omitempty"`
	Stickied              bool     `json:"stickied"`
	Subreddit             string   `json:"subreddit"`
	SubredditID           string   `json:"subreddit_id"`
	SubredditNamePrefixed string   `json:"subreddit_name_prefixed"`
	SubredditSubscribers  *int     `json:"subreddit_subscribers,omitempty"`
	SubredditType         string   `json:"subreddit_type"`
	SuggestedSort         *string  `json:"suggested_sort,omitempty"`
	Thumbnail             *string  `json:"thumbnail,omitempty"`
	Title                 *string  `json:"title,omitempty"`
	TotalAwardsReceived   int      `json:"total_awards_received"`
	URL                   *string  `json:"url,omitempty"`
	Ups                   int      `json:"ups"`
	UpvoteRatio           *float64 `json:"upvote_ratio,omitempty"`
	Visited               *bool    `json:"visited,omitempty"`
	WhitelistStatus       *string  `json:"whitelist_status,omitempty"`
	Wls                   *int     `json:"wls,omitempty"`
}

type Document []struct {
	Data struct {
		Children []struct {
			Data Data   `json:"data"`
			Kind string `json:"kind"`
		} `json:"children"`
		Dist      *int   `json:"dist"`
		GeoFilter string `json:"geo_filter"`
		Modhash   string `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"`
}
type Icon struct {
	Format string `json:"format"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type HeightURLWidth struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type Key struct {
	AwardingsRequired  int              `json:"awardings_required"`
	Icon               Icon             `json:"icon"`
	ResizedIcons       []HeightURLWidth `json:"resized_icons"`
	ResizedStaticIcons []HeightURLWidth `json:"resized_static_icons"`
	StaticIcon         HeightURLWidth   `json:"static_icon"`
}

func getURL(url string) ([]byte, error) {
	url = strings.ReplaceAll(url, "www.reddit.com", "old.reddit.com")
	url += ".json"

	hash := sha256.Sum256([]byte(url))
	hashStr := hex.EncodeToString(hash[:])
	cacheFile := filepath.Join("/tmp/", "redditepub-"+hashStr+".json")
	cacheBody, err := ioutil.ReadFile(cacheFile)
	if err == nil {
		return cacheBody, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	log.Printf("fetching %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Reddit Epub/0.1")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(cacheFile, body, 0600); err != nil {
		return nil, err
	}
	return body, nil
}

func fetchBody(url string) (*Document, error) {
	body, err := getURL(url)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func run() error {
	flag.Parse()
	if *title == "" {
		return errors.Errorf("must set title")
	}

	e := epub.NewEpub(*title)
	for _, nextPage := range flag.Args() {
		for nextPage != "" {
			doc, err := fetchBody(nextPage)
			if err != nil {
				return err
			}
			post := (*doc)[0].Data.Children[0].Data
			author := post.Author
			e.SetAuthor(author)

			log.Printf("%s - %s: %s", *post.Title, author, nextPage)
			body := "# " + *post.Title + "\n\n" + *post.Selftext

			for _, comment := range (*doc)[1].Data.Children {
				if comment.Data.Author == author {
					body += "\n\n---\n\n" + *comment.Data.Body
				}
			}

			body += fmt.Sprintf("\n\n[%s](%s)", nextPage, nextPage)

			html := blackfriday.Run([]byte(body))
			e.AddSection(string(html), *post.Title, "", "")

			cdoc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
			if err != nil {
				return err
			}
			nextPage = ""
			cdoc.Find("a[href]").EachWithBreak(func(i int, elem *goquery.Selection) bool {
				text := strings.ToLower(elem.Text())
				if strings.Contains(text, "next") || strings.Contains(text, "forward") {
					href := elem.AttrOr("href", "")
					if strings.Contains(href, "reddit.com") {
						nextPage = href
						return false
					}
				}
				return true
			})
		}
	}

	if _, err := e.WriteTo(os.Stdout); err != nil {
		return err
	}

	return nil
}
