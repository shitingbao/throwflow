package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"douyin/internal/pkg/tool"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// 抖音开放平台达人视频表
type OpenDouyinVideo struct {
	ClientKey     string                 `json:"client_key" bson:"client_key"`
	OpenId        string                 `json:"open_id" bson:"open_id"`
	AwemeId       uint64                 `json:"aweme_id" bson:"aweme_id"`
	AccountId     string                 `json:"account_id" bson:"account_id"`
	Nickname      string                 `json:"nickname" bson:"nickname"`
	Avatar        string                 `json:"avatar" bson:"avatar"`
	Title         string                 `json:"title" bson:"title"`
	Cover         string                 `json:"cover" bson:"cover"`
	CreateTime    int64                  `json:"create_time" bson:"create_time"`
	IsReviewed    bool                   `json:"is_reviewed" bson:"is_reviewed"`
	ItemId        string                 `json:"item_id" bson:"item_id"`
	Statistics    domain.VideoStatistics `json:"statistics" bson:"statistics"`
	IsTop         bool                   `json:"is_top" bson:"is_top"`
	MediaType     int                    `json:"media_type" bson:"media_type"`
	ShareUrl      string                 `json:"share_url" bson:"share_url"`
	VideoId       string                 `json:"video_id" bson:"video_id"`
	VideoStatus   int32                  `json:"video_status" bson:"video_status"` // 表示视频状态。1:细化为5、6、7三种状态;2:不适宜公开;4:审核中;5:公开视频;6:好友可见;7:私密视频, 99:盟码标签，表示接口没有返回对应的数据
	ProductId     string                 `json:"product_id" bson:"product_id"`
	ProductName   string                 `json:"product_name" bson:"product_name"`
	ProductImg    string                 `json:"product_img" bson:"product_img"`
	IsUpdateCover uint8                  `json:"is_update_cover" bson:"is_update_cover"`
}

type TotalOpenDouyinVideo struct {
	Total int64 `json:"total" bson:"total"`
}

type openDouyinVideoRepo struct {
	data *Data
	log  *log.Helper
}

func (odv *OpenDouyinVideo) ToDomain() *domain.OpenDouyinVideo {
	return &domain.OpenDouyinVideo{
		ClientKey:     odv.ClientKey,
		OpenId:        odv.OpenId,
		AwemeId:       odv.AwemeId,
		AccountId:     odv.AccountId,
		Nickname:      odv.Nickname,
		Avatar:        odv.Avatar,
		Title:         odv.Title,
		Cover:         odv.Cover,
		CreateTime:    odv.CreateTime,
		IsReviewed:    odv.IsReviewed,
		ItemId:        odv.ItemId,
		Statistics:    odv.Statistics,
		IsTop:         odv.IsTop,
		MediaType:     odv.MediaType,
		ShareUrl:      odv.ShareUrl,
		VideoId:       odv.VideoId,
		VideoStatus:   odv.VideoStatus,
		ProductId:     odv.ProductId,
		ProductName:   odv.ProductName,
		ProductImg:    odv.ProductImg,
		IsUpdateCover: odv.IsUpdateCover,
	}
}

func NewOpenDouyinVideoRepo(data *Data, logger log.Logger) biz.OpenDouyinVideoRepo {
	return &openDouyinVideoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (odvr *openDouyinVideoRepo) GetByClientKeyAndOpenId(ctx context.Context, clientId, openId, mediaType string, videoStatus int32) (*domain.OpenDouyinVideo, error) {
	var openDouyinVideo OpenDouyinVideo

	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	if len(mediaType) > 0 {
		mediaTypes := strings.Split(mediaType, ",")

		imediaTypes := make([]uint64, 0, 0)

		for _, lmediaType := range mediaTypes {
			if imediaType, err := strconv.ParseUint(lmediaType, 10, 64); err == nil {
				imediaTypes = append(imediaTypes, imediaType)
			}
		}

		if err := collection.FindOne(ctx, bson.M{
			"$and": []bson.M{
				bson.M{"media_type": bson.D{{"$in", imediaTypes}}},
				bson.M{"video_status": bson.D{{"$eq", videoStatus}}},
				bson.M{"client_key": clientId},
				bson.M{"open_id": openId},
			},
		}, &options.FindOneOptions{
			Sort: bson.D{
				bson.E{"create_time", -1},
			},
		}).Decode(&openDouyinVideo); err != nil {
			return nil, err
		}
	} else {
		if err := collection.FindOne(ctx, bson.M{
			"$and": []bson.M{
				bson.M{"video_status": bson.D{{"$eq", videoStatus}}},
				bson.M{"client_key": clientId},
				bson.M{"open_id": openId},
			},
		}, &options.FindOneOptions{
			Sort: bson.D{
				bson.E{"create_time", -1},
			},
		}).Decode(&openDouyinVideo); err != nil {
			return nil, err
		}
	}

	return openDouyinVideo.ToDomain(), nil
}

func (odvr *openDouyinVideoRepo) List(ctx context.Context, isExistProduct uint8, videoIds, keyword, videoStatus, mediaType string, openDouyinTokens []*domain.OpenDouyinToken, pageNum, pageSize int64) ([]*domain.OpenDouyinVideo, error) {
	list := make([]*domain.OpenDouyinVideo, 0)
	var openDouyinVideos []*OpenDouyinVideo
	var aopenDouyinTokens []bson.M
	var avideoIds bson.A

	svideoIds := tool.RemoveEmptyString(strings.Split(videoIds, ","))

	for _, svideoId := range svideoIds {
		avideoIds = append(avideoIds, svideoId)
	}

	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var cursor *mongo.Cursor
	var and []bson.M
	var err error

	if len(svideoIds) > 0 {
		and = append(and, bson.M{"video_id": bson.D{{"$in", avideoIds}}})
	}

	if len(videoStatus) > 0 {
		ivideoStatus, _ := strconv.ParseUint(videoStatus, 10, 64)

		and = append(and, bson.M{"video_status": bson.D{{"$eq", int32(ivideoStatus)}}})
	}

	if len(mediaType) > 0 {
		mediaTypes := strings.Split(mediaType, ",")

		imediaTypes := make([]uint64, 0, 0)

		for _, lmediaType := range mediaTypes {
			if imediaType, err := strconv.ParseUint(lmediaType, 10, 64); err == nil {
				imediaTypes = append(imediaTypes, imediaType)
			}
		}

		and = append(and, bson.M{"media_type": bson.D{{"$in", imediaTypes}}})
	}

	if len(openDouyinTokens) > 0 {
		for _, openDouyinToken := range openDouyinTokens {
			aopenDouyinTokens = append(aopenDouyinTokens, bson.M{"$and": []bson.M{
				bson.M{"client_key": openDouyinToken.ClientKey},
				bson.M{"open_id": openDouyinToken.OpenId},
			},
			})
		}

		and = append(and, bson.M{"$or": aopenDouyinTokens})
	}

	if isExistProduct == 1 {
		and = append(and, bson.M{"product_id": bson.D{{"$ne", ""}}})
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		var or []bson.M

		or = append(or, bson.M{"nickname": bsonx.Regex(keyword, "")})
		or = append(or, bson.M{"title": bsonx.Regex(keyword, "")})

		if isExistProduct == 1 {
			or = append(or, bson.M{"product_name": bsonx.Regex(keyword, "")})
			or = append(or, bson.M{"product_id": bsonx.Regex(keyword, "")})
		}

		and = append(and, bson.M{"$or": or})
	}

	matchStage := bson.D{
		{"$match", bson.M{
			"$and": and,
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.M{"create_time": -1}},
	}

	if pageNum > 0 {
		skipStage := bson.D{
			{"$skip", (pageNum - 1) * pageSize},
		}

		limitStage := bson.D{
			{"$limit", pageSize},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage})
	} else {
		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, sortStage})
	}

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &openDouyinVideos)

	if err != nil {
		return nil, err
	}

	for _, openDouyinVideo := range openDouyinVideos {
		list = append(list, openDouyinVideo.ToDomain())
	}

	return list, nil
}

func (odvr *openDouyinVideoRepo) ListProduct(ctx context.Context, keyword string, openDouyinTokens []*domain.OpenDouyinToken, pageNum, pageSize int64) ([]*domain.OpenDouyinVideo, error) {
	list := make([]*domain.OpenDouyinVideo, 0)
	var openDouyinVideos []*OpenDouyinVideo
	var aopenDouyinTokens []bson.M

	for _, openDouyinToken := range openDouyinTokens {
		aopenDouyinTokens = append(aopenDouyinTokens, bson.M{"$and": []bson.M{
			bson.M{"client_key": openDouyinToken.ClientKey},
			bson.M{"open_id": openDouyinToken.OpenId},
		},
		})
	}

	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var matchStage bson.D

	if l := utf8.RuneCountInString(keyword); l > 0 {
		matchStage = bson.D{
			{"$match", bson.M{
				"$and": []bson.M{
					bson.M{"$or": aopenDouyinTokens},
					bson.M{"product_id": bson.D{{"$ne", ""}}},
					bson.M{"$or": []bson.M{
						bson.M{"product_id": bsonx.Regex(keyword, "")},
						bson.M{"product_name": bsonx.Regex(keyword, "")},
					}},
				},
			}},
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{
				"$and": []bson.M{
					bson.M{"$or": aopenDouyinTokens},
					bson.M{"product_id": bson.D{{"$ne", ""}}},
				},
			}},
		}
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":          "$product_id",
			"product_name": bson.D{{"$first", "$product_name"}},
			"product_img":  bson.D{{"$first", "$product_img"}},
		}},
	}

	projectAStage := bson.D{
		{"$project", bson.M{
			"_id":          0,
			"product_id":   "$_id",
			"product_name": 1,
			"product_img":  1,
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.M{"product_id": 1}},
	}

	var cursor *mongo.Cursor
	var err error

	if pageNum > 0 {
		skipStage := bson.D{
			{"$skip", (pageNum - 1) * pageSize},
		}

		limitStage := bson.D{
			{"$limit", pageSize},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectAStage, sortStage, skipStage, limitStage})
	} else {
		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectAStage, sortStage})
	}

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &openDouyinVideos)

	if err != nil {
		return nil, err
	}

	for _, openDouyinVideo := range openDouyinVideos {
		list = append(list, openDouyinVideo.ToDomain())
	}

	return list, nil
}

func (odvr *openDouyinVideoRepo) ListVideoId(ctx context.Context, pageNum, pageSize int64) ([]*domain.OpenDouyinVideo, error) {
	list := make([]*domain.OpenDouyinVideo, 0)
	var openDouyinVideos []*OpenDouyinVideo

	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var matchStage bson.D
	matchStage = bson.D{
		{"$match", bson.M{
			"$and": []bson.M{
				bson.M{"video_status": 1},
				bson.M{"is_update_cover": bson.D{{"$exists", false}}},
			},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id": bson.D{{"video_id", "$video_id"}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"_id":      0,
			"video_id": "$_id.video_id",
		}},
	}

	var cursor *mongo.Cursor
	var err error

	if pageNum > 0 {
		skipStage := bson.D{
			{"$skip", (pageNum - 1) * pageSize},
		}

		limitStage := bson.D{
			{"$limit", pageSize},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage, skipStage, limitStage})
	} else {
		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	}

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &openDouyinVideos)

	if err != nil {
		return nil, err
	}

	for _, openDouyinVideo := range openDouyinVideos {
		list = append(list, &domain.OpenDouyinVideo{
			VideoId: openDouyinVideo.VideoId,
		})
	}

	return list, nil
}

func (odvr *openDouyinVideoRepo) ListProductsByTokens(ctx context.Context, productOutId uint64, claimTime string, openDouyinTokens []*domain.OpenDouyinToken) ([]*domain.OpenDouyinVideo, error) {
	list := make([]*domain.OpenDouyinVideo, 0)
	var openDouyinVideos []*OpenDouyinVideo
	var aopenDouyinTokens []bson.M

	for _, openDouyinToken := range openDouyinTokens {
		aopenDouyinTokens = append(aopenDouyinTokens, bson.M{"$and": []bson.M{
			{"client_key": openDouyinToken.ClientKey},
			{"open_id": openDouyinToken.OpenId},
		}})
	}

	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	t, err := time.Parse("2006-01-02 15:04:05", claimTime)

	if err != nil {
		return nil, err
	}

	unixTimestamp := t.Unix()
	matchStage := bson.M{
		"$or":        aopenDouyinTokens,
		"product_id": strconv.FormatUint(productOutId, 10),
		"create_time": bson.M{
			"$gte": unixTimestamp,
		},
	}

	cursor, err := collection.Find(ctx, matchStage)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &openDouyinVideos); err != nil {
		return nil, err
	}

	for _, openDouyinVideo := range openDouyinVideos {
		list = append(list, openDouyinVideo.ToDomain())
	}

	return list, nil
}

func (odvr *openDouyinVideoRepo) Count(ctx context.Context, isExistProduct uint8, videoIds, keyword, videoStatus, mediaType string, openDouyinTokens []*domain.OpenDouyinToken) (int64, error) {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var and []bson.M
	var aopenDouyinTokens bson.A
	var avideoIds bson.A

	for _, openDouyinToken := range openDouyinTokens {
		aopenDouyinTokens = append(aopenDouyinTokens, bson.M{"$and": []bson.M{
			bson.M{"client_key": openDouyinToken.ClientKey},
			bson.M{"open_id": openDouyinToken.OpenId},
		},
		})
	}

	svideoIds := tool.RemoveEmptyString(strings.Split(videoIds, ","))

	for _, svideoId := range svideoIds {
		avideoIds = append(avideoIds, svideoId)
	}

	if len(svideoIds) > 0 {
		and = append(and, bson.M{"video_id": bson.D{{"$in", avideoIds}}})
	}

	if len(videoStatus) > 0 {
		ivideoStatus, _ := strconv.ParseUint(videoStatus, 10, 64)

		and = append(and, bson.M{"video_status": bson.D{{"$eq", int32(ivideoStatus)}}})
	}

	if len(mediaType) > 0 {
		mediaTypes := strings.Split(mediaType, ",")

		imediaTypes := make([]uint64, 0, 0)

		for _, lmediaType := range mediaTypes {
			if imediaType, err := strconv.ParseUint(lmediaType, 10, 64); err == nil {
				imediaTypes = append(imediaTypes, imediaType)
			}
		}

		and = append(and, bson.M{"media_type": bson.D{{"$in", imediaTypes}}})
	}

	and = append(and, bson.M{"$or": aopenDouyinTokens})

	if isExistProduct == 1 {
		and = append(and, bson.M{"product_id": bson.D{{"$ne", ""}}})
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		var or []bson.M

		or = append(or, bson.M{"nickname": bsonx.Regex(keyword, "")})
		or = append(or, bson.M{"Title": bsonx.Regex(keyword, "")})

		if isExistProduct == 1 {
			or = append(or, bson.M{"product_name": bsonx.Regex(keyword, "")})
			or = append(or, bson.M{"product_id": bsonx.Regex(keyword, "")})
		}

		and = append(and, bson.M{"$or": or})
	}

	return collection.CountDocuments(ctx, bson.M{"$and": and})
}

func (odvr *openDouyinVideoRepo) CountProduct(ctx context.Context, keyword string, openDouyinTokens []*domain.OpenDouyinToken) (int64, error) {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var aopenDouyinTokens []bson.M

	for _, openDouyinToken := range openDouyinTokens {
		aopenDouyinTokens = append(aopenDouyinTokens, bson.M{"$and": []bson.M{
			bson.M{"client_key": openDouyinToken.ClientKey},
			bson.M{"open_id": openDouyinToken.OpenId},
		},
		})
	}

	var matchStage bson.D

	if l := utf8.RuneCountInString(keyword); l > 0 {
		matchStage = bson.D{
			{"$match", bson.M{
				"$and": []bson.M{
					bson.M{"$or": aopenDouyinTokens},
					bson.M{"product_id": bson.D{{"$ne", ""}}},
					bson.M{"$or": []bson.M{
						bson.M{"product_id": bsonx.Regex(keyword, "")},
						bson.M{"product_name": bsonx.Regex(keyword, "")},
					}},
				},
			}},
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{
				"$and": []bson.M{
					bson.M{"$or": aopenDouyinTokens},
					bson.M{"product_id": bson.D{{"$ne", ""}}},
				},
			}},
		}
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":          "$product_id",
			"product_name": bson.D{{"$first", "$product_name"}},
			"product_img":  bson.D{{"$first", "$product_img"}},
		}},
	}

	countStage := bson.D{
		{"$count", "total"},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, countStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var totalOpenDouyinVideos []*TotalOpenDouyinVideo

	err = cursor.All(ctx, &totalOpenDouyinVideos)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, totalOpenDouyinVideo := range totalOpenDouyinVideos {
		total = totalOpenDouyinVideo.Total
	}

	return total, nil
}

func (odvr *openDouyinVideoRepo) CountVideoId(ctx context.Context) (int64, error) {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	var matchStage bson.D
	matchStage = bson.D{
		{"$match", bson.M{
			"$and": []bson.M{
				bson.M{"video_status": 1},
				bson.M{"is_update_cover": bson.D{{"$exists", false}}},
			},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id": bson.D{{"video_id", "$video_id"}},
		}},
	}

	countStage := bson.D{
		{"$count", "total"},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, countStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var totalOpenDouyinVideos []*TotalOpenDouyinVideo

	err = cursor.All(ctx, &totalOpenDouyinVideos)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, totalOpenDouyinVideo := range totalOpenDouyinVideos {
		total = totalOpenDouyinVideo.Total
	}

	return total, nil
}

func (odvr *openDouyinVideoRepo) SaveIndex(ctx context.Context) {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	isNotExistIndexA := true
	isNotExistIndexB := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "client_key_-1_open_id_-1_video_id_-1" {
				isNotExistIndexA = false
			}

			if indexSpecification.Name == "video_id_-1" {
				isNotExistIndexB = false
			}
		}

		if isNotExistIndexA {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "client_key", Value: -1},
					{Key: "open_id", Value: -1},
					{Key: "video_id", Value: -1},
				},
			})
		}

		if isNotExistIndexB {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "video_id", Value: -1},
				},
			})
		}
	}
}

func (odvr *openDouyinVideoRepo) Upsert(ctx context.Context, in *domain.OpenDouyinVideo) error {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"client_key", in.ClientKey},
		{"open_id", in.OpenId},
		{"video_id", in.VideoId},
	}, bson.D{
		{"$setOnInsert", bson.D{
			{"cover", in.Cover},
		}},
		{"$set", bson.D{
			{"aweme_id", in.AwemeId},
			{"account_id", in.AccountId},
			{"nickname", in.Nickname},
			{"avatar", in.Avatar},
			{"title", in.Title},
			{"create_time", in.CreateTime},
			{"is_reviewed", in.IsReviewed},
			{"item_id", in.ItemId},
			{"statistics", in.Statistics},
			{"is_top", in.IsTop},
			{"media_type", in.MediaType},
			{"share_url", in.ShareUrl},
			{"video_status", in.VideoStatus},
			{"product_id", in.ProductId},
			{"product_name", in.ProductName},
			{"product_img", in.ProductImg},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (odvr *openDouyinVideoRepo) UpdateIsUpdateCoverAndProductId(ctx context.Context, videoId, cover, productId string) error {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	if _, err := collection.UpdateMany(ctx, bson.D{
		{"video_id", videoId},
	}, bson.D{
		{"$set", bson.D{
			{"cover", cover},
			{"product_id", productId},
			{"is_update_cover", 1},
		}},
	}); err != nil {
		return err
	}

	return nil
}

func (odvr *openDouyinVideoRepo) UpdateVideoStatus(ctx context.Context, videoStatus int32, videoIds []string) error {
	collection := odvr.data.mdb.Database(odvr.data.conf.Mongo.Dbname).Collection("open_douyin_video")

	if _, err := collection.UpdateMany(ctx, bson.D{
		{"video_id", bson.D{{"$nin", videoIds}}},
	}, bson.D{
		{"$set", bson.D{
			{"video_status", videoStatus},
		}},
	}); err != nil {
		return err
	}

	return nil
}
