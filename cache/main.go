package main

import (
	"flag"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"strconv"
	"sync"
)

type Article struct {
	ID      int
	Content string
}

var (
	n       int
	cache   = make(map[int]*Article, 0)
	rwMutex sync.RWMutex
	sg      singleflight.Group
)

func init() {
	flag.IntVar(&n, "n", 5, "模拟的并发数，默认 5")
}

func main() {
	flag.Parse()

	wg := sync.WaitGroup{}

	wg.Add(n)
	// 模拟并发访问
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			article := fetchArticle(1)
			log.Println(article)
		}()
	}
	wg.Wait()
}

func fetchArticle(id int) *Article {
	article := findArticleFromCache(id)

	if article != nil && article.ID > 0 {
		return article
	}

	v, err, shared := sg.Do(strconv.Itoa(id), func() (interface{}, error) {
		return findArticleFromDB(id), nil
	})

	if err != nil {
		log.Println("singleflight do error:", err)
		return nil
	}

	// 打印 shared
	fmt.Println("shared===", shared)

	return v.(*Article)
}

// 模拟从缓存获取数据
func findArticleFromCache(id int) *Article {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return cache[id]
}

// 模拟从数据库中获取数据
func findArticleFromDB(id int) *Article {
	log.Printf("SELECT * FROM article WHERE id=%d", id)
	article := &Article{ID: id, Content: "xyz"}
	rwMutex.Lock()
	defer rwMutex.Unlock()
	cache[id] = article
	return article
}
