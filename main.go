package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

const workMinutes = 5
const restMinutes = 5
const nbRounds = 2

func main() {
	workBar := progressbar.NewOptions(workMinutes,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("[red][Work session progress][reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[red]=[reset]",
			SaucerHead:    "[red]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	restBar := progressbar.NewOptions(restMinutes,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("[blue][Take a break, you deserve it][reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]=[reset]",
			SaucerHead:    "[blue]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	done := make(chan struct{})
	soundsCh := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	streamerWork, err := playBeep("./end_of_work_session.mp3")
	if err != nil {

	}
	defer streamerWork.Close()
	streamerRest, err := playBeep("./end_of_rest_session.mp3")
	if err != nil {

	}
	streamerEnd, err := playBeep("./end_of_all_sessions.mp3")
	if err != nil {

	}
	defer streamerEnd.Close()

	go func() {
		for i := 0; i < nbRounds; i++ {
			now := time.Now()
			workBar.Reset()
			workBar.RenderBlank()
		workFree:
			for {
				select {
				case <-time.Tick(1 * time.Second):
					{
						count := minutesSince(now)
						workBar.Add(1)
						if count > workMinutes {
							break workFree
						}
					}
				case <-sigs:
					{
						done <- struct{}{}
						return
					}
				}
			}
			speaker.Play(beep.Seq(streamerWork, beep.Callback(func() {
				soundsCh <- struct{}{}
			})))
			<-soundsCh
			streamerWork.Seek(0)
			now = time.Now()
			restBar.Reset()
			restBar.RenderBlank()
		restFree:
			for {
				select {
				case <-time.Tick(1 * time.Second):
					{
						count := minutesSince(now)
						restBar.Add(1)
						if count > restMinutes {
							break restFree
						}
					}
				case <-sigs:
					{
						done <- struct{}{}
						return
					}
				}
			}
			//make a sound here
			speaker.Play(beep.Seq(streamerRest, beep.Callback(func() {
				soundsCh <- struct{}{}
			})))
			streamerRest.Seek(0)
			<-soundsCh
			if err != nil {
				continue
			}
			fmt.Print("\033[H\033[2J")
		}
		done <- struct{}{}
	}()
	<-done

	speaker.Play(beep.Seq(streamerEnd, beep.Callback(func() {
		soundsCh <- struct{}{}
	})))
	<-soundsCh

	fmt.Printf("\n%d round have passed\n", nbRounds)
}

func minutesSince(t time.Time) int {
	minutes := time.Since(t).Seconds()
	// minutes := time.Since(t).Minutes()
	return int(minutes)
}

func playBeep(fileName string) (beep.StreamSeekCloser, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return nil, err
	}
	return streamer, nil
}
