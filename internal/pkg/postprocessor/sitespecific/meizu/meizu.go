package meizu

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/internetarchive/Zeno/pkg/models"
)

func IsURL(URL *models.URL) bool {
	host := URL.GetParsed().Host
	return strings.Contains(host, "meizu.cn") || strings.Contains(host, "meizu.com")
}

func IsThreadURL(URL *models.URL) bool {
	host := URL.GetParsed().Host
	path := URL.GetParsed().Path
	return strings.Contains(host, "meizu.cn") && strings.HasPrefix(path, "/thread/")
}

func getThreadID(URL *models.URL) string {
	path := URL.GetParsed().Path
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

// https://www.meizu.cn/thread/23149497

// static config:
// https://myplus-api.meizu.cn/myplus-qing/ug/app/start
// https://myplus-api.meizu.cn/myplus-qing/g/config/app
//
// thread content:
// https://myplus-api.meizu.cn/myplus-qing/ug/content/23149497
//
// comments: (iterate?)
// https://myplus-api.meizu.cn/myplus-qing/ug/comment/list?contentId=23149497&currentPage=0&secondary=1&sortBy=0
//
// 热门话题 sidebar: (nessary?)
// https://myplus-api.meizu.cn/myplus-qing/ug/topic/list/hot/v2?pageSize=10
//
// 热门圈子 sidebar: (nessary?)
// https://myplus-api.meizu.cn/myplus-qing/ug/forum/hotForum
func AddAssetsFromThreadURL(URL *models.URL) ([]*models.URL, error) {
	threadID := getThreadID(URL)
	if threadID == "" {
		return nil, fmt.Errorf("unable to extract thread ID from URL: %s", URL.Raw)
	}
	assets := []*models.URL{}
	assets = append(assets, &models.URL{Raw: "https://myplus-api.meizu.cn/myplus-qing/ug/app/start"})
	assets = append(assets, &models.URL{Raw: "https://myplus-api.meizu.cn/myplus-qing/g/config/app"})

	return assets, nil
}

// https://myplus-api.meizu.cn/myplus-qing/ug/content/+
func IsContentURL(URL *models.URL) bool {
	host := URL.GetParsed().Host
	path := URL.GetParsed().Path
	return strings.Contains(host, "meizu.cn") && strings.HasPrefix(path, "/myplus-qing/ug/content/")
}

// https://myplus-api.meizu.cn/myplus-qing/ug/comment/list
func IsCommentListURL(URL *models.URL) bool {
	host := URL.GetParsed().Host
	path := URL.GetParsed().Path
	return strings.Contains(host, "meizu.cn") && strings.HasPrefix(path, "/myplus-qing/ug/comment/list")
}

// .data.posts.hasMore
func CommentListHasMoreComments(URL *models.URL) bool {
	defer URL.RewindBody()
	dec := json.NewDecoder(URL.GetBody())
	var data struct {
		Data struct {
			Posts struct {
				HasMore bool `json:"hasMore"`
			} `json:"posts"`
		} `json:"data"`
	}

	if err := dec.Decode(&data); err != nil {
		return false
	}

	return data.Data.Posts.HasMore
}

func CommentCurrentPage(URL *models.URL) int {
	url := URL.GetParsed()
	query := url.Query()
	currentPage := query.Get("currentPage")
	if currentPage == "" {
		return 0
	}

	pageNum, err := strconv.Atoi(currentPage)
	if err != nil {
		panic(err)
	}

	return pageNum
}

// download.meizu.com
