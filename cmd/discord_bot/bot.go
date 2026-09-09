package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lbenedar/fresta/cmd/api/core"
	"github.com/lbenedar/fresta/cmd/discord_bot/config"
	"github.com/lbenedar/fresta/proto/grpc_bots"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DiscordBot struct {
	app       *core.Application
	dgSession *discordgo.Session
	cfg       *config.Config
}

func CreateNewBot(app *core.Application) (*DiscordBot, error) {
	apiKey := os.Getenv("DISCORD_API_KEY")
	if apiKey == "" {
		return nil, ErrorTokenNotExist
	}

	guildId := os.Getenv("GUILD_ID")
	if guildId == "" {
		return nil, ErrorGuildIdNotExist
	}

	bot := &DiscordBot{app: app}

	cfg, err := config.ParseConfigs()
	if err != nil {
		return nil, err
	}

	cfg.GuildId = guildId

	s, err := discordgo.New("Bot " + apiKey)
	if err != nil {
		return nil, err
	}

	bot.cfg = cfg
	bot.dgSession = s

	return bot, nil
}

func (b *DiscordBot) ListenForConfigChanges() {
	b.app.Slogger.Info("Listener for config files has been set")
	for {
		time.Sleep(5 * time.Second)
		changedFiles, err := config.CheckIfFilesChanged()
		if err != nil {
			b.app.Slogger.Info("Got error on check file change", "err", err)
			continue
		}
		if len(changedFiles) == 0 {
			continue
		}
		b.app.Slogger.Info("Files has been changed", "files", changedFiles)

		changed := config.UpdateConfigsByNames(b.cfg, changedFiles)
		if changed {
			b.app.Slogger.Info("Data", "data", b.cfg.Commands)
			b.dgSession.ApplicationCommandBulkOverwrite(b.dgSession.State.User.ID, "b.cfg.GuildId", b.cfg.Commands)
			b.app.Slogger.Info("Commands has been updated")
		}
	}
}

func (b *DiscordBot) ListenForShutdown(closeChan chan any, errChan chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	b.app.Slogger.Info("Caught signal", "signal", s.String())

	err := b.dgSession.Close()
	if err != nil {
		errChan <- err
	}

	err = b.app.FoundryApp.Shutdown()
	if err != nil {
		errChan <- err
	}

	b.app.FoundryApp.Transport.DB.Close()

	closeChan <- nil
}

func (b *DiscordBot) ListenAndServe() error {
	setupHandlers(b)

	if err := b.dgSession.Open(); err != nil {
		return err
	}
	defer b.dgSession.Close()

	err := config.InitCommands(b.dgSession, b.cfg)
	if err != nil {
		return err
	}

	closeChan := make(chan any)
	errChan := make(chan error)

	go b.ListenForConfigChanges()
	go b.ListenForShutdown(closeChan, errChan)

	for {
		select {
		case <-closeChan:
			return nil
		case err := <-errChan:
			return err
		case err := <-b.cfg.Errs:
			b.app.Slogger.Error("Received error in config", "err", err.Error())
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
		return
	}

	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		log.Fatalf("Ошибка! Отсутствует API HOST в переменных окружения: %v", err)
		return
	}

	conn, err := grpc.NewClient(apiHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
		return
	}
	defer conn.Close()

	client := grpc_bots.NewBotsGatewayClient(conn)

	// Отправляем запрос
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetStatus(ctx, &grpc_bots.GetStatusRequest{Type: grpc_bots.StatusType_STATUS_DATA})
	if err != nil {
		log.Fatalf("Ошибка при вызове GetStatus: %v", err)
	}

	log.Printf("Получены данные о статусе: Success=%s, Data=%s", res.Success, res.Data)

}
