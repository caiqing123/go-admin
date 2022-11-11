package site

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"api/pkg/app"
	"api/pkg/book/store"
)

type SyncStore struct {
	lock  *sync.Mutex
	jobs  [][]bool
	Store *store.Store
}

func (s *SyncStore) Init() {
	s.jobs = make([][]bool, len(s.Store.Volumes))
	s.lock = &sync.Mutex{}
	for k := range s.jobs {
		s.jobs[k] = make([]bool, len(s.Store.Volumes[k].Chapters))
	}
}

func (s *SyncStore) GetJob() (vi, ci int, url string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for vi, vol := range s.Store.Volumes {
		for ci, ch := range vol.Chapters {
			if !s.jobs[vi][ci] {
				if len(ch.Text) == 0 {
					s.jobs[vi][ci] = true
					return vi, ci, ch.URL, nil
				}
			}
		}
	}
	return 0, 0, "", io.EOF
}

func Job(syncStore *SyncStore, jobch chan error) {

	var (
		content []string
	)

	defer func(jobch chan error) {
		jobch <- io.EOF
	}(jobch)

	for {
		vi, ci, BookURL, err := syncStore.GetJob()
		if err != nil {
			if err != io.EOF {
				jobch <- err
			}
			return
		}

	A:
		for i := 0; i < 3; i++ { //重试次数
			content, err = Chapter(BookURL)
			if err != nil {
				log.Printf("Error: %s %s", err, BookURL)
				time.Sleep(500 * time.Millisecond) //爬取错误间隔
				continue A
			}
			syncStore.SaveJob(vi, ci, content)
			jobch <- nil
			time.Sleep(200 * time.Millisecond) //爬取间隔
			break A
		}
	}
}

func (s *SyncStore) SaveJob(vi, ci int, text []string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Store.Volumes[vi].Chapters[ci].Text = text
	s.Store.LastUpdate = time.Now()
	//log.Printf("SaveJob，%s-%s", s.Store.BookName, s.Store.Volumes[vi].Chapters[ci].Name)
}

func Download(chapter *store.Store) {
	ssss := &SyncStore{
		Store: chapter,
	}
	ssss.Init()
	var threadNum = 10
	var chCount = 0
	var isDone = 0
	for _, v := range chapter.Volumes {
		chCount += len(v.Chapters)
		for _, v2 := range v.Chapters {
			if len(v2.Text) != 0 {
				isDone++
			}
		}
	}
	if isDone != 0 {
		log.Printf("[读入] 已缓存:%d", isDone)
	}

	// End Print
	defer func(s *store.Store) {
		var chCount = 0
		var isDone = 0
		for _, v := range chapter.Volumes {
			chCount += len(v.Chapters)
			for _, v2 := range v.Chapters {
				if len(v2.Text) != 0 {
					isDone++
				}
			}
		}
		if isDone != 0 {
			log.Printf("[爬取结束] 已缓存:%d", isDone)
		}
	}(chapter)

	if isDone < chCount {

		Jobch := make(chan error)
		for i := 0; i < threadNum; i++ {
			go Job(ssss, Jobch)
		}

		var ii = 0
	AA:
		for {
			select {
			case err := <-Jobch:
				if err != nil {
					if err == io.EOF {
						ii++
						if ii >= threadNum {
							log.Printf("缓存完成")
							break AA
						}
					} else {
						log.Printf("Job Error: %s", err)
					}
				}
			}
		}
		close(Jobch)
	}
}

func DownloadWs(chapter *store.Store, ctx context.Context, id string, group string, hookfn func(context.Context, string, string, []byte), path string) {
	ssss := &SyncStore{
		Store: chapter,
	}
	ssss.Init()
	var threadNum = 10
	var chCount = 0
	var isDone = 0
	for _, v := range chapter.Volumes {
		chCount += len(v.Chapters)
		for _, v2 := range v.Chapters {
			if len(v2.Text) != 0 {
				isDone++
			}
		}
	}

	ext := ".txt"
	if strings.Contains(chapter.BookURL, "bookstack.cn") {
		ext = ".epub"
	}

	src := `{"progress":"%v","type":"progress","book_id":"` + chapter.BookName + "_" + id + `","download_url":"` + app.URL(path+ext) + `"}`

	if isDone != 0 {
		hookfn(ctx, id, group, []byte(fmt.Sprintf(src, "100")))
		log.Printf("[读入] 已缓存:%d", isDone)
	}

	// End Print
	defer func(s *store.Store) {
		var chCount = 0
		var isDone = 0
		for _, v := range chapter.Volumes {
			chCount += len(v.Chapters)
			for _, v2 := range v.Chapters {
				if len(v2.Text) != 0 {
					isDone++
				}
			}
		}
		if isDone != 0 {
			log.Printf("[爬取结束] 已缓存:%d", isDone)
		}
	}(chapter)

	if isDone < chCount {
		progress := 0
		hookfn(ctx, id, group, []byte(fmt.Sprintf(src, "0")))
		Jobch := make(chan error)
		for i := 0; i < threadNum; i++ {
			go Job(ssss, Jobch)
		}

		var ii = 0
	AA:
		for {
			select {
			case err := <-Jobch:
				if err != nil {
					if err == io.EOF {
						ii++
						if ii >= threadNum {
							hookfn(ctx, id, group, []byte(fmt.Sprintf(src, "100")))
							log.Printf("缓存完成")
							break AA
						}
					} else {
						log.Printf("Job Error: %s", err)
					}
				} else {
					progress++
					go hookfn(ctx, id, group, []byte(fmt.Sprintf(src, strconv.FormatFloat(float64(float32(progress)/float32(len(chapter.Volumes[0].Chapters))*100), 'f', 10, 32))))
				}
			}
		}
		close(Jobch)
	}
}
