package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
    TOKEN   string
    PORT    int
    commands = []*discordgo.ApplicationCommand {
        {
            Name: "chicken",
            Description: "Posts an image of TF2 Scout turning into a chicken",
        },
        {
            Name: "sosig",
            Description: "sosig command",
        },
        {
            Name: "iplookup",
            Description: "posts information about a specific IP address",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Name: "ipaddress",
                    Description: "The address you would like to search for",
                    Type: discordgo.ApplicationCommandOptionString,
                    Required: true,
                },
            },
        },
        {
            Name: "fumo",
            Description: "posts a random fumo image",
        },
        {
            Name: "resolve",
            Description: "closes the support thread",
        },
    }

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "chicken": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            // Posts an image of TF2 Scout turning into a chicken
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/view/chicken-gif-19565842",
				},
			})
        },
        "sosig": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            // Posts an image of a man doing "tricks" with a sausage
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/view/sosig-gif-23013003",
					
				},
			})
        },
        "iplookup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            // gets ipaddress from interaction
            ipaddress := i.ApplicationCommandData().Options[0].StringValue()
            
            // Make API call for JSON data
            resp, err := http.Get("http://ipwhois.app/json/" + strings.ToLower(ipaddress))
            handlePanic(err)

            defer resp.Body.Close()

            // On successful API call
            if resp.StatusCode == 200 {
                type Response struct {
                    IP        string `json:"ip"`
                    ISP       string `json:"isp"`
                    Country   string `json:"country"`
                    Region    string `json:"region"`
                    City      string `json:"city"`
                    Timezone  string `json:"timezone"`
                    GMTOffset string `json:"timezone_gmt"`
                }

                body, _ := ioutil.ReadAll(resp.Body)
                var f Response
                err := json.Unmarshal(body, &f)
                handlePanic(err)

                // Generate Discord embed
                embed := &discordgo.MessageEmbed {
                    Color:       0xff1100, // Red
                    Title:       "IP lookup results for " + strings.ToLower(ipaddress),
                    Description: fmt.Sprintf("IP: %s\nISP: %s\nCountry: %s\nRegion: %s\nCity: %s\nTimezone: %s\nGMT-Offset: %s", f.IP, f.ISP, f.Country, f.Region, f.City, f.Timezone, f.GMTOffset),
                }

                // Send Discord embed
                s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData{
                        Embeds: []*discordgo.MessageEmbed {
                            embed, 
                        }, 
                    },
                })
            } else {
                log.Println("Received HTTP status code:", resp.StatusCode)
            }
        },
        "fumo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            // Make API call for JSON data
            resp, err := http.Get("https://fumoapi.nosesisaid.me/random")
            handlePanic(err)

            defer resp.Body.Close()

            // On successful API call
            if resp.StatusCode == 200 {
                type Response struct {
                    URL     string  `json:"url"`
                }

                body, _ := ioutil.ReadAll(resp.Body)
                var f Response
                err := json.Unmarshal(body, &f)
                handlePanic(err)

                // Send random image of Fumo
                s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData{
                        Content: f.URL,
                    },
                })
            } else {
                log.Println("Received HTTP status code:", resp.StatusCode)
            }
        },
        "resolve": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            // Set channel variable
            channel, err := s.Channel(i.ChannelID)
            handlePanic(err)

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "closing thread",
					
				},
			})

            // Define support channel ID
            SUPPORT_CHANNEL_ID := os.Getenv("SUPPORT_CHANNEL_ID")

            fmt.Printf("Channel ID: %s\nSupport Channel ID: %s\n", i.ChannelID, SUPPORT_CHANNEL_ID)
            // Set thread to archived
            if channel.IsThread() {
                s.ChannelEditComplex(i.ChannelID, &discordgo.ChannelEdit{Archived: true})
                log.Printf("Archived thread with ID: %s", i.ChannelID)
            }
        },
    }


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

    dcSession.Open()

    dcSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate){
        if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
            h(s,i)
        }
    })

    _, err := dcSession.ApplicationCommandBulkOverwrite(dcSession.State.User.ID, "621749786722172930", commands)
    handlePanic(err)
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
    TOKEN := os.Getenv("DISCORD_TOKEN")
    handlePanic(err)

    // Create a new Discord session using the provided token
    dcSession, err := discordgo.New("Bot " + TOKEN)
    handlePanic(err)

    dcSession.Identify.Presence.Status = string(discordgo.StatusOnline)
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
    if m.Author.ID == s.State.User.ID || m.Author.Bot == true {
        return
    }

    // Define support channel ID
    SUPPORT_CHANNEL_ID := os.Getenv("SUPPORT_CHANNEL_ID")

    // Check if message was sent in the designated support channel
    if m.ChannelID == SUPPORT_CHANNEL_ID {
        // Create thread based off message sender
        thread, err := s.MessageThreadStart(m.ChannelID, m.ID, m.Author.Username, 1440)
        handlePanic(err)

        log.Printf("Created thread with ID: %s", m.ChannelID)

        // Send messages on thread creation
        s.ChannelMessageSend(thread.ID, "Hi there! I have created this support thread for you.")
        s.ChannelMessageSend(thread.ID, "If you no longer need assistance, please use the `/resolve` command to archive this thread.")
    }
}

func handlePanic(err error) {
    if err != nil {
        panic(err)
    }
}
