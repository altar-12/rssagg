package main

import (
	"time"

	"github.com/altar-12/rssagg/internal/database"
	"github.com/google/uuid"
)

type Usera struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

type Feeda struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"api_key"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollowa struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Posta struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Url         string     `json:"url"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func ConvertDBUserToUsera(user database.User) User {
	return User{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Name,
		user.ApiKey,
	}
}

func ConvertDBFeedToFeeda(feed database.Feed) Feed {
	return Feed{
		feed.ID,
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.Name,
		feed.Url,
		feed.UserID,
	}
}

func ConvertDBFeedFollowToFeedFollowa(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		feedFollow.ID,
		feedFollow.CreatedAt,
		feedFollow.UpdatedAt,
		feedFollow.UserID,
		feedFollow.FeedID,
	}
}

func ConvertDBPostToPosta(post database.Post) Post {
	var publishedAt *time.Time
	if post.PublishedAt.Valid {
		publishedAt = &post.PublishedAt.Time
	}
	return Post{
		post.ID,
		post.CreatedAt,
		post.UpdatedAt,
		post.Title,
		post.Description,
		post.Url,
		publishedAt,
		post.FeedID,
	}
}

func ConvertDBFeedsToFeedsa(DBFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, DBFeed := range DBFeeds {
		feeds = append(feeds, ConvertDBFeedToFeed(DBFeed))
	}
	return feeds
}

func ConvertDBFeedFollowsToFeedFollowsa(DBFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}
	for _, DBFeedFollow := range DBFeedFollows {
		feedFollows = append(feedFollows, ConvertDBFeedFollowToFeedFollow(DBFeedFollow))
	}
	return feedFollows
}

func ConvertDBPostsToPostsa(DBPosts []database.Post) []Post {
	posts := []Post{}
	for _, post := range DBPosts {
		posts = append(posts, ConvertDBPostToPost(post))
	}
	return posts
}
