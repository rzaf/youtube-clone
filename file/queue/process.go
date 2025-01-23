package queue

import (
	"fmt"
	"time"

	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/file/models"
)

func RunQueue() {
	time.Sleep(3500 * time.Millisecond)
	fmt.Printf("starting queue \n")
	processes, err := models.GetRemainingProcesses()
	fmt.Printf("err: %v\n", err)
	fmt.Printf("remaining processes length: %v\n", len(processes))
	if err == nil {
		queueMutext.Lock()
		for i := range processes {
			fmt.Printf("remaining procesing[%d]: %v\n", i, processes[i])
			mediaType := pbHelper.MediaType_VIDEO
			switch processes[i].Name {
			case "video_format":
				mediaType = pbHelper.MediaType_VIDEO
			case "photo_format":
				mediaType = pbHelper.MediaType_PHOTO
			case "music_format":
				mediaType = pbHelper.MediaType_MUSIC
			}
			queue = append(queue, NewFileFormat(processes[i].Url, mediaType))
			size++
		}
		queueMutext.Unlock()
	}
	for {
		if size != 0 {
			e := Top()
			fmt.Printf("starting procces for :'%s'\n", e.url())
			err := e.process()
			fmt.Printf("procces:'%s' finished . err:%v\n", e.url(), err)
			if err == nil {
				err2 := models.SetProcessDone(e.name(), e.url())
				fmt.Printf("err2 %v \n", err2)
			} else {
				err2 := models.SetProcessFailed(e.url())
				fmt.Printf("err2 %v \n", err2)
			}
			pop()
		}
		time.Sleep(500 * time.Millisecond)
	}
}
