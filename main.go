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

    // Load dotenv into environment
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
    // dcSession.AddHandler(dcOnChannelDelete)

    // Open a websocket connection to Discord and begin listening
    err = dcSession.Open()
    handlePanic(err)
    
	dcSession.Identify.Intents = discordgo.IntentsGuildMessages

    //guild, err = dcSession.Guild(dcSession.GuildID)
	//handlePanic(err)

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

        // Divide message
        message := strings.Fields(m.Content)

        // Switch case for command functions
        switch message[1] {
        case "chicken":
            s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/chicken-gif-19565842")
        case "iplookup":

            // Prompt for IP address/hostname if not provided
            if len(message) <= 2 {
                s.ChannelMessageSend(m.ChannelID, "Please provide an IP address/hostname.")
                return
            }

            // Make API call for JSON data
            resp, err := http.Get("http://ipwhois.app/json/" + message[2])

            if err != nil {
                log.Println(err)
            }
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
                if err != nil {
                    log.Println(err)
                }

                if err != nil {
                    log.Println(err)
                }

                embed := &discordgo.MessageEmbed {
                    Color: 0xff1100, // Red
                    Title: "IP lookup results for " + message[2],
                    Description: fmt.Sprintf("ISP: %s\nCountry: %s\nRegion: %s\nCity: %s\nTimezone: %s\nGMT-Offset: %s", f.ISP, f.Country, f.Region, f.City, f.Timezone, f.GMTOffset),
                }

                s.ChannelMessageSendEmbed(m.ChannelID, embed)

            } else {
                log.Println("Received HTTP status code:", resp.StatusCode)
            }
        default:
            s.ChannelMessageSend(m.ChannelID, "Command not found.")
        }
    }
}

func handlePanic(err error) {
    if err != nil {
        panic(err)
    }
}
