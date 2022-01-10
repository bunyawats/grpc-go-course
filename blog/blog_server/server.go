package main

import (
	"context"
	"fmt"
	"grpc-go-course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServicServer
}

type blogItem struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	AuthorID string        `bson:"author_id"`
	Content  string        `bson:"content"`
	Title    string        `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context,
	reg *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := reg.GetBlog()

	data := blogItem{
		AuthorID: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	result, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := result.InsertedID.(bson.ObjectId)
	if !ok {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprint("Can not convert to OID"),
		)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func main() {

	//if we crash the go code, we ger the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	// connect to MongoDB
	//connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to cannect to MongoDB: %v\n ", err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/private/server.crt"
		keyFile := "ssl/private/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServicServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("faild to serve %v\n", err)
		}
	}()

	defer func() {

		fmt.Println("Stoping the server")
		s.Stop()
		fmt.Println("Closing the listener")
		lis.Close()
		fmt.Println("Closing MongoDB connection")
		client.Disconnect(ctx)
		fmt.Println("End of Program")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	//Block until a signal is received
	<-ch

}
