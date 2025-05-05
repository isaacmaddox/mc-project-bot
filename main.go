package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/cmd"
	"github.com/isaacmaddox/mc-project-bot/cmd/project"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var discord *discordgo.Session

var (
	TOKEN string
	GUILD_ID       = *flag.String("guildId", "1367938202349863023", "the unique id for the server the bot runs in")
	RemoveCommands = flag.Bool("remove", false, "unregisters all commands when the program is stopped")
)

func init() {
	util.ErrorCheck(godotenv.Load(), "Could not load env: %v")
	TOKEN = os.Getenv("BOT_TOKEN")
}

var commands = []*discordgo.ApplicationCommand{
	// cmd.HiCommand,
	project.ProjectCommand,
	cmd.TestCommand,
	cmd.ClearCommand,
}

var commandHandlers = map[string]func(discord *discordgo.Session, i *discordgo.InteractionCreate){
	// "hi": cmd.HiHandler,
	"project": project.ProjectHandler,
	"test": cmd.TestHandler,
	"clear": cmd.ClearHandler,
}

func init() {
	flag.Parse()

	discord = util.Extract(discordgo.New("Bot " + TOKEN))

	discord.AddHandler(func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(discord, i)
		}
	})
}

func main() {
	db.Init_database()
	
	discord.AddHandler(func(discord *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", discord.State.User.Username, discord.State.User.ID)
	})

	util.ErrorCheck(discord.Open(), "Couldn't start bot: %v")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	log.Println("Adding commands...")

	for i, v := range commands {
		log.Printf("Adding the %v command\n", v.Name)
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, GUILD_ID, v)

		if err != nil {
			log.Fatalf("Couldn't add the %v command: %v", v.Name, err)
		}

		registeredCommands[i] = cmd
	}

	defer discord.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press CTRL+C to stop")
	<-stop

	log.Println("Shutting down...")

	if *RemoveCommands {
		log.Println("Unregistering commands...")

		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, GUILD_ID, v.ID)

			if err != nil {
				log.Fatalf("Couldn't delete the %v command: %v", v.Name, err)
			}
		}
	}
}