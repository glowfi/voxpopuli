package scraper

// SubredditRule represents a subreddit rule
type SubredditRule struct {
	ShortName   string `json:"short_name"`
	Description string `json:"description"`
}

// FlairEmoji represents a flair emoji
type FlairEmoji struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

// FlairText represents a flair text
type FlairText struct {
	Text string `json:"text"`
}

// SubredditFlair represents a subreddit flair
type SubredditFlair struct {
	RichText        []interface{} `json:"rich_text"`
	FullText        string        `json:"full_text"`
	BackgroundColor string        `json:"background_color"`
}

// Subreddit represents a subreddit
type Subreddit struct {
	ID                    string           `json:"id"`
	Title                 string           `json:"title"`
	PublicDescription     string           `json:"public_description"`
	CommunityIcon         string           `json:"community_icon"`
	BannerBackgroundImage string           `json:"banner_background_image"`
	Topic                 string           `json:"topic"`
	Rules                 []SubredditRule  `json:"rules"`
	Flairs                []SubredditFlair `json:"flairs"`
	UserFlairs            []SubredditFlair `json:"user_flairs"`
	KeyColor              string           `json:"key_color"`
	PrimaryColor          string           `json:"primary_color"`
	BannerBackgroundColor string           `json:"banner_background_color"`
	CreatedUTC            float64          `json:"created_utc"`
	CreatedHuman          string           `json:"created_human"`
	Subscribers           int              `json:"subscribers"`
	SubscribersHuman      string           `json:"subscribers_human"`
	Members               []User           `json:"members"`
	Moderators            []User           `json:"moderators"`
	Over18                bool             `json:"over18"`
	SpoilersEnabled       bool             `json:"spoilers_enabled"`
}

// Awards represents an award
type Awards struct {
	Title     string `json:"title"`
	ImageLink string `json:"image_link"`
}

// Comment represents a comment
type Comment struct {
	Author          string    `json:"author"`
	AuthorFullname  string    `json:"author_fullname"`
	AuthorFlairText string    `json:"author_flair_text"`
	Body            string    `json:"body"`
	BodyHTML        string    `json:"body_html"`
	Ups             int       `json:"ups"`
	Score           int       `json:"score"`
	CreatedUTC      float64   `json:"created_utc"`
	Replies         []Comment `json:"replies"`
}

// ImageMultiResolution represents an image multi-resolution
type ImageMultiResolution struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	U string `json:"u"`
}

// ImageResolution represents an image resolution
type ImageResolution struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// Image represents an image
type Image struct {
	ID          string        `json:"id"`
	Resolutions []interface{} `json:"resolutions"`
	Type        string        `json:"_type"`
}

// Gif represents a GIF
type Gif struct {
	ID          string        `json:"id"`
	Resolutions []interface{} `json:"resolutions"`
	Type        string        `json:"_type"`
}

// Video represents a video
type Video struct {
	ID               string `json:"id"`
	FallbackURL      string `json:"fallback_url"`
	HLSURL           string `json:"hls_url"`
	DashURL          string `json:"dash_url"`
	ScrubberMediaURL string `json:"scrubber_media_url"`
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Type             string `json:"_type"`
}

// VideoMulti represents a video multi
type VideoMulti struct {
	ID      string `json:"id"`
	HLSURL  string `json:"hlsUrl"`
	DashURL string `json:"dashUrl"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Type    string `json:"_type"`
}

// GalleryImageResolutions represents a gallery image resolutions
type GalleryImageResolutions struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
	U  string `json:"u"`
}

// Gallery represents a gallery
type Gallery struct {
	ID     string                      `json:"id"`
	Images [][]GalleryImageResolutions `json:"images"`
	Type   string                      `json:"_type"`
}

// Link represents a link
type Link struct {
	ID   string `json:"id"`
	Link string `json:"link"`
	Type string `json:"_type"`
}

// MediaContent represents media content
type MediaContent struct {
	Type    string      `json:"_type"`
	Content interface{} `json:"content"`
}

// Post represents a post
type Post struct {
	ID              string       `json:"id"`
	Subreddit       string       `json:"subreddit"`
	SubredditID     string       `json:"subreddit_id"`
	Title           string       `json:"title"`
	Author          string       `json:"author"`
	AuthorFullname  string       `json:"author_fullname"`
	AuthorFlairText string       `json:"author_flair_text"`
	LinkFlairText   string       `json:"link_flair_text"`
	NumComments     int          `json:"num_comments"`
	Ups             int          `json:"ups"`
	Awards          []Awards     `json:"awards"`
	Comments        []Comment    `json:"comments"`
	MediaContent    MediaContent `json:"media_content"`
	CreatedUTC      float64      `json:"created_utc"`
	CreatedHuman    string       `json:"created_human"`
	Text            string       `json:"text"`
	TextHTML        string       `json:"text_html"`
	Over18          bool         `json:"over_18"`
	Spoiler         bool         `json:"spoiler"`
}

// Trophies represents a trophy
type Trophies struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"image_link"`
}

// User represents a user
type User struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	CakeDayUTC        int        `json:"cake_day_utc"`
	CakeDayHuman      string     `json:"cake_day_human"`
	Age               string     `json:"age"`
	AvatarImg         string     `json:"avatar_img"`
	BannerImg         string     `json:"banner_img"`
	PublicDescription string     `json:"public_description"`
	Over18            bool       `json:"over18"`
	KeyColor          string     `json:"keycolor"`
	PrimaryColor      string     `json:"primarycolor"`
	IconColor         string     `json:"iconcolor"`
	Suspended         bool       `json:"suspended"`
	Trophies          []Trophies `json:"trophies"`
}
