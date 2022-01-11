package main

import (
	"context"
	"fmt"
	"grpc-go-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	//if we crash the go code, we ger the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Hello Blog Client")

	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/public/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServicClient(cc)
	//fmt.Printf("Create client %f\n\n", c)

	// create the blog
	fmt.Println("Create the blog")
	blog := &blogpb.Blog{
		AuthorId: "Bunyawat",
		Title:    "My First Blog",
		Content:  "Content of my first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(),
		&blogpb.CreateBlogRequest{
			Blog: blog,
		})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog had been create %v\n", createBlogRes)
	blogId := createBlogRes.GetBlog().GetId()

	// read the blog
	fmt.Println("Read the blog")

	_, readErr := c.ReadBlog(context.Background(),
		&blogpb.ReadBlogRequest{
			BlogId: "random_id",
		})
	if readErr != nil {
		fmt.Printf("Error happend while reading %v\n", readErr)
	}

	req := &blogpb.ReadBlogRequest{
		BlogId: blogId,
	}
	readBlogRes, readErr2 := c.ReadBlog(context.Background(), req)
	if readErr2 != nil {
		fmt.Printf("Error happend while reading %v\n ", readErr2)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update the blog
	fmt.Println("Update the blog")
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Updated Author Id",
		Title:    "Updated Title",
		Content:  "Updated Content",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(),
		&blogpb.UpdateBlogRequest{
			Blog: newBlog,
		})
	if updateErr != nil {
		fmt.Printf("Error happened while updateing: %v\n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	//delete Blog
	fmt.Println("Delete the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(),
		&blogpb.DeleteBlogRequest{
			BlogId: blogId,
		})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

}
