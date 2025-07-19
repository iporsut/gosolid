package main

import (
	"cmp"
	"context"
	"errors"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID    int
	Title string
	Body  string
}

var (
	inmemoryPostDB = make(map[int]Post)
	idPostCounter  int
	idPostMutex    sync.Mutex
)

type DB struct{}

var ErrNotFound = errors.New("not found")

func (d *DB) AddPost(ctx context.Context, newPost Post) (Post, error) {
	idPostMutex.Lock()
	defer idPostMutex.Unlock()
	idPostCounter++
	newPost.ID = idPostCounter
	inmemoryPostDB[idPostCounter] = newPost

	return newPost, nil
}

func (d *DB) GetPostByID(ctx context.Context, id int) (Post, error) {
	post, ok := inmemoryPostDB[id]
	if !ok {
		return Post{}, ErrNotFound
	}
	return post, nil
}

func (d *DB) GetAllPost(ctx context.Context) ([]Post, error) {
	posts := slices.SortedFunc(maps.Values(inmemoryPostDB), func(p1, p2 Post) int { return cmp.Compare(p1.ID, p2.ID) })
	return posts, nil
}

func (d *DB) UpdatePost(ctx context.Context, updatePost Post) (Post, error) {
	inmemoryPostDB[updatePost.ID] = updatePost
	return updatePost, nil
}

func (d *DB) DeletePostByID(ctx context.Context, id int) error {
	delete(inmemoryPostDB, id)

	return nil
}

type NewPostReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NewPostResp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type GetPostResp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type ListPostDataResp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePostReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePostResp struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewPostHandler(db *DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var newPostReq NewPostReq

		if err := c.ShouldBindJSON(&newPostReq); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		post, err := db.AddPost(c.Request.Context(), Post{
			Title: newPostReq.Title,
			Body:  newPostReq.Body,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		newPostResp := NewPostResp{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
		}

		c.JSON(http.StatusOK, newPostResp)
	}
}

func GetPostHandler(db *DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if idParam == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		post, err := db.GetPostByID(c.Request.Context(), id)
		if err != nil {
			if err == ErrNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		getPostResp := GetPostResp{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
		}
		c.JSON(http.StatusOK, getPostResp)
	}
}

func ListPostHanlder(db *DB) func(*gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.GetAllPost(c.Request.Context())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		listPostDataResps := make([]ListPostDataResp, 0, len(posts))
		for _, post := range posts {
			listPostDataResps = append(listPostDataResps, ListPostDataResp{
				ID:    post.ID,
				Title: post.Title,
				Body:  post.Body,
			})
		}

		c.JSON(http.StatusOK, listPostDataResps)
	}
}

func UpdatePostHanlder(db *DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if idParam == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		post, err := db.GetPostByID(c.Request.Context(), id)
		if err != nil {
			if err == ErrNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var updatePostReq UpdatePostReq

		if err := c.ShouldBindJSON(&updatePostReq); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		post.Body = updatePostReq.Body
		post.Title = updatePostReq.Title

		post, err = db.UpdatePost(c.Request.Context(), post)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		resp := UpdatePostResp{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func DeletePostHandler(db *DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if idParam == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		_, err = db.GetPostByID(c.Request.Context(), id)
		if err != nil {
			if err == ErrNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = db.DeletePostByID(c.Request.Context(), id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func main() {
	e := gin.Default()

	var db DB

	e.POST("/posts", NewPostHandler(&db))
	e.GET("/posts/:id", GetPostHandler(&db))
	e.GET("/posts", ListPostHanlder(&db))
	e.PATCH("/posts/:id", UpdatePostHanlder(&db))
	e.DELETE("/posts/:id", DeletePostHandler(&db))

	if err := e.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
