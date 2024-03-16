package main

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// Current work session progressbar
var workBar = progressbar.NewOptions(workMinutes,
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

// Current rest session progressbar
var restBar = progressbar.NewOptions(restMinutes,
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
