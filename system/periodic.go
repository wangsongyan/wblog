package system

import (
	"fmt"
	"github.com/denisbakhtin/sitemap"
	"github.com/wangsongyan/wblog/models"
	"os"
	"path"
	"time"
)

func CreateXMLSitemap() {

	folder := path.Join(GetConfiguration().Public, "sitemap")
	os.MkdirAll(folder, os.ModePerm)
	domain := GetConfiguration().Domain
	now := time.Now()
	items := make([]sitemap.Item, 0)

	items = append(items, sitemap.Item{
		Loc:        domain,
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	posts, err := models.ListPublishedPost("")
	if err == nil {
		for _, post := range posts {
			items = append(items, sitemap.Item{
				Loc:        fmt.Sprintf("%s/post/%d", domain, post.ID),
				LastMod:    post.UpdatedAt,
				Changefreq: "weekly",
				Priority:   0.9,
			})
		}
	}

	pages, err := models.ListPublishedPage()
	if err == nil {
		for _, page := range pages {
			items = append(items, sitemap.Item{
				Loc:        fmt.Sprintf("%s/page/%d", domain, page.ID),
				LastMod:    page.UpdatedAt,
				Changefreq: "monthly",
				Priority:   0.8,
			})
		}
	}

	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/static/sitemap/"); err != nil {
		return
	}
}
