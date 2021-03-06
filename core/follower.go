package core

import (
	"time"

	m "github.com/my0sot1s/social/mirrors"
	"github.com/my0sot1s/social/utils"
)

func (p *Social) CountFollowerByOwner(owner string) (error, int) {
	err, count := p.Db.CountFollower(owner)
	if utils.ErrLog(err); err != nil {
		return nil, -1
	}
	return nil, count
}

// LoadFollowerByOwner ai đang follow bạn `owner`
func (p *Social) LoadFollowerByOwner(owner string) (error, []*m.User) {
	err, follower := p.Db.GetFollower(owner)
	listUserID := make([]string, 0)
	for _, f := range follower {
		listUserID = append(listUserID, f.GetFollower())
	}
	if len(listUserID) == 0 {
		return nil, nil
	}
	if err != nil {
		utils.ErrLog(err)
		return err, nil
	}
	err2, users := p.Db.GetUserByIds(listUserID)
	// get List user by 1 request
	if err2 != nil {
		utils.ErrLog(err2)
		return err2, nil
	}
	return nil, users
}

// GetFollowingByUid bạn đang follow những ai `userId`
func (p *Social) LoadFollowingByUid(uid string) (error, []*m.User) {
	err, follower := p.Db.GetFollowing(uid)
	listUserID := make([]string, 0)
	for _, f := range follower {
		listUserID = append(listUserID, f.GetOwn())
	}
	if len(listUserID) == 0 {
		return nil, nil
	}
	if err != nil {
		utils.ErrLog(err)
		return err, nil
	}

	err2, users := p.Db.GetUserByIds(listUserID)
	if err2 != nil {
		utils.ErrLog(err2)
		return err2, nil
	}
	return nil, users
}

func (p *Social) CombineFollowToFeed(uid, owner string) {
	// find post
	_, ps := p.Db.GetPost(-100, "", owner)
	// make array
	listFeed := make([]*m.Feed, 0)
	for _, v := range ps {
		listFeed = append(listFeed, &m.Feed{
			Created:    v.GetCreated(),
			PostID:     v.GetID(),
			ConsumerID: uid,
			Author:     owner,
		})
	}
	err, _ := p.Db.CreateFeeds(listFeed)
	if err != nil {
		utils.ErrLog(err)
	}
}

func (p *Social) UpsertFollowAnUser(uid, owner string) error {
	follow := &m.Follower{
		Follower: uid,
		Own:      owner,
		Created:  time.Now(),
	}
	err := p.Db.FollowUser(follow)
	if err != nil {
		return err
	}
	// add to feed
	go p.CombineFollowToFeed(uid, owner)
	return nil
}

func (p *Social) ReduceFeeds(uid, owner string) {
	err := p.Db.DeleteFeed(uid, owner)
	if err != nil {
		utils.ErrLog(err)
	}
}

func (p *Social) RemoveFollowAnUser(uid, owner string) error {
	err := p.Db.UnfollowUser(owner, uid)
	// remove all feed
	utils.Log(uid, ",", owner)
	if err != nil {
		return err
	}
	go p.ReduceFeeds(uid, owner)
	return nil
}

func (p *Social) CheckFollow(owner, uid string) bool {
	numb, err := p.Db.IsFollow(owner, uid)
	if numb > 0 {
		return true
	}
	if err != nil {
		utils.ErrLog(err)
	}
	return false
}
