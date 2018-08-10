package service

import "time"

func RunTimeReport(f func(), delay time.Duration) {
	for ; ; {
		f()
		time.Sleep(delay)
	}
}

func main() {
	go RunTimeReport(func() {
		println(time.Now().Format("2006-01-02 15:04:05"))
	}, 1*time.Second)
	time.Sleep(1 * time.Hour)
}
