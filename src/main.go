package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	// "strings"
	"time"
	"flag"
	// "fmt"

	// External Packages
	"github.com/sevlyar/go-daemon"
	"github.com/fatih/color"
)


// CONFIGURE GLOBAL STD OUTPUT COLORS
var (
	errOut  = color.New(color.FgRed).Add(color.Bold)
	infoOut = color.New(color.FgHiMagenta)
	stdOut  = color.New()
)


type commandline_arguments struct{
	 remind string
	 title string
	 summary string
	 icon string
	 urgency int
	 daemon bool
}
func parseInput() commandline_arguments{
	var FlagRemind=flag.String("Remind", "", "Time to Remind")
	var FlagTitle=flag.String("Title", "", "Message for title")
	var FlagSummary=flag.String("Summary", "", "Message for summary")
	var FlagIcon=flag.String("Icon", "", "Custom Icon to use")
	var FlagUrgent=flag.Int("Urgency", 1, "Set urgancy level")
	var FlagBool=flag.Bool("Daemon", false, "Daemonize process or not")
	flag.Parse()
	flags := commandline_arguments{*FlagRemind, *FlagTitle,*FlagSummary,*FlagIcon,*FlagUrgent,*FlagBool}
	return flags
}
// Prints Help Menu
func printHelp() {
	cyan := color.New(color.FgHiCyan).SprintFunc()

	errOut.Println("Two Arguments Required:")
	infoOut.Println("\tArgument 1: [Time {amount(s/m/h)}]")
	infoOut.Println("\tArgument 2: [Message]")

	infoOut.Println("Examples: ")
	stdOut.Printf("\tapp %s \n", cyan("{time[s/m/h]} {message}"))
	stdOut.Printf("\tapp %s \n", cyan("2s \"Hello World\""))
	stdOut.Printf("\tapp %s \n", cyan("2 \"Hello World\""))
}

// Simple wrapper that returns the conversion of string to int
func getIntStr(sVal string) int {
	intVal, err := strconv.Atoi(sVal)
	if err != nil {
		errOut.Println("First Argument is time to Sleep! [int]")
		os.Exit(1)
	}
	return intVal
}

func main() {
	var args = parseInput()
	// Icon := args.icon
	Remind := args.remind
	// Summary := args.summary
	Title := args.title
	// Urgency := args.urgency
	isDaemon := args.daemon

	// VERIFY ARGUMENTS 
	// Title and Reminder must be enabled
	if(Title ==""){
		printHelp()
		errOut.Println("Title of notification is not set")
		os.Exit(-1)
	} else if (Remind == "") {
		printHelp()
		errOut.Println("Reminder time is not set")
		os.Exit(-1)
	}
	// DETERMINE SLEEP AMOUNT
	var dTime time.Duration
	tTypeStr := Remind[len(os.Args[1])-1]
	waitTime := Remind[:len(os.Args[1])-1]
	waitType := "Seconds"

	switch tTypeStr {
	case 's': // Specifically Seconds
		dTime = time.Duration(getIntStr(waitTime)) * time.Second
	case 'm': // Minutes
		dTime = time.Duration(getIntStr(waitTime)) * time.Minute
		waitType = "Minutes"
	case 'h': // Hours
		dTime = time.Duration(getIntStr(waitTime)) * time.Hour
		waitType = "Hours"
	default: // Defaulted to Seconds
		waitTime = os.Args[1]
		dTime = time.Duration(getIntStr(waitTime)) * time.Second
	}
	infoOut.Printf("Waiting for %s %s to output '%s'\n", waitTime, waitType, os.Args[2])

	// Obtain Icon Full Path
	binPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	iconPath := binPath + "/Notification.png"

	// Deamonize if Flag
	if isDaemon {
		infoOut.Println("Daemonized Process, running in the Background 😈")

		// Setup Daemon
		ctx := &daemon.Context{
			PidFileName: binPath + "/timed-notify.pid",
			PidFilePerm: 0644,
			LogFileName: binPath + "/timed-notify.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        os.Args,
		}

		// Release the DAEMON!
		d, err := ctx.Reborn()
		if err != nil {
			errOut.Printf("Unable to run: %s\n", err)
		}
		if d != nil {
			os.Exit(0)
		}
		ctx.Release()
	}

	// SET SLEEP TIME
	time.Sleep(dTime)

	// INITIATE NOTIFICATION
	cmd := exec.Command("notify-send", os.Args[2], "-i", iconPath)
	cmd.Start()
}
