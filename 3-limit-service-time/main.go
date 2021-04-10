//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	m         sync.Mutex
}

const secondsLimit = 10

func (u *User) incTime(value int64) {
	u.m.Lock()
	defer u.m.Unlock()
	u.TimeUsed += value
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed

func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	if u.TimeUsed >= secondsLimit {
		return false
	}
	done := make(chan bool)
	go func() {
		process()
		done <- true
	}()
	for {
		if u.TimeUsed >= secondsLimit {
			return false
		}
		select {
		case <-time.After(time.Second):
			u.incTime(1)
			continue
		case usedSeconds := <-done:
			fmt.Printf("Done task, time - %v \n", usedSeconds)

			return true
		}

	}

}

func main() {
	RunMockServer()
}
