package core

import (
	dbase "github.com/my0sot1s/social/db"
	m "github.com/my0sot1s/social/mongo"
	"github.com/my0sot1s/social/redis"
)

type Core struct {
	Db    ICore
	token *JWTAuthentication
	rd    *redis.RedisCli
}

type ICore interface {
	GetUserByUname(username string) (error, *m.User)
	GetUserByEmail(email string) (error, *m.User)
	GetPost(limit, page int, userID string) (error, []*m.Post)
	GetPostById(postID string) (error, *m.Post)
	GetFeed(limit, page int, userID string) (error, []*m.Feed)
	GetFollower(own string) (error, []*m.Follower)
	GetFollowing(follower string) (error, []*m.Follower)
	CountLike(postID string) (error, int)
	GetAlbum(AlbumID string) (error, *m.Album)
	GetAlbumByAuthor(limit, page int, userId string) (error, []*m.Album)
	GetComments(limit, page int, postID string) (error, []*m.Comment)
	GetLikes(postID string) (error, []*m.Like)
	IsUserLikePost(pid, uid string) (error, bool)
	GetPosts(pIDs []string) (error, []*m.Post)
	GetUserOwns(uIDs []string) (error, []*m.User)
	CreatePost(p *m.Post) (error, *m.Post)
	CreateComment(c *m.Comment) (error, *m.Comment)
	CreateFeed(f *m.Feed) (error, *m.Feed)
	CreateFeeds(feeds []*m.Feed) (error, []interface{})
	CreateUser(u *m.User) (error, *m.User)
	ModifyFollower(t *m.Follower) (error, *m.Follower)
	CreateAlbum(a *m.Album) (error, *m.Album)
	HitLikePost(postID, userID string) error
	UnlikePost(postID, userID string) error
	//
	GetUsersLikePost(userIDs []string) (error, []*m.User)
	FollowUser(f *m.Follower) error
	UnfollowUser(own, uid string) error
	//
}

func (c *Core) Config(db *dbase.DB, rd *redis.RedisCli, privateKeyPath, PublicKeyPath string) {
	// connect to drive Mongo
	c.Db = db
	// connect token access
	c.token = &JWTAuthentication{}
	c.token.Config(privateKeyPath, PublicKeyPath)
	// connect drive Redis
	c.rd = rd
}
