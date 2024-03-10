package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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
			fmt.Print("\033[H\033[2J") //clear the console
		}
		done <- struct{}{}
	}()

	// go func() {
	// 	<-sigs
	// 	done <- struct{}{}
	// }()
	<-done

	fmt.Printf("\n %d round have passed\n", nbRounds)
}

func minutesSince(t time.Time) int {
	minutes := time.Since(t).Seconds()
	// minutes := time.Since(t).Minutes()
	return int(minutes)
}
