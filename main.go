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
)

const workMinutes = 20
const restMinutes = 5
const nbRounds = 2

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run() error {
	done := make(chan struct{})
	soundsCh := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	streamerWork, err := playBeep("./end_of_work_session.mp3")
	if err != nil {
		return err
	}
	streamerRest, err := playBeep("./end_of_rest_session.mp3")
	if err != nil {
		return err
	}
	streamerEnd, err := playBeep("./end_of_all_sessions.mp3")
	if err != nil {
		return err
	}

	defer func() {
		defer streamerWork.Close()
		defer streamerWork.Close()
		defer streamerEnd.Close()
	}()

	go func() {
		currentRoundNumber := 0
		for ; currentRoundNumber < nbRounds; currentRoundNumber++ {
			now := time.Now()
			workBar.Reset()
			workBar.RenderBlank()
		workFree:
			for {
				select {
				case <-time.Tick(1 * time.Minute):
					{
						count := minutesSince(now)
						workBar.Add(1)
						if count >= workMinutes {
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
				case <-time.Tick(1 * time.Minute):
					{
						count := minutesSince(now)
						restBar.Add(1)
						if count >= restMinutes {
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

		fmt.Printf("\n%d round have passed\n", currentRoundNumber)
	}()
	<-done

	speaker.Play(beep.Seq(streamerEnd, beep.Callback(func() {
		soundsCh <- struct{}{}
	})))
	<-soundsCh

	return nil
}

func minutesSince(t time.Time) int {
	// minutes := time.Since(t).Seconds()
	minutes := time.Since(t).Minutes()
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
