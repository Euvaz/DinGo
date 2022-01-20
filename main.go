package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
    "log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	TOKEN   string
    PORT    int
)

func init() {
    // Initialize environment
    godotenv.Load()
}

func main() {
    PORT := os.Getenv("PORT")

    if PORT == "" {
        log.Println("$PORT not found within environment")
    }
    
    log.Println("Starting")

    dcSession := initializeDiscord()
    log.Println("Discord handlers added")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
    handlePanic(dcSession.Close())
}

// General types and functions
func initializeLogging() *os.File {
	log.SetFlags(log.Ltime | log.Lshortfile)
	logFilename := "logs.txt"
	if os.Getenv("DYNO") != "" {
		logFilename = "/tmp/logs.txt"
	}
	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	return file
}

func initializeDiscord() *discordgo.Session {
    var err error

    // Load Discord Token
    TOKEN:= os.Getenv("DISCORD_TOKEN")
    handlePanic(err)

    // Create a new Discord session using the provided token
    dcSession, err := discordgo.New("Bot " + TOKEN)
    handlePanic(err)

    dcSession.AddHandler(dcOnMessageCreate)

    // Open a websocket connection to Discord and begin listening
    err = dcSession.Open()
    handlePanic(err)
    
	dcSession.Identify.Intents = discordgo.IntentsGuildMessages

    return dcSession
}

func dcOnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
    
    // Define bot prefix
    PREFIX := os.Getenv("PREFIX")

    // Check for command prefix
    if strings.HasPrefix(m.Content, PREFIX) {

        // Switch case for command functions
        switch command := strings.Fields(m.Content); strings.ToLower(command[1]) {
        case "chicken":
            dcCommandChicken(command, s, m)
        case "sosig":
            dcCommandSosig(command, s, m)
        case "iplookup":
            dcCommandIPLookup(command, s, m)
        default:
            s.ChannelMessageSend(m.ChannelID, "Command not found.")
        }
    }
}

func dcCommandChicken(command []string, s *discordgo.Session, m *discordgo.MessageCreate) {
    // Posts an image of TF2 Scout turning into a chicken
    s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/chicken-gif-19565842")
}

func dcCommandSosig(command []string, s *discordgo.Session, m *discordgo.MessageCreate) {
    // Posts an image of a man doing "tricks" with a sausage
    s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/sosig-gif-23013003")
}

func dcCommandIPLookup(command []string, s *discordgo.Session, m *discordgo.MessageCreate) {
    // Prompt for IP address/hostname if not provided
    if len(command) <= 2 {
        s.ChannelMessageSend(m.ChannelID, "Please provide an IP address/hostname.")
        return
    }
    // Make API call for JSON data
    resp, err := http.Get("http://ipwhois.app/json/" + command[2])
    handlePanic(err)
    
    defer resp.Body.Close()

    // On successful API call
    if resp.StatusCode == 200 {
        type Response struct {
            ISP         string  `json:"isp"`
            Country     string  `json:"country"`
            Region      string  `json:"region"`
            City        string  `json:"city"`
            Timezone    string  `json:"timezone"`
            GMTOffset   string  `json:"timezone_gmt"`
        } 
                
        body, _ := ioutil.ReadAll(resp.Body)
        var f Response
        err := json.Unmarshal(body, &f)
        handlePanic(err)

        // Generate Discord embed
        embed := &discordgo.MessageEmbed {
            Color: 0xff1100, // Red
            Title: "IP lookup results for " + command[2],
            Description: fmt.Sprintf("ISP: %s\nCountry: %s\nRegion: %s\nCity: %s\nTimezone: %s\nGMT-Offset: %s", f.ISP, f.Country, f.Region, f.City, f.Timezone, f.GMTOffset),
        }

        // Send Discord embed
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    } else {
        log.Println("Received HTTP status code:", resp.StatusCode)
    }
}

func handlePanic(err error) {
    if err != nil {
        panic(err)
    }
}
