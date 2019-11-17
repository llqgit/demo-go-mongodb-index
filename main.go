package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 创建一个数据库链接
func CreateDb() *mongo.Client {
	uri := "mongodb://localhost:27017"                                 // 本地链接
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) // 链接超时 2 秒
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return db
}

// 测试的实体数据
type User struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

func main() {
	DB := CreateDb() // 创建一个数据库链接
	collection := DB.Database("test").Collection("user")
	ctx, _ := context.WithCancel(context.Background())

	data := User{
		Name: "llqgit",
		Age:  27,
	}

	var err error

	// 增加
	addRes, err := collection.InsertOne(ctx, data)
	fmt.Printf("增加结果：%+v : %+v\n", addRes, err)

	// 删除
	delRes, err := collection.DeleteOne(ctx, bson.M{"name": "llqgit"})
	fmt.Printf("删除结果：%+v : %+v\n", delRes, err)

	// 再增加
	addRes, err = collection.InsertOne(ctx, data)
	fmt.Printf("增加结果：%+v : %+v\n", addRes, err)

	// 修改
	filter := bson.D{{"name", "hello"}}
	update := bson.D{{"$set", bson.D{{"name", "llqgit-1"}}}}
	modRes, err := collection.UpdateOne(ctx, filter, update)
	fmt.Printf("修改结果：%+v : %+v\n", modRes, err)

	// 查询
	findRes := collection.FindOne(ctx, bson.M{"name": "llqgit-1"})
	var test User
	err = findRes.Decode(&test) // 这个对象必须为一个地址，一定要加 &
	fmt.Printf("查询结果：%+v : %+v\n", test, err)

	// 创建索引
	indexName, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"name", 1}},
	})
	if err != nil {
		fmt.Println("创建索引", err)
	} else {
		fmt.Println("创建索引", indexName)
	}

	// 删除索引
	dropIndexRet, err := collection.Indexes().DropOne(ctx, indexName)
	if err != nil {
		fmt.Println("删除索引", err)
	} else {
		fmt.Println("删除索引", dropIndexRet)
	}
}
