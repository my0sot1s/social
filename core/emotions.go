package core

import (
	"encoding/json"
	"time"

	m "github.com/my0sot1s/social/mirrors"
	"github.com/my0sot1s/social/utils"
)

const timeExpired = 24 * 60

func (c *Social) CreateEmotion(mediasStr, by string) (error, *m.Emotion) {
	var medias []*m.Media
	if errMedia := utils.Str2T(mediasStr, &medias); errMedia != nil {
		utils.ErrLog(errMedia)
		return errMedia, nil
	}

	t := time.Now()
	emo := &m.Emotion{
		Medias:  medias,
		By:      by,
		Created: int(t.Unix()),
	}
	if err := c.rd.LPushItem(by, timeExpired, emo); err != nil {
		utils.ErrLog(err)
		return err, nil
	}

	return nil, emo
}

func (c *Social) GetEmotions(key string) []*m.Emotion {
	lEmo := make([]*m.Emotion, 0)
	if mk, err := c.rd.LRangeAll(key); err == nil {
		for _, v := range mk {
			// map[string]interface{}
			tempE := &m.Emotion{}
			var b []byte
			b, er := json.Marshal(v)
			if er != nil {
				utils.ErrLog(er)
				continue
			}
			json.Unmarshal(b, &tempE)
			lEmo = append(lEmo, tempE)
		}
		return lEmo
	}
	return nil
}

func (c *Social) GetEmotionsByMultipleKeys(key ...string) map[string][]*m.Emotion {
	lMult := make(map[string][]*m.Emotion)
	for _, v := range key {
		emoArray := c.GetEmotions(v)
		if len(emoArray) > 0 {
			lMult[v] = emoArray
		}
	}
	return lMult
}

func (c *Social) GetEmotionByUId(uid string) map[string][]*m.Emotion {
	err, users := c.LoadFollowingByUid(uid)
	utils.ErrLog(err)
	lUids := make([]string, 1)
	lUids[0] = uid
	for _, v := range users {
		lUids = append(lUids, v.GetID())
	}
	return c.GetEmotionsByMultipleKeys(lUids...)
}
