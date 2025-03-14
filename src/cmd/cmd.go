package cmd

import (
	"context"
	"flag"
	"fmt"
	"ratblog/config"
	"ratblog/contract"
	"ratblog/internal/entity"
	"ratblog/internal/repository"
	serviceV1 "ratblog/internal/services/v1"
)

var (
	createTags = flag.Bool("create-tags", false, "Create tags raw data if not exists")

	cfg = config.LoadConfig("config/")
)

func createDefulteTagsIfNotExists(ctx context.Context, repo contract.Repository) {
	tags := []string{
		"Go",
		"Java",
		"Python",
		"JavaScript",
		"Ruby",
		"PHP",
		"Swift",
		"Kotlin",
		"Rust",
		"Scala",
		"TypeScript",
	}

	for _, tag := range tags {
		_ = repo.CreateTag(ctx, &entity.Tag{
			Name: tag,
		})
	}
}

func Run() error {
	flag.Parse()

	ctx := context.TODO()
	postgreeUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	repo := repository.New(ctx, postgreeUrl)

	if *createTags {
		createDefulteTagsIfNotExists(ctx, repo)
	}

	echoService := serviceV1.SetupHttpServerV1(ctx, repo, cfg)
	addres := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	if err := echoService.Start(addres); err != nil {
		return err
	}

	return nil
}
