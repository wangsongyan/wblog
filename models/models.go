package models

import (
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

// I don't need soft delete,so I use customized BaseModel instead gorm.Model
type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// table pages
type Page struct {
	BaseModel
	Body        string // body
	View        int    // view count
	IsPublished string // published or not
}

// table posts
type Post struct {
	Page
	Title string // title
	Tags  []*Tag `gorm:"-"` // tags of post
}

// table tags
type Tag struct {
	BaseModel
	Name  string // tag name
	Total int    `gorm:"-"` // count of post
}

// table post_tags
type PostTag struct {
	BaseModel
	PostId uint // post id
	TagId  uint // tag id
}

var DB *gorm.DB

func InitDB() *gorm.DB {
	//db, err := gorm.Open("sqlite3", "wblog.db")
	db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB = db

	db.AutoMigrate(&Page{}, &Post{}, &Tag{}, &PostTag{})
	db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")

	return db
}

// Page
func (page *Page) Insert() error {
	return DB.Create(page).Error
}

func (page *Page) Update() error {
	return DB.Model(page).Update(Page{Body: page.Body, IsPublished: page.IsPublished}).Error
}

func (page *Page) Delete() error {
	return DB.Delete(page).Error
}

func GetPageById(id string) (*Page, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var page Page
	err = DB.First(&page, "id = ?", pid).Error
	return &page, err
}

func ListPage() ([]*Page, error) {
	var pages []*Page
	err := DB.First(pages).Error
	return pages, err
}

// Post
func (post *Post) Insert() error {
	return DB.Create(post).Error
}

func (post *Post) Update() error {
	p := Post{Title: post.Title}
	p.Body = post.Body
	return DB.Model(post).Update(p).Error
}

func (post *Post) Delete() error {
	return DB.Delete(post).Error
}

func ListPost(tag string) ([]*Post, error) {
	var posts []*Post
	var err error
	if len(tag) > 0 {
		tagId, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return nil, err
		}
		rows, err := DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ?", tagId).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post Post
			rows.Scan(&post)
			posts = append(posts, &post)
		}
	} else {
		err = DB.First(posts).Error
	}
	return posts, err
}

func GetPostById(id string) (*Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var post Post
	err = DB.First(&post, "id = ?", pid).Error
	return &post, err
}

// Tag
func (tag *Tag) Insert() error {
	return DB.FirstOrCreate(tag, "name = ?", tag.Name).Error
}

func ListTag() ([]*Tag, error) {
	var tags []*Tag
	rows, err := DB.Raw("select t.*,count(*) total from tags t inner join post_tags pt on t.id = pt.tag_id group by pt.tag_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

func ListTagByPostId(id string) ([]*Tag, error) {
	var tags []*Tag
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	rows, err := DB.Raw("select t.* from tags t inner join post_tags pt on t.id = pt.tag_id where pt.post_id = ?", uint(pid)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}
