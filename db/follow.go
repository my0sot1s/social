package db

import (
	m "github.com/my0sot1s/social/mirrors"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) CountFollower(own string) (error, int) {
	collection := db.Db.C(followerCollection)
	count, err := collection.Find(bson.M{"own": own}).Count()
	if err != nil {
		return err, -1
	}
	return nil, count
}

// GetFollower return những ng đang follow own
func (db *DB) GetFollower(own string) (error, []*m.Follower) {
	collection := db.Db.C(followerCollection)
	var follower []*m.Follower
	err := collection.Find(bson.M{"own": own}).All(&follower)
	if err != nil {
		return err, nil
	}
	return nil, follower
}

// GetFollowing return những ng ban đang follow
func (db *DB) GetFollowing(follower string) (error, []*m.Follower) {
	collection := db.Db.C(followerCollection)
	var following []*m.Follower
	err := collection.Find(bson.M{"follower": follower}).All(&following)
	if err != nil {
		return err, nil
	}
	return nil, following
}

func (db *DB) ModifyFollower(t *m.Follower) (error, *m.Follower) {
	collection := db.Db.C(followerCollection)
	t.ID = bson.NewObjectId()
	err := collection.Insert(&t)
	if err != nil {
		return err, nil
	}
	return nil, t
}

func (db *DB) FollowUser(f *m.Follower) error {
	collection := db.Db.C(followerCollection)
	err := collection.Insert(f)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UnfollowUser(own, uid string) error {
	collection := db.Db.C(followerCollection)
	err := collection.Remove(bson.M{"own": own, "follower": uid})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IsFollow(own, uid string) (int, error) {
	collection := db.Db.C(followerCollection)
	num, err := collection.Find(bson.M{"own": own, "follower": uid}).Count()
	if err != nil {
		return 0, err
	}

	return num, nil
}
