package dao

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"guoguo/source"
	"guoguo/source/tool"
	"user/model"
)

var dao *Dao

type Dao struct {
	collection *mongo.Collection
	logger     *logrus.Logger
	tool       *tool.Tool
}

func init() {
	dao = new(Dao)
	dao.collection = source.GetDatabase().Collection("user")
	dao.logger = source.GetLogger()
	dao.tool = tool.GetTool()
	//dao.CtxManager = tool.GetCtxManager()
}

func getDao() *Dao {
	return dao
}

//func (d *Dao) ctx() context.Context {
//	return source.NewCtx()
//}

func (d *Dao) Find(ctx context.Context, filter interface{}, opt *options.FindOptions) ([]*model.User, error) {
	cur, err := d.collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	users := make([]*model.User, 0)
	for cur.Next(ctx) {
		item := new(model.User)
		err = cur.Decode(item)
		if err != nil {
			return nil, err
		}
		users = append(users, item)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}

	defer func() {
		_ = cur.Close(ctx)
	}()

	return users, err
}

func (d *Dao) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opt *options.UpdateOptions) (*mongo.UpdateResult, error) {
	res, err := d.collection.UpdateOne(ctx, filter, update, opt)
	return res, err
}

func (d *Dao) FindOne(ctx context.Context, filter interface{}, opt *options.FindOneOptions) (*model.User, error) {
	one := d.collection.FindOne(ctx, filter, opt)
	if one == nil {
		return nil, errors.New("not find")
	}
	u := new(model.User)
	err := one.Decode(u)
	if err != nil {
		return nil, err
	}
	u.ID = d.tool.ConObjectIDToString(u.ObjID)
	return u, err
}

func (d *Dao) InsertOne(ctx context.Context, doc interface{}, opt *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	one, err := d.collection.InsertOne(ctx, doc, opt)
	return one, err
}