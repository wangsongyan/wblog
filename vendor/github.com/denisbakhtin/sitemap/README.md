XML sitemap
=======

#Usage
```
  import (
    "github.com/denisbakhtin/sitemap"
    "time"
    ...
  )

  func CreateSitemap() {
    folder = "public_sitemap_folder"
    domain := "http://mydomain.com"
    now := time.Now()
    items := make([]sitemap.Item, 1)

    //Home page
    items = append(items, sitemap.Item{
      Loc:        fmt.Sprintf("%s", domain),
      LastMod:    now,
      Changefreq: "daily",
      Priority:   1,
    })

    //pages
    pages := models.GetPublishedPages() //get slice of pages
    for i := range pages {
      items = append(items, sitemap.Item{
        Loc:        fmt.Sprintf("%s/pages/%d", domain, pages[i].Id), //page url
        LastMod:    pages[i].UpdatedAt, //page modification timestamp (time.Time)
        Changefreq: "monthly", //or "weekly", "daily", "hourly", ...
        Priority:   0.8,
      })
    }
    if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"),
  items); err != nil {
      log.Error(err)
      return
    }
    if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml",
  domain+"/public/sitemap/"); err != nil {
      log.Error(err)
      return
    }
    //done
  }

```
For periodical sitemap generation you may use something like this: https://github.com/jasonlvhit/gocron

#TODO
- Image, Video sitemap
- Split large sitemaps by threshold (50.000 items)
- Ping search engines?

Generate gzipped Google sitemap and sitemap index with golang
