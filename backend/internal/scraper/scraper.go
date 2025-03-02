package scraper

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func connectPostgres(user, password, address, dbName string) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, address, dbName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// add query logging hook
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.RegisterModel((*models.Topic)(nil))
	db.RegisterModel((*models.Voxsphere)(nil))
	db.RegisterModel((*models.Award)(nil))
	db.RegisterModel((*models.User)(nil))
	db.RegisterModel((*models.Trophy)(nil))
	db.RegisterModel((*models.Post)(nil))
	db.RegisterModel((*models.Emoji)(nil))
	db.RegisterModel((*models.CustomEmoji)(nil))
	db.RegisterModel((*models.UserTrophy)(nil))
	db.RegisterModel((*models.VoxsphereMember)(nil))
	db.RegisterModel((*models.VoxsphereModerator)(nil))
	db.RegisterModel((*models.UserFlair)(nil))
	db.RegisterModel((*models.UserFlairCustomEmoji)(nil))
	db.RegisterModel((*models.UserFlairEmoji)(nil))
	db.RegisterModel((*models.UserFlairDescription)(nil))
	db.RegisterModel((*models.PostFlair)(nil))
	db.RegisterModel((*models.PostFlairCustomEmoji)(nil))
	db.RegisterModel((*models.PostFlairEmoji)(nil))
	db.RegisterModel((*models.PostFlairDescription)(nil))
	db.RegisterModel((*models.PostAward)(nil))
	db.RegisterModel((*models.Post)(nil))
	db.RegisterModel((*models.PostMedia)(nil))
	db.RegisterModel((*models.Image)(nil))
	db.RegisterModel((*models.ImageMetadata)(nil))
	db.RegisterModel((*models.Gif)(nil))
	db.RegisterModel((*models.GifMetadata)(nil))
	db.RegisterModel((*models.Gallery)(nil))
	db.RegisterModel((*models.GalleryMetadata)(nil))
	db.RegisterModel((*models.Video)(nil))
	db.RegisterModel((*models.Link)(nil))

	// drop all rows of the user,trophies table
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Topic)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Voxsphere)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Award)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.User)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Trophy)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Emoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.CustomEmoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserTrophy)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.VoxsphereMember)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.VoxsphereModerator)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlair)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairCustomEmoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairEmoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.UserFlairDescription)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlair)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairCustomEmoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairEmoji)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostFlairDescription)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostAward)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Post)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.PostMedia)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Image)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.ImageMetadata)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Gif)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.GifMetadata)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Gallery)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.GalleryMetadata)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Video)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}
	if _, err := db.NewTruncateTable().Cascade().Model((*models.Link)(nil)).Exec(context.Background()); err != nil {
		fmt.Println("truncate table failed:", err)
	}

	return db, nil
}

func Run() {
	// get db instance
	db, err := connectPostgres("postgres", "postgres", "127.0.0.1:5432", "voxpopuli")
	if err != nil {
		panic(err)
	}
	_ = db
}
