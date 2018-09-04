package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Tag 标签
type FileMeta struct {
	ID     bson.ObjectId   `bson:"_id"`
	Name   string          `bson:"name"`
	MD5    string          `bson:"md5"`
	PathID bson.ObjectId   `bson:"path"`  // 物理文件归属者ID
	OwnID  []bson.ObjectId `bson:"ownId"` // 引用拥有者
}

// TagModel 标签数据库
type FileModel struct {
	DB *mgo.Collection
}

// Add 添加文件
func (m *FileModel) Add(name, md5, ownID string) error {
	return m.DB.Insert(FileMeta{
		ID:     bson.NewObjectId(),
		Name:   name,
		MD5:    md5,
		PathID: bson.ObjectIdHex(ownID),
		OwnID:  []bson.ObjectId{bson.ObjectIdHex(ownID)},
	})
}

func (m *FileModel) FindByMD5(md5 string) (res FileMeta, err error) {
	err = m.DB.Find(bson.M{"md5": md5}).One(&res)
	return
}

func (m *FileModel) FindByID(id string) (res FileMeta, err error) {
	if !bson.IsObjectIdHex(id) {
		return
	}
	err = m.DB.FindId(bson.ObjectIdHex(id)).One(&res)
	return
}

func (m *FileModel) AddOwn(ownID, token string) error {
	return m.DB.Update(bson.M{"md5": bson.ObjectIdHex(token)}, bson.M{"$addToSet": bson.M{"ownId": bson.ObjectIdHex(ownID)}})
}

func (m *FileModel) DeleteOwn(ownID, token string) error {
	return m.DB.Update(bson.M{"md5": bson.ObjectIdHex(token)}, bson.M{"$pull": bson.M{"ownId": bson.ObjectIdHex(ownID)}})
}

func (m *FileModel) DeleteFileMeta(ID string) error {
	return m.DB.RemoveId(bson.ObjectIdHex(ID))
}

func (m *FileModel) IsExist(md5 string) bool {
	n, err := m.DB.Find(bson.M{"md5": md5}).Count()
	return n != 0 && err == nil
}
