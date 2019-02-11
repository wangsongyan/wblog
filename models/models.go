package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/wangsongyan/wblog/system"
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
	Title       string // title
	Body        string // body
	View        int    // view count
	IsPublished bool   // published or not
}

// table posts
type Post struct {
	BaseModel
	Title        string     // title
	Body         string     // body
	View         int        // view count
	IsPublished  bool       // published or not
	Tags         []*Tag     `gorm:"-"` // tags of post
	Comments     []*Comment `gorm:"-"` // comments of post
	CommentTotal int        `gorm:"-"` // count of comment
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

// table users
type User struct {
	gorm.Model
	Email         string    `gorm:"unique_index;default:null"` //邮箱
	Telephone     string    `gorm:"unique_index;default:null"` //手机号码
	Password      string    `gorm:"default:null"`              //密码
	VerifyState   string    `gorm:"default:'0'"`               //邮箱验证状态
	SecretKey     string    `gorm:"default:null"`              //密钥
	OutTime       time.Time //过期时间
	GithubLoginId string    `gorm:"unique_index;default:null"` // github唯一标识
	GithubUrl     string    //github地址
	IsAdmin       bool      //是否是管理员
	AvatarUrl     string    // 头像链接
	NickName      string    // 昵称
	LockState     bool      `gorm:"default:'0'"` //锁定状态
}

// table comments
type Comment struct {
	BaseModel
	UserID    uint   // 用户id
	Content   string // 内容
	PostID    uint   // 文章id
	ReadState bool   `gorm:"default:'0'"` // 阅读状态
	//Replies []*Comment // 评论
	NickName  string `gorm:"-"`
	AvatarUrl string `gorm:"-"`
	GithubUrl string `gorm:"-"`
}

// table subscribe
type Subscriber struct {
	gorm.Model
	Email          string    `gorm:"unique_index"` //邮箱
	VerifyState    bool      `gorm:"default:'0'"`  //验证状态
	SubscribeState bool      `gorm:"default:'1'"`  //订阅状态
	OutTime        time.Time //过期时间
	SecretKey      string    // 秘钥
	Signature      string    //签名
}

// table link
type Link struct {
	gorm.Model
	Name string //名称
	Url  string //地址
	Sort int    `gorm:"default:'0'"` //排序
	View int    //访问次数
}

// query result
type QrArchive struct {
	ArchiveDate time.Time //month
	Total       int       //total
	Year        int       // year
	Month       int       // month
}

type SmmsFile struct {
	BaseModel
	FileName  string `json:"filename"`
	StoreName string `json:"storename"`
	Size      int    `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Hash      string `json:"hash"`
	Delete    string `json:"delete"`
	Url       string `json:"url"`
	Path      string `json:"path"`
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {

	db, err := gorm.Open("sqlite3", system.GetConfiguration().DSN)
	//db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Asia/Shanghai")
	if err == nil {
		DB = db
		//db.LogMode(true)
		db.AutoMigrate(&Page{}, &Post{}, &Tag{}, &PostTag{}, &User{}, &Comment{}, &Subscriber{}, &Link{}, &SmmsFile{})
		db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
		return db, err
	}
	return nil, err
}

// Page
func (page *Page) Insert() error {
	return DB.Create(page).Error
}

func (page *Page) Update() error {
	return DB.Model(page).Updates(map[string]interface{}{
		"title":        page.Title,
		"body":         page.Body,
		"is_published": page.IsPublished,
	}).Error
}

func (page *Page) UpdateView() error {
	return DB.Model(page).Updates(map[string]interface{}{
		"view": page.View,
	}).Error
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

func ListPublishedPage() ([]*Page, error) {
	return _listPage(true)
}

func ListAllPage() ([]*Page, error) {
	return _listPage(false)
}

func _listPage(published bool) ([]*Page, error) {
	var pages []*Page
	var err error
	if published {
		err = DB.Where("is_published = ?", true).Find(&pages).Error
	} else {
		err = DB.Find(&pages).Error
	}
	return pages, err
}

func CountPage() int {
	var count int
	DB.Model(&Page{}).Count(&count)
	return count
}

// Post
func (post *Post) Insert() error {
	return DB.Create(post).Error
}

func (post *Post) Update() error {
	return DB.Model(post).Updates(map[string]interface{}{
		"title":        post.Title,
		"body":         post.Body,
		"is_published": post.IsPublished,
	}).Error
}

func (post *Post) UpdateView() error {
	return DB.Model(post).Updates(map[string]interface{}{
		"view": post.View,
	}).Error
}

func (post *Post) Delete() error {
	return DB.Delete(post).Error
}

func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body))))
	runes := []rune(sanitized)
	if len(runes) > 300 {
		sanitized = string(runes[:300])
	}
	excerpt := template.HTML(sanitized + "...")
	return excerpt
}

func ListPublishedPost(tag string, pageIndex, pageSize int) ([]*Post, error) {
	return _listPost(tag, true, pageIndex, pageSize)
}

func ListAllPost(tag string) ([]*Post, error) {
	return _listPost(tag, false, 0, 0)
}

func _listPost(tag string, published bool, pageIndex, pageSize int) ([]*Post, error) {
	var posts []*Post
	var err error
	if len(tag) > 0 {
		tagId, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return nil, err
		}
		var rows *sql.Rows
		if published {
			if pageIndex > 0 {
				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc limit ? offset ?", tagId, true, pageSize, (pageIndex-1)*pageSize).Rows()
			} else {
				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc", tagId, true).Rows()
			}
		} else {
			rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? order by created_at desc", tagId).Rows()
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post Post
			DB.ScanRows(rows, &post)
			posts = append(posts, &post)
		}
	} else {
		if published {
			if pageIndex > 0 {
				err = DB.Where("is_published = ?", true).Order("created_at desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&posts).Error
			} else {
				err = DB.Where("is_published = ?", true).Order("created_at desc").Find(&posts).Error
			}
		} else {
			err = DB.Order("created_at desc").Find(&posts).Error
		}
	}
	return posts, err
}

func MustListMaxReadPost() (posts []*Post) {
	posts, _ = ListMaxReadPost()
	return
}

func ListMaxReadPost() (posts []*Post, err error) {
	err = DB.Where("is_published = ?", true).Order("view desc").Limit(5).Find(&posts).Error
	return
}

func MustListMaxCommentPost() (posts []*Post) {
	posts, _ = ListMaxCommentPost()
	return
}

func ListMaxCommentPost() (posts []*Post, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = DB.Raw("select p.*,c.total comment_total from posts p inner join (select post_id,count(*) total from comments group by post_id) c on p.id = c.post_id order by c.total desc limit 5").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		DB.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return
}

func CountPostByTag(tag string) (count int, err error) {
	var (
		tagId uint64
	)
	if len(tag) > 0 {
		tagId, err = strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return
		}
		err = DB.Raw("select count(*) from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ?", tagId, true).Row().Scan(&count)
	} else {
		err = DB.Raw("select count(*) from posts p where p.is_published = ?", true).Row().Scan(&count)
	}
	return
}

func CountPost() int {
	var count int
	DB.Model(&Post{}).Count(&count)
	return count
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

func MustListPostArchives() []*QrArchive {
	archives, _ := ListPostArchives()
	return archives
}

func ListPostArchives() ([]*QrArchive, error) {
	var archives []*QrArchive
	//querysql := `select DATE_FORMAT(created_at,'%Y-%m') as month,count(*) as total from posts where is_published = ? group by month order by month desc`
	querysql := `select strftime('%Y-%m',created_at) as month,count(*) as total from posts where is_published = ? group by month order by month desc`
	rows, err := DB.Raw(querysql, true).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var archive QrArchive
		var month string
		rows.Scan(&month, &archive.Total)
		//DB.ScanRows(rows, &archive)
		archive.ArchiveDate, _ = time.Parse("2006-01", month)
		archive.Year = archive.ArchiveDate.Year()
		archive.Month = int(archive.ArchiveDate.Month())
		archives = append(archives, &archive)
	}
	return archives, nil
}

func ListPostByArchive(year, month string, pageIndex, pageSize int) ([]*Post, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if len(month) == 1 {
		month = "0" + month
	}
	condition := fmt.Sprintf("%s-%s", year, month)
	if pageIndex > 0 {
		//querysql := `select * from posts where date_format(created_at,'%Y-%m') = ? and is_published = ? order by created_at desc limit ? offset ?`
		querysql := `select * from posts where strftime('%Y-%m',created_at) = ? and is_published = ? order by created_at desc limit ? offset ?`
		rows, err = DB.Raw(querysql, condition, true, pageSize, (pageIndex-1)*pageSize).Rows()
	} else {
		//querysql := `select * from posts where date_format(created_at,'%Y-%m') = ? and is_published = ? order by created_at desc`
		querysql := `select * from posts where strftime('%Y-%m',created_at) = ? and is_published = ? order by created_at desc`
		rows, err = DB.Raw(querysql, condition, true).Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		var post Post
		DB.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return posts, nil
}

func CountPostByArchive(year, month string) (count int, err error) {
	if len(month) == 1 {
		month = "0" + month
	}
	condition := fmt.Sprintf("%s-%s", year, month)
	//querysql := `select count(*) from posts where date_format(created_at,'%Y-%m') = ? and is_published = ? order by created_at desc`
	querysql := `select count(*) from posts where strftime('%Y-%m',created_at) = ? and is_published = ?`
	err = DB.Raw(querysql, condition, true).Row().Scan(&count)
	return
}

// Tag
func (tag *Tag) Insert() error {
	return DB.FirstOrCreate(tag, "name = ?", tag.Name).Error
}

func ListTag() ([]*Tag, error) {
	var tags []*Tag
	rows, err := DB.Raw("select t.*,count(*) total from tags t inner join post_tags pt on t.id = pt.tag_id inner join posts p on pt.post_id = p.id where p.is_published = ? group by pt.tag_id", true).Rows()
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

func MustListTag() []*Tag {
	tags, _ := ListTag()
	return tags
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

func CountTag() int {
	var count int
	DB.Model(&Tag{}).Count(&count)
	return count
}

func ListAllTag() ([]*Tag, error) {
	var tags []*Tag
	err := DB.Model(&Tag{}).Find(&tags).Error
	return tags, err
}

// post_tags
func (pt *PostTag) Insert() error {
	return DB.FirstOrCreate(pt, "post_id = ? and tag_id = ?", pt.PostId, pt.TagId).Error
}

func DeletePostTagByPostId(postId uint) error {
	return DB.Delete(&PostTag{}, "post_id = ?", postId).Error
}

// user
// insert user
func (user *User) Insert() error {
	return DB.Create(user).Error
}

// update user
func (user *User) Update() error {
	return DB.Save(user).Error
}

//
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.First(&user, "email = ?", username).Error
	return &user, err
}

//
func (user *User) FirstOrCreate() (*User, error) {
	err := DB.FirstOrCreate(user, "github_login_id = ?", user.GithubLoginId).Error
	return user, err
}

func IsGithubIdExists(githubId string, id uint) (*User, error) {
	var user User
	err := DB.First(&user, "github_login_id = ? and id != ?", githubId, id).Error
	return &user, err
}

func GetUser(id interface{}) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	return &user, err
}

func (user *User) UpdateProfile(avatarUrl, nickName string) error {
	return DB.Model(user).Update(User{AvatarUrl: avatarUrl, NickName: nickName}).Error
}

func (user *User) UpdateEmail(email string) error {
	if len(email) > 0 {
		return DB.Model(user).Update("email", email).Error
	} else {
		return DB.Model(user).Update("email", gorm.Expr("NULL")).Error
	}
}

func (user *User) UpdateGithubUserInfo() error {
	var githubLoginId interface{}
	if len(user.GithubLoginId) == 0 {
		githubLoginId = gorm.Expr("NULL")
	} else {
		githubLoginId = user.GithubLoginId
	}
	return DB.Model(user).Update(map[string]interface{}{
		"github_login_id": githubLoginId,
		"avatar_url":      user.AvatarUrl,
		"github_url":      user.GithubUrl,
	}).Error
}

func (user *User) Lock() error {
	return DB.Model(user).Update(map[string]interface{}{
		"lock_state": user.LockState,
	}).Error
}

func ListUsers() ([]*User, error) {
	var users []*User
	err := DB.Find(&users, "is_admin = ?", false).Error
	return users, err
}

// Comment
func (comment *Comment) Insert() error {
	return DB.Create(comment).Error
}

func (comment *Comment) Update() error {
	return DB.Model(comment).UpdateColumn("read_state", true).Error
}

func SetAllCommentRead() error {
	return DB.Model(&Comment{}).Where("read_state = ?", false).Update("read_state", true).Error
}

func ListUnreadComment() ([]*Comment, error) {
	var comments []*Comment
	err := DB.Where("read_state = ?", false).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func MustListUnreadComment() []*Comment {
	comments, _ := ListUnreadComment()
	return comments
}

func (comment *Comment) Delete() error {
	return DB.Delete(comment, "user_id = ?", comment.UserID).Error
}

func ListCommentByPostID(postId string) ([]*Comment, error) {
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	rows, err := DB.Raw("select c.*,u.github_login_id nick_name,u.avatar_url,u.github_url from comments c inner join users u on c.user_id = u.id where c.post_id = ? order by created_at desc", uint(pid)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		DB.ScanRows(rows, &comment)
		comments = append(comments, &comment)
	}
	return comments, err
}

/*func GetComment(id interface{}) (*Comment, error) {
	var comment Comment
	err := DB.First(&comment, id).Error
	return &comment, err
}*/

func CountComment() int {
	var count int
	DB.Model(&Comment{}).Count(&count)
	return count
}

// Subscriber
func (s *Subscriber) Insert() error {
	return DB.FirstOrCreate(s, "email = ?", s.Email).Error
}

func (s *Subscriber) Update() error {
	return DB.Model(s).Update(map[string]interface{}{
		"verify_state":    s.VerifyState,
		"subscribe_state": s.SubscribeState,
		"out_time":        s.OutTime,
		"signature":       s.Signature,
		"secret_key":      s.SecretKey,
	}).Error
}

func ListSubscriber(invalid bool) ([]*Subscriber, error) {
	var subscribers []*Subscriber
	db := DB.Model(&Subscriber{})
	if invalid {
		db.Where("verify_state = ? and subscribe_state = ?", true, true)
	}
	err := db.Find(&subscribers).Error
	return subscribers, err
}

func CountSubscriber() (int, error) {
	var count int
	err := DB.Model(&Subscriber{}).Where("verify_state = ? and subscribe_state = ?", true, true).Count(&count).Error
	return count, err
}

func GetSubscriberByEmail(mail string) (*Subscriber, error) {
	var subscriber Subscriber
	err := DB.Find(&subscriber, "email = ?", mail).Error
	return &subscriber, err
}

func GetSubscriberBySignature(key string) (*Subscriber, error) {
	var subscriber Subscriber
	err := DB.Find(&subscriber, "signature = ?", key).Error
	return &subscriber, err
}

func GetSubscriberById(id uint) (*Subscriber, error) {
	var subscriber Subscriber
	err := DB.First(&subscriber, id).Error
	return &subscriber, err
}

// Link
func (link *Link) Insert() error {
	return DB.FirstOrCreate(link, "url = ?", link.Url).Error
}

func (link *Link) Update() error {
	return DB.Save(link).Error
}

func (link *Link) Delete() error {
	return DB.Delete(link).Error
}

func ListLinks() ([]*Link, error) {
	var links []*Link
	err := DB.Order("sort asc").Find(&links).Error
	return links, err
}

func MustListLinks() []*Link {
	links, _ := ListLinks()
	return links
}

func GetLinkById(id uint) (*Link, error) {
	var link Link
	err := DB.FirstOrCreate(&link, "id = ?", id).Error
	return &link, err
}

/*func GetLinkByUrl(url string) (*Link, error) {
	var link Link
	err := DB.Find(&link, "url = ?", url).Error
	return &link, err
}*/

func (sf SmmsFile) Insert() (err error) {
	err = DB.Create(&sf).Error
	return
}
