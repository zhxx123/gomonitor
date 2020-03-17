package cron

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo/bson"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
)

func yesterdayCron() {
	var yesterdaySignupUserCount uint // 昨日注册用户数
	var yesterdayPV uint              // 昨日PV
	var yesterdayUV uint              // 昨日UV

	todayTime := utils.GetTodayTime().Unix()
	yesterdayTime := utils.GetYesterdayTime().Unix()

	if err := db.DB.Model(&model.User{}).Where("register_at >= ? AND register_at < ?", yesterdayTime, todayTime).Count(&yesterdaySignupUserCount).Error; err != nil {
		logrus.Error(err)
		return
	}

	var pvCount map[string]uint
	pvErr := db.MongoDB.C("userVisit").Pipe(
		[]bson.M{
			{"$match": bson.M{
				"date": bson.M{
					"$gte": yesterdayTime,
					"$lt":  todayTime,
				},
			}},
			{"$count": "pv"},
		},
	).AllowDiskUse().One(&pvCount)

	if pvErr != nil {
		fmt.Println(pvErr)
	} else {
		yesterdayPV = pvCount["pv"]
	}

	var uvCount map[string]uint
	uvErr := db.MongoDB.C("userVisit").Pipe(
		[]bson.M{
			{"$match": bson.M{
				"date": bson.M{
					"$gte": yesterdayTime,
					"$lt":  todayTime,
				},
			}},
			{
				"$group": bson.M{
					"_id": "$clientID",
				},
			},
			{"$count": "uv"},
		},
	).AllowDiskUse().One(&uvCount)

	if uvErr != nil {
		logrus.Error(uvErr)
	} else {
		yesterdayUV = uvCount["uv"]
	}

	yesterdayStr := utils.GetYesterdayYMD("-")
	_, err := db.MongoDB.C("yesterdayStats").Upsert(bson.M{
		"date": yesterdayStr,
	}, bson.M{
		"$set": bson.M{
			"date":            yesterdayStr,
			"signupUserCount": yesterdaySignupUserCount,
			"pv":              yesterdayPV,
			"uv":              yesterdayUV,
		},
	})
	if err != nil {
		logrus.Error(err)
		return
	}
}
