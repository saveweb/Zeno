package zhubai

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func IsURL(url *url.URL) bool {
	return strings.Contains(url.Host, "zhubai.love")
}

func getSubdomain(url *url.URL) string {
	return strings.Split(url.Host, ".")[0]
}

func IsOnHostPage(url *url.URL) bool {
	return url.Path == "/" || url.Path == ""
}

func GenPublicationsInitUrl(url *url.URL) (*url.URL, error) {
	newURL, err := url.Parse(url.Scheme + "://" + url.Host + "/api/publications/" + getSubdomain(url) + "?id_type=token")
	if err != nil {
		return nil, err
	}
	return newURL, nil
}

func GenPublicationsPostInitUrl(url *url.URL) (*url.URL, error) {
	newURL, err := url.Parse(url.Scheme + "://" + url.Host + "/api/publications/" + getSubdomain(url) + "/posts?publication_id_type=token")
	if err != nil {
		return nil, err
	}
	return newURL, nil
}

func GenPostApiUrl(url *url.URL, postID string) (*url.URL, error) {
	newURL, err := url.Parse(url.Scheme + "://" + url.Host + "/api/posts/" + postID)
	if err != nil {
		return nil, err
	}
	return newURL, nil
}

func GenSelfUrl(url *url.URL) (*url.URL, error) {
	// https://news.zhubai.love/api/self
	newURL, err := url.Parse(url.Scheme + "://" + url.Host + "/api/self")
	if err != nil {
		return nil, err
	}
	return newURL, nil
}

func IsPublicationsPostURL(url *url.URL) bool {
	return strings.Contains(url.Path, "/api/publications/") && strings.Contains(url.Path, "/posts")
}

type PostsAPIResponse struct {
	Data []struct {
		ID            string `json:"id"`
		CreatedAt     int    `json:"created_at"`
		IsPaidContent bool   `json:"is_paid_content"`
		Title         string `json:"title"`
	}
	Pagination struct {
		Next string `json:"next"`
	}
}

func ParsePostsAPIResponse(resp *http.Response) (*PostsAPIResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var postsAPIResponse PostsAPIResponse
	err = json.Unmarshal(body, &postsAPIResponse)
	if err != nil {
		return nil, err
	}

	return &postsAPIResponse, nil
}
