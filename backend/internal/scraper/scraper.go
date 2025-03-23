package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/forPelevin/gomoji"
	"github.com/glowfi/voxpopuli/backend/internal/helper"
	"github.com/glowfi/voxpopuli/backend/internal/threadsafe"
	"github.com/glowfi/voxpopuli/backend/pkg/models"
	awardsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/award"
	commentsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/comments"
	customemojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/custom_emoji"
	emojirepo "github.com/glowfi/voxpopuli/backend/pkg/repo/emoji"
	mediasrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/media"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postflairsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post_flair"
	relationrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/relation"
	rulesrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/rule"
	topicsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/topic"
	trophiesrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/trophy"
	"github.com/glowfi/voxpopuli/backend/pkg/repo/user"
	usersrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user"
	userflairsrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/user_flair"
	voxspheresrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/voxsphere"
	"github.com/google/uuid"
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
	db.RegisterModel((*models.UserUserFlair)(nil))
	db.RegisterModel((*models.PostFlair)(nil))
	db.RegisterModel((*models.PostFlairCustomEmoji)(nil))
	db.RegisterModel((*models.PostFlairEmoji)(nil))
	db.RegisterModel((*models.PostFlairDescription)(nil))
	db.RegisterModel((*models.PostPostFlair)(nil))
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

func insertTopics(ctx context.Context, db *bun.DB, filename string) error {
	var topicsJson map[string][]string

	// Open the file
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &topicsJson); err != nil {
		return err
	}

	topicsExpectedTotalRecords := 0
	for _, topics := range topicsJson {
		topicsExpectedTotalRecords += len(topics)
	}

	topicCI := NewConcurrentInserter(topicsExpectedTotalRecords, func(topics []models.Topic) error {
		topicRepo := topicsrepo.NewRepo(db)

		if _, err := topicRepo.AddTopics(ctx, topics...); err != nil {
			return err
		}

		return nil
	}, 300*time.Millisecond, 500)

	wg := new(sync.WaitGroup)
	wg.Add(topicsExpectedTotalRecords)

	go func() {
		for parentTopic, topicsList := range topicsJson {
			for _, topic := range topicsList {
				newTopic := models.Topic{
					ID:       uuid.New(),
					Name:     topic,
					Category: parentTopic,
				}
				topicCI.ResC <- newTopic
				wg.Done()
			}
		}
	}()

	if err := topicCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func insertTrophies(ctx context.Context, db *bun.DB, filename string) error {
	var trophiesJson []Trophies

	// Open the file
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &trophiesJson); err != nil {
		return err
	}

	trophiesExpectedTotalRecords := len(trophiesJson)

	trophiesCI := NewConcurrentInserter(trophiesExpectedTotalRecords, func(trophies []models.Trophy) error {
		trophiesRepo := trophiesrepo.NewRepo(db)

		if _, err := trophiesRepo.AddTrophies(ctx, trophies...); err != nil {
			return err
		}

		return nil
	}, 300*time.Millisecond, 500)

	wg := new(sync.WaitGroup)
	wg.Add(trophiesExpectedTotalRecords)

	go func() {
		for _, trophy := range trophiesJson {
			newTrophy := models.Trophy{
				ID:          uuid.New(),
				Title:       trophy.Title,
				Description: trophy.Description,
				ImageLink:   trophy.ImageLink,
			}
			trophiesCI.ResC <- newTrophy
			wg.Done()
		}
	}()

	if err := trophiesCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func insertAwards(ctx context.Context, db *bun.DB, filename string) error {
	var awardsJson []Awards

	// Open the file
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &awardsJson); err != nil {
		return err
	}

	awardsExpectedTotalRecords := len(awardsJson)

	awardsCI := NewConcurrentInserter(awardsExpectedTotalRecords, func(awards []models.Award) error {
		awardsRepo := awardsrepo.NewRepo(db)

		if _, err := awardsRepo.AddAwards(ctx, awards...); err != nil {
			return err
		}

		return nil
	}, 300*time.Millisecond, 500)

	wg := new(sync.WaitGroup)
	wg.Add(awardsExpectedTotalRecords)

	go func() {
		for _, award := range awardsJson {
			newAward := models.Award{
				ID:        uuid.New(),
				Title:     award.Title,
				ImageLink: award.ImageLink,
			}
			awardsCI.ResC <- newAward
			wg.Done()
		}
	}()

	if err := awardsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func getTrophyByTitle(_ context.Context, trophies []models.Trophy, trophyTitle string) (models.Trophy, error) {
	for _, trophy := range trophies {
		if trophy.Title == trophyTitle {
			return trophy, nil
		}
	}

	return models.Trophy{}, fmt.Errorf("failed to get trophy with title %v", trophyTitle)
}

func insertUsers(ctx context.Context, db *bun.DB, usersFilename string) error {
	var usersJson []User

	// Open the file
	jsonBytes, err := os.ReadFile(usersFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &usersJson); err != nil {
		return err
	}

	usersExpectedTotalRecords := len(usersJson)

	usersCI := NewConcurrentInserter(usersExpectedTotalRecords, func(users []models.User) error {
		usersRepo := usersrepo.NewRepo(db)

		if _, err := usersRepo.AddUsers(ctx, users...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 100)

	wg := new(sync.WaitGroup)
	wg.Add(usersExpectedTotalRecords)

	userTrophiesMap := threadsafe.NewThreadSafeMap[uuid.UUID, []models.Trophy]()

	trophyRepo := trophiesrepo.NewRepo(db)
	trophies, err := trophyRepo.Trophies(ctx)
	if err != nil {
		return err
	}

	var errs error

	go func() {
		for _, user := range usersJson {
			userID := uuid.New()

			newUser := models.User{
				ID:                userID,
				Name:              user.Name,
				PublicDescription: &user.PublicDescription,
				AvatarImg:         &user.AvatarImg,
				BannerImg:         &user.BannerImg,
				Iconcolor:         &user.IconColor,
				Keycolor:          &user.KeyColor,
				Primarycolor:      &user.PrimaryColor,
				Over18:            user.Over18,
				Suspended:         user.Suspended,
				CreatedAt:         time.Unix(int64(user.CakeDayUTC), 0),
				CreatedAtUnix:     int64(user.CakeDayUTC),
				UpdatedAt:         time.Unix(int64(user.CakeDayUTC), 0),
			}

			// insert user trophies
			userTrophies := make([]models.Trophy, 0)
			for _, trophy := range user.Trophies {
				t, err := getTrophyByTitle(ctx, trophies, trophy.Title)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				userTrophies = append(userTrophies, t)
			}

			userTrophiesMap.Put(userID, userTrophies)

			usersCI.ResC <- newUser
			wg.Done()
		}
	}()

	if err := usersCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	if errs != nil {
		return errs
	}

	userTrophiesExpectedTotalRecords := 0
	for _, trophies := range userTrophiesMap.Load() {
		userTrophiesExpectedTotalRecords += len(trophies)
	}

	userTrophiesCI := NewConcurrentInserter(userTrophiesExpectedTotalRecords, func(userTrophies []models.UserTrophy) error {
		userTrophiesRepo := relationrepo.NewRepo(db)

		if _, err := userTrophiesRepo.LinkUserTrophies(ctx, userTrophies...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 3000)

	wg.Add(userTrophiesExpectedTotalRecords)

	go func() {
		for userID, trophies := range userTrophiesMap.Load() {
			for _, trophy := range trophies {
				newUserTrophy := models.UserTrophy{
					UserID:   userID,
					TrophyID: trophy.ID,
				}
				userTrophiesCI.ResC <- newUserTrophy
				wg.Done()
			}
		}
	}()

	if err := userTrophiesCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func getTopicByName(_ context.Context, topics []models.Topic, topicName string) (models.Topic, error) {
	for _, topic := range topics {
		if topic.Name == topicName {
			return topic, nil
		}
	}

	return models.Topic{}, fmt.Errorf("failed to get topic with title %v", topicName)
}

func getuserByName(_ context.Context, users []models.User, userName string) (models.User, error) {
	for _, user := range users {
		if user.Name == userName {
			return user, nil
		}
	}

	return models.User{}, fmt.Errorf("failed to get user with name %v", userName)
}

func getpostByTitleAndAuthorIDAndVoxsphereID(
	_ context.Context,
	posts []models.Post,
	postTitle string,
	authorID uuid.UUID,
	voxsphereID uuid.UUID,
) (models.Post, error) {
	for _, post := range posts {
		if post.Title == postTitle && post.AuthorID == authorID && post.VoxsphereID == voxsphereID {
			return post, nil
		}
	}

	return models.Post{}, fmt.Errorf("failed to get post with title %v", postTitle)
}

func insertVoxspheres(ctx context.Context, db *bun.DB, subredditsFilename string) error {
	var subredditJson map[string][]Subreddit

	// Open the file
	jsonBytes, err := os.ReadFile(subredditsFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &subredditJson); err != nil {
		return err
	}

	voxspheresExpectedTotalRecords := 0

	for _, sreddits := range subredditJson {
		voxspheresExpectedTotalRecords += len(sreddits)
	}

	voxspheresCI := NewConcurrentInserter(voxspheresExpectedTotalRecords, func(voxspheres []models.Voxsphere) error {
		voxspheresRepo := voxspheresrepo.NewRepo(db)

		if _, err := voxspheresRepo.AddVoxspheres(ctx, voxspheres...); err != nil {
			return err
		}

		return nil
	}, 1000*time.Millisecond, 500)

	wg := new(sync.WaitGroup)
	wg.Add(voxspheresExpectedTotalRecords)

	topicRepo := topicsrepo.NewRepo(db)
	topics, err := topicRepo.Topics(ctx)
	if err != nil {
		return err
	}

	var errs error

	voxsphereRulesMap := threadsafe.NewThreadSafeMap[uuid.UUID, []SubredditRule]()
	voxsphereMembersMap := threadsafe.NewThreadSafeMap[uuid.UUID, []User]()
	voxsphereModeratorsMap := threadsafe.NewThreadSafeMap[uuid.UUID, []User]()

	go func() {
		for _, subreddits := range subredditJson {
			for _, subreddit := range subreddits {

				voxsphereID := uuid.New()

				topic, err := getTopicByName(ctx, topics, subreddit.Topic)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}

				newVoxsphere := models.Voxsphere{
					ID:                    voxsphereID,
					TopicID:               topic.ID,
					Title:                 subreddit.Title,
					PublicDescription:     &subreddit.PublicDescription,
					CommunityIcon:         &subreddit.CommunityIcon,
					BannerBackgroundImage: &subreddit.BannerBackgroundImage,
					BannerBackgroundColor: &subreddit.BannerBackgroundColor,
					KeyColor:              &subreddit.KeyColor,
					PrimaryColor:          &subreddit.PrimaryColor,
					Over18:                subreddit.Over18,
					SpoilersEnabled:       subreddit.SpoilersEnabled,
					CreatedAt:             time.Unix(int64(subreddit.CreatedUTC), 0),
					CreatedAtUnix:         int64(subreddit.CreatedUTC),
					UpdatedAt:             time.Unix(int64(subreddit.CreatedUTC), 0),
				}

				voxsphereRulesMap.Put(voxsphereID, subreddit.Rules)
				voxsphereMembersMap.Put(voxsphereID, subreddit.Members)
				voxsphereModeratorsMap.Put(voxsphereID, subreddit.Moderators)

				voxspheresCI.ResC <- newVoxsphere
				wg.Done()
			}
		}
	}()

	if err := voxspheresCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	if errs != nil {
		return errs
	}

	// insert rules
	voxsphereRulesExpectedTotalRecords := 0
	for _, voxsphereRules := range voxsphereRulesMap.Load() {
		voxsphereRulesExpectedTotalRecords += len(voxsphereRules)
	}

	wg.Add(voxsphereRulesExpectedTotalRecords)

	rulesCI := NewConcurrentInserter(voxsphereRulesExpectedTotalRecords, func(rules []models.Rule) error {
		rulesRepo := rulesrepo.NewRepo(db)

		if _, err := rulesRepo.AddRules(ctx, rules...); err != nil {
			return err
		}

		return nil
	}, 1000*time.Millisecond, 500)

	go func() {
		for voxsphereID, rules := range voxsphereRulesMap.Load() {
			for _, rule := range rules {
				newRule := models.Rule{
					ID:          uuid.New(),
					VoxsphereID: voxsphereID,
					ShortName:   rule.ShortName,
					Description: rule.Description,
				}
				rulesCI.ResC <- newRule
				wg.Done()
			}
		}
	}()

	if err := rulesCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userRepo := usersrepo.NewRepo(db)
	users, err := userRepo.Users(ctx)
	if err != nil {
		return err
	}

	// insert members
	voxsphereMembersExpectedTotalRecords := 0
	for _, voxsphereMembers := range voxsphereMembersMap.Load() {
		voxsphereMembersExpectedTotalRecords += len(voxsphereMembers)
	}

	wg.Add(voxsphereMembersExpectedTotalRecords)
	errs = nil

	memebersCI := NewConcurrentInserter(voxsphereMembersExpectedTotalRecords, func(voxsphereMembers []models.VoxsphereMember) error {
		relationRepo := relationrepo.NewRepo(db)

		if _, err := relationRepo.LinkVoxsphereMembers(ctx, voxsphereMembers...); err != nil {
			return err
		}

		return nil
	}, 300*time.Millisecond, 100_000)

	go func() {
		for voxsphereID, members := range voxsphereMembersMap.Load() {
			for _, member := range members {
				user, err := getuserByName(ctx, users, member.Name)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				newVoxsphereMember := models.VoxsphereMember{
					VoxsphereID: voxsphereID,
					UserID:      user.ID,
				}

				memebersCI.ResC <- newVoxsphereMember
				wg.Done()
			}
		}
	}()

	if err := memebersCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	if errs != nil {
		return errs
	}

	// insert moderators
	voxsphereModeratorsExpectedTotalRecords := 0
	for _, voxsphereModerators := range voxsphereModeratorsMap.Load() {
		voxsphereModeratorsExpectedTotalRecords += len(voxsphereModerators)
	}

	wg.Add(voxsphereModeratorsExpectedTotalRecords)
	errs = nil

	moderatorsCI := NewConcurrentInserter(voxsphereModeratorsExpectedTotalRecords, func(voxsphereModerators []models.VoxsphereModerator) error {
		relationRepo := relationrepo.NewRepo(db)

		if _, err := relationRepo.LinkVoxsphereModerators(ctx, voxsphereModerators...); err != nil {
			return err
		}

		return nil
	}, 300*time.Millisecond, 100_000)

	go func() {
		for voxsphereID, moderators := range voxsphereModeratorsMap.Load() {
			for _, moderator := range moderators {
				user, err := getuserByName(ctx, users, moderator.Name)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				newVoxsphereModerator := models.VoxsphereModerator{
					VoxsphereID: voxsphereID,
					UserID:      user.ID,
				}

				moderatorsCI.ResC <- newVoxsphereModerator
				wg.Done()
			}
		}
	}()

	if err := moderatorsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func getvoxsphereByTitle(_ context.Context, voxspheres []models.Voxsphere, voxsphereTitle string) (models.Voxsphere, error) {
	for _, voxsphere := range voxspheres {
		if voxsphere.Title == voxsphereTitle {
			return voxsphere, nil
		}
	}

	return models.Voxsphere{}, fmt.Errorf("failed to get voxsphere with title %v", voxsphereTitle)
}

func getawardByTitle(_ context.Context, awards []models.Award, awardTitle string) (models.Award, error) {
	for _, award := range awards {
		if award.Title == awardTitle {
			return award, nil
		}
	}

	return models.Award{}, fmt.Errorf("failed to get award with title %v", awardTitle)
}

// mapToStruct is a helper function to unmarshal the content into the appropriate struct
func mapToStruct(content interface{}, out interface{}) error {
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return json.Unmarshal(contentBytes, out)
}

func handleResolutions(resolutions *[]interface{}) error {
	for i, res := range *resolutions {
		switch r := res.(type) {
		case map[string]interface{}:
			// Check if it has the "height" and "width" fields for ImageResolution
			if height, ok := r["height"].(float64); ok {
				if width, ok := r["width"].(float64); ok {
					imageRes := ImageResolution{
						Height: int(height),
						Width:  int(width),
						URL:    r["url"].(string),
					}
					// Update the original resolutions slice
					(*resolutions)[i] = imageRes
					fmt.Printf("Updated ImageResolution: %+v\n", imageRes)
					continue
				}
			}
			// Check if it has the "x", "y", and "u" fields for ImageMultiResolution
			if x, ok := r["x"].(float64); ok {
				if y, ok := r["y"].(float64); ok {
					imageMultiRes := ImageMultiResolution{
						X: int(x),
						Y: int(y),
						U: r["u"].(string),
					}
					// Update the original resolutions slice
					(*resolutions)[i] = imageMultiRes
					fmt.Printf("Updated ImageMultiResolution: %+v\n", imageMultiRes)
					continue
				}
			}
		default:
			return fmt.Errorf("Unknown resolution type")
		}
	}
	return nil
}

func insertPostMedias(
	ctx context.Context,
	db *bun.DB,
	postMediasMap *threadsafe.ThreadSafeMap[uuid.UUID, MediaContent],
) error {
	var postMedias []models.PostMedia

	var images []models.Image
	var imageMetadatas []models.ImageMetadata

	var gifs []models.Gif
	var gifMetadatas []models.GifMetadata

	var galleries []models.Gallery
	var galleryMetadatas []models.GalleryMetadata

	var videos []models.Video

	var links []models.Link

	mediaRepo := mediasrepo.NewRepo(db)

	for postID, postMediaContent := range postMediasMap.Load() {
		postMediaID := uuid.New()
		postMedia := models.PostMedia{
			ID:     postMediaID,
			PostID: postID,
		}

		switch postMediaContent.Type {
		case "image":
			var image Image
			err := mapToStruct(postMediaContent.Content, &image)
			if err != nil {
				return fmt.Errorf("Error casting to image for post id %v %v", postID, err)
			}
			if err := handleResolutions(&image.Resolutions); err != nil {
				return err
			}

			imageID := uuid.New()
			newImage := models.Image{
				ID:      imageID,
				MediaID: postMediaID,
			}
			images = append(images, newImage)

			// Determine the type of each resolution
			for _, res := range image.Resolutions {
				switch r := res.(type) {
				case ImageResolution:
					newImageMetadata := models.ImageMetadata{
						ID:            uuid.New(),
						ImageID:       imageID,
						Height:        int32(r.Height),
						Width:         int32(r.Width),
						Url:           r.URL,
						CreatedAt:     time.Now(),
						CreatedAtUnix: time.Now().Unix(),
						UpdatedAt:     time.Now(),
					}
					imageMetadatas = append(imageMetadatas, newImageMetadata)
				default:
					return fmt.Errorf("Unknown image resolution type for post id %v", postID)
				}
			}
			postMedia.MediaType = models.MediaTypeImage
			postMedias = append(postMedias, postMedia)
		case "gif":
			var gif Gif
			err := mapToStruct(postMediaContent.Content, &gif)
			if err != nil {
				return fmt.Errorf("Error casting to gif for post id %v %v", postID, err)
			}
			if err := handleResolutions(&gif.Resolutions); err != nil {
				return err
			}

			gifID := uuid.New()
			newGif := models.Gif{
				ID:      gifID,
				MediaID: postMediaID,
			}
			gifs = append(gifs, newGif)

			// Determine the type of each resolution
			for _, res := range gif.Resolutions {
				switch r := res.(type) {
				case ImageResolution:
					newGifMetadata := models.GifMetadata{
						ID:            uuid.New(),
						GifID:         gifID,
						Height:        int32(r.Height),
						Width:         int32(r.Width),
						Url:           r.URL,
						CreatedAt:     time.Now(),
						CreatedAtUnix: time.Now().Unix(),
						UpdatedAt:     time.Now(),
					}
					gifMetadatas = append(gifMetadatas, newGifMetadata)
				default:
					return fmt.Errorf("Unknown gif resolution type for post id %v", postID)
				}
			}
			postMedia.MediaType = models.MediaTypeGif
			postMedias = append(postMedias, postMedia)

		case "video":
			var video Video
			err := mapToStruct(postMediaContent.Content, &video)
			if err != nil {
				return fmt.Errorf("Error casting to video for post id %v %v", postID, err)
			}

			videoID := uuid.New()
			newVideo := models.Video{
				ID:            videoID,
				MediaID:       postMediaID,
				Url:           video.FallbackURL,
				Height:        int32(video.Height),
				Width:         int32(video.Width),
				CreatedAt:     time.Now(),
				CreatedAtUnix: time.Now().Unix(),
				UpdatedAt:     time.Now(),
			}
			videos = append(videos, newVideo)
			postMedia.MediaType = models.MediaTypeVideo
			postMedias = append(postMedias, postMedia)
		case "link":
			var link Link
			err := mapToStruct(postMediaContent.Content, &link)
			if err != nil {
				return fmt.Errorf("Error casting to link for post id %v %v", postID, err)
			}
			linkID := uuid.New()
			newLink := models.Link{
				ID:            linkID,
				MediaID:       postMediaID,
				Link:          link.Link,
				CreatedAt:     time.Now(),
				CreatedAtUnix: time.Now().Unix(),
				UpdatedAt:     time.Now(),
			}
			links = append(links, newLink)
			postMedia.MediaType = models.MediaTypeLink
			postMedias = append(postMedias, postMedia)
		case "gallery":
			var gallery Gallery
			err := mapToStruct(postMediaContent.Content, &gallery)
			if err != nil {
				return fmt.Errorf("Error casting to gallery for post id %v %v", postID, err)
			}
			galleryID := uuid.New()
			newGallery := models.Gallery{
				ID:      galleryID,
				MediaID: postMediaID,
			}
			galleries = append(galleries, newGallery)
			index := 0
			for _, galleryImages := range gallery.Images {
				for _, galleryImage := range galleryImages {
					newGalleryMetadataID := uuid.New()
					newGalleryMetadata := models.GalleryMetadata{
						ID:            newGalleryMetadataID,
						GalleryID:     galleryID,
						OrderIndex:    int32(index),
						Height:        int32(galleryImage.Y),
						Width:         int32(galleryImage.X),
						Url:           galleryImage.U,
						CreatedAt:     time.Now(),
						CreatedAtUnix: time.Now().Unix(),
						UpdatedAt:     time.Now(),
					}
					galleryMetadatas = append(galleryMetadatas, newGalleryMetadata)
				}
				index += 1
			}

			postMedia.MediaType = models.MediaTypeGallery
			postMedias = append(postMedias, postMedia)
		case "multi":
			var multipleContent []any
			err := mapToStruct(postMediaContent.Content, &multipleContent)
			if err != nil {
				return fmt.Errorf("Error casting to multiple media content for post id %v %v", postID, err)
			}
			for _, content := range multipleContent {
				switch r := content.(type) {
				case map[string]interface{}:
					mediaType, found := r["_type"]
					if !found {
						return fmt.Errorf("url retrieval failed from multiple media content")
					}
					if mediaType == "image" {
						resolutions, found := r["resolutions"]
						if !found {
							return fmt.Errorf("resolutions retrieval failed from image multiple media content")
						}

						ImageID := uuid.New()
						newImage := models.Image{
							ID:      ImageID,
							MediaID: postMediaID,
						}
						images = append(images, newImage)

						var imageResolutions []ImageMultiResolution
						if err := mapToStruct(resolutions, &imageResolutions); err != nil {
							return err
						}

						for _, imageResolution := range imageResolutions {
							imageMetadatas = append(imageMetadatas, models.ImageMetadata{
								ID:            uuid.New(),
								ImageID:       ImageID,
								Height:        int32(imageResolution.Y),
								Width:         int32(imageResolution.X),
								Url:           imageResolution.U,
								CreatedAt:     time.Now(),
								CreatedAtUnix: time.Now().Unix(),
								UpdatedAt:     time.Now(),
							})
						}
					} else if mediaType == "gif" {
						resolutions, found := r["resolutions"]
						if !found {
							return fmt.Errorf("resolutions retrieval failed from gif multiple media content")
						}

						gifID := uuid.New()
						newGif := models.Gif{
							ID:      gifID,
							MediaID: postMediaID,
						}
						gifs = append(gifs, newGif)

						var gifResolutions []ImageMultiResolution
						if err := mapToStruct(resolutions, &gifResolutions); err != nil {
							return err
						}

						for _, imageResolution := range gifResolutions {
							gifMetadatas = append(gifMetadatas, models.GifMetadata{
								ID:            uuid.New(),
								GifID:         gifID,
								Height:        int32(imageResolution.Y),
								Width:         int32(imageResolution.X),
								Url:           imageResolution.U,
								CreatedAt:     time.Now(),
								CreatedAtUnix: time.Now().Unix(),
								UpdatedAt:     time.Now(),
							})
						}

					} else if mediaType == "video" {
						var video VideoMulti
						if err := mapToStruct(content, &video); err != nil {
							return err
						}
						newVideo := models.Video{
							ID:            uuid.New(),
							MediaID:       postMediaID,
							Url:           video.HLSURL,
							Height:        int32(video.Y),
							Width:         int32(video.X),
							CreatedAt:     time.Now(),
							CreatedAtUnix: time.Now().Unix(),
							UpdatedAt:     time.Now(),
						}
						videos = append(videos, newVideo)
					} else {
						return fmt.Errorf("Unknown mediaType for multiple media content")
					}
				default:
					return fmt.Errorf("Unknown multiple content type")
				}
			}
			postMedia.MediaType = models.MediaTypeMulti
			postMedias = append(postMedias, postMedia)

		default:
			fmt.Println("Unknown media type:", postMediaContent.Type)
		}
	}

	wg := new(sync.WaitGroup)

	postMediasCI := NewConcurrentInserter(len(postMedias), func(postMedias []models.PostMedia) error {
		if _, err := mediaRepo.AddPostMedias(ctx, postMedias...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(postMedias))

	go func() {
		for _, postMedia := range postMedias {
			postMediasCI.ResC <- postMedia
			wg.Done()
		}
	}()

	if err := postMediasCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	imagesCI := NewConcurrentInserter(len(images), func(images []models.Image) error {
		if _, err := mediaRepo.AddImages(ctx, images...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(images))

	go func() {
		for _, image := range images {
			imagesCI.ResC <- image
			wg.Done()
		}
	}()

	if err := imagesCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	imageMetadatasCI := NewConcurrentInserter(len(imageMetadatas), func(imageMetadatas []models.ImageMetadata) error {
		if _, err := mediaRepo.AddImageMetadatas(ctx, imageMetadatas...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(imageMetadatas))

	go func() {
		for _, imageMetadata := range imageMetadatas {
			imageMetadatasCI.ResC <- imageMetadata
			wg.Done()
		}
	}()

	if err := imageMetadatasCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	gifsCI := NewConcurrentInserter(len(gifs), func(gifs []models.Gif) error {
		if _, err := mediaRepo.AddGifs(ctx, gifs...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(gifs))

	go func() {
		for _, gif := range gifs {
			gifsCI.ResC <- gif
			wg.Done()
		}
	}()

	if err := gifsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	gifMetadatasCI := NewConcurrentInserter(len(gifMetadatas), func(gifMetadatas []models.GifMetadata) error {
		if _, err := mediaRepo.AddGifMetadatas(ctx, gifMetadatas...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(gifMetadatas))

	go func() {
		for _, gifMetadata := range gifMetadatas {
			gifMetadatasCI.ResC <- gifMetadata
			wg.Done()
		}
	}()

	if err := gifMetadatasCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	galleriesCI := NewConcurrentInserter(len(galleries), func(galleries []models.Gallery) error {
		if _, err := mediaRepo.AddGalleries(ctx, galleries...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(galleries))

	go func() {
		for _, gallery := range galleries {
			galleriesCI.ResC <- gallery
			wg.Done()
		}
	}()

	if err := galleriesCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	galleriesMetadatasCI := NewConcurrentInserter(len(galleryMetadatas), func(galleryMetadatas []models.GalleryMetadata) error {
		if _, err := mediaRepo.AddGalleryMetadatas(ctx, galleryMetadatas...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(galleryMetadatas))

	go func() {
		for _, galleryMetadata := range galleryMetadatas {
			galleriesMetadatasCI.ResC <- galleryMetadata
			wg.Done()
		}
	}()

	if err := galleriesMetadatasCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	videosCI := NewConcurrentInserter(len(videos), func(videos []models.Video) error {
		if _, err := mediaRepo.AddVideos(ctx, videos...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(videos))

	go func() {
		for _, video := range videos {
			videosCI.ResC <- video
			wg.Done()
		}
	}()

	if err := videosCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	linksCI := NewConcurrentInserter(len(links), func(links []models.Link) error {
		if _, err := mediaRepo.AddLinks(ctx, links...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(links))

	go func() {
		for _, link := range links {
			linksCI.ResC <- link
			wg.Done()
		}
	}()

	if err := linksCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func insertPostComments(
	ctx context.Context,
	db *bun.DB,
	postCommentsMap *threadsafe.ThreadSafeMap[uuid.UUID, []Comment],
	parentCommentID uuid.UUID,
	resC chan models.Comment,
	users []models.User,
) error {
	for postID, comments := range postCommentsMap.Load() {
		for _, comment := range comments {
			newRepliesMap := threadsafe.NewThreadSafeMap[uuid.UUID, []Comment]()
			newRepliesMap.Put(postID, comment.Replies)
			commentID := uuid.New()

			user, err := getuserByName(ctx, users, comment.Author)
			if err != nil {
				return err
			}

			newComment := models.Comment{
				ID:              commentID,
				AuthorID:        user.ID,
				ParentCommentID: parentCommentID,
				PostID:          postID,
				Body:            comment.Body,
				BodyHtml:        comment.BodyHTML,
				Ups:             int32(comment.Ups),
				Score:           int32(comment.Score),
				CreatedAt:       time.Unix(int64(comment.CreatedUTC), 0),
				CreatedAtUnix:   int64(comment.CreatedUTC),
				UpdatedAt:       time.Unix(int64(comment.CreatedUTC), 0),
			}

			resC <- newComment

			if err := insertPostComments(ctx, db, newRepliesMap, commentID, resC, users); err != nil {
				return err
			}

			newRepliesMap.Clear()
		}
	}

	return nil
}

func insertPosts(ctx context.Context, db *bun.DB, postsFilename string) error {
	var postsJson []Post

	// Open the file
	jsonBytes, err := os.ReadFile(postsFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &postsJson); err != nil {
		return err
	}

	postsExpectedTotalRecords := len(postsJson)

	wg := new(sync.WaitGroup)
	wg.Add(postsExpectedTotalRecords)

	userRepo := usersrepo.NewRepo(db)
	users, err := userRepo.Users(ctx)
	if err != nil {
		return err
	}

	voxsphereRepo := voxspheresrepo.NewRepo(db)
	voxspheres, err := voxsphereRepo.Voxspheres(ctx)
	if err != nil {
		return err
	}

	postsCI := NewConcurrentInserter(postsExpectedTotalRecords, func(posts []models.Post) error {
		postRepo := postsrepo.NewRepo(db)

		if _, err := postRepo.AddPosts(ctx, posts...); err != nil {
			return err
		}

		return nil
	}, 1000*time.Millisecond, 500)

	var errs error

	postAwardsMap := threadsafe.NewThreadSafeMap[uuid.UUID, []Awards]()
	postMediasMap := threadsafe.NewThreadSafeMap[uuid.UUID, MediaContent]()
	postCommentsMap := threadsafe.NewThreadSafeMap[uuid.UUID, []Comment]()

	go func() {
		// posts
		for _, post := range postsJson {
			postID := uuid.New()

			user, err := getuserByName(ctx, users, post.Author)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}

			voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, post.Subreddit)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}

			newPost := models.Post{
				ID:            postID,
				AuthorID:      user.ID,
				VoxsphereID:   voxsphere.ID,
				Title:         post.Title,
				Text:          post.Text,
				TextHtml:      post.TextHTML,
				Ups:           int32(post.Ups),
				Over18:        post.Over18,
				Spoiler:       post.Spoiler,
				CreatedAt:     time.Unix(int64(post.CreatedUTC), 0),
				CreatedAtUnix: int64(post.CreatedUTC),
				UpdatedAt:     time.Unix(int64(post.CreatedUTC), 0),
			}

			postAwardsMap.Put(postID, post.Awards)
			postMediasMap.Put(postID, post.MediaContent)
			postCommentsMap.Put(postID, post.Comments)

			postsCI.ResC <- newPost
			wg.Done()
		}
	}()

	if err := postsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	if errs != nil {
		return errs
	}

	// insert post medias
	if err := insertPostMedias(ctx, db, postMediasMap); err != nil {
		return err
	}

	commentsExpectedTotalRecords := 0
	for _, comments := range postCommentsMap.Load() {
		commentsExpectedTotalRecords += len(comments)
	}

	commentsCI := NewConcurrentInserter(commentsExpectedTotalRecords, func(comments []models.Comment) error {
		commentRepo := commentsrepo.NewRepo(db)

		if _, err := commentRepo.AddComments(ctx, comments...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 3000)

	// insert post comments
	ctxnew, cancel := context.WithCancel(ctx)
	go func() {
		if err := insertPostComments(ctx, db, postCommentsMap, uuid.Nil, commentsCI.ResC, users); err != nil {
			cancel()
		}
	}()

	if err := commentsCI.Serve(ctxnew); err != nil {
		return err
	}

	// insert post awards
	postAwardsExpectedTotalRecords := 0
	for _, postAwards := range postAwardsMap.Load() {
		postAwardsExpectedTotalRecords += len(postAwards)
	}

	postAwardsCI := NewConcurrentInserter(postAwardsExpectedTotalRecords, func(postAward []models.PostAward) error {
		relationRepo := relationrepo.NewRepo(db)

		if _, err := relationRepo.LinkPostAwards(ctx, postAward...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 3000)

	wg.Add(postAwardsExpectedTotalRecords)
	errs = nil

	awardRepo := awardsrepo.NewRepo(db)
	awards, err := awardRepo.Awards(ctx)
	if err != nil {
		return err
	}

	go func() {
		for postID, postAwards := range postAwardsMap.Load() {
			for _, postAward := range postAwards {

				award, err := getawardByTitle(ctx, awards, postAward.Title)
				if err != nil {
					errs = errors.Join(errs, err)
				}

				newPostAward := models.PostAward{
					PostID:  postID,
					AwardID: award.ID,
				}

				postAwardsCI.ResC <- newPostAward
				wg.Done()
			}
		}
	}()

	if err := postAwardsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	if errs != nil {
		return errs
	}

	return nil
}

func insertEmojis(ctx context.Context, db *bun.DB, subredditsFilename string) error {
	var subredditJson map[string][]Subreddit

	// Open the file
	jsonBytes, err := os.ReadFile(subredditsFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytes), &subredditJson); err != nil {
		return err
	}

	goEmojis := gomoji.AllEmojis()
	codePoints := make(map[string]bool)

	for _, goEmoji := range goEmojis {
		result := strings.ToLower(goEmoji.CodePoint)
		joined := strings.Join(strings.Fields(result), "-")
		codePoints[joined] = true
	}

	emojiRepo := emojirepo.NewRepo(db)

	emojisExpectedTotalRecords := len(codePoints)

	emojiCI := NewConcurrentInserter(emojisExpectedTotalRecords, func(emojis []models.Emoji) error {
		if _, err := emojiRepo.AddEmojis(ctx, emojis...); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond, 500)

	wg := new(sync.WaitGroup)
	wg.Add(emojisExpectedTotalRecords)

	go func() {
		for codePoint := range codePoints {
			newEmoji := models.Emoji{
				ID:    uuid.New(),
				Title: codePoint,
			}
			emojiCI.ResC <- newEmoji
			wg.Done()
		}
	}()

	if err := emojiCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	voxspheresRepo := voxspheresrepo.NewRepo(db)
	voxspheres, err := voxspheresRepo.Voxspheres(ctx)
	if err != nil {
		return err
	}

	customEmojis := make(map[string]models.CustomEmoji)

	for _, subreddits := range subredditJson {
		for _, subreddit := range subreddits {
			voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, subreddit.Title)
			if err != nil {
				return err
			}
			for _, postFlair := range subreddit.Flairs {
				richTexts := postFlair.RichText
				for _, richtText := range richTexts {
					switch r := richtText.(type) {
					case map[string]any:
						_, found := r["url"]
						if found {
							var flairEmoji FlairEmoji
							if err := mapToStruct(r, &flairEmoji); err != nil {
								return err
							}
							if helper.IsCustomEmoji(flairEmoji.Text) {
								key := fmt.Sprintf("%v-%v", voxsphere.ID, flairEmoji.Text)
								customEmojis[key] = models.CustomEmoji{
									ID:          uuid.New(),
									VoxsphereID: voxsphere.ID,
									Url:         flairEmoji.URL,
									Title:       flairEmoji.Text,
								}
							}
						}
					default:
						return fmt.Errorf("Unknown type for rich text")
					}
				}
			}
			for _, userFlair := range subreddit.UserFlairs {
				richTexts := userFlair.RichText
				for _, richtText := range richTexts {
					switch r := richtText.(type) {
					case map[string]any:
						_, found := r["url"]
						if found {
							var flairEmoji FlairEmoji
							if err := mapToStruct(r, &flairEmoji); err != nil {
								return err
							}
							if helper.IsCustomEmoji(flairEmoji.Text) {
								key := fmt.Sprintf("%v-%v", voxsphere.ID, flairEmoji.Text)
								customEmojis[key] = models.CustomEmoji{
									ID:          uuid.New(),
									VoxsphereID: voxsphere.ID,
									Url:         flairEmoji.URL,
									Title:       flairEmoji.Text,
								}
							}
						}
					default:
						return fmt.Errorf("Unknown type for rich text")
					}
				}
			}
		}
	}

	customemojiRepo := customemojirepo.NewRepo(db)

	customemojisExpectedTotalRecords := len(customEmojis)

	customemojiCI := NewConcurrentInserter(customemojisExpectedTotalRecords, func(customemojis []models.CustomEmoji) error {
		if _, err := customemojiRepo.AddCustomEmojis(ctx, customemojis...); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(customemojisExpectedTotalRecords)

	go func() {
		for _, customEmoji := range customEmojis {
			customemojiCI.ResC <- customEmoji
			wg.Done()
		}
	}()

	if err := customemojiCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func getCustomEmojiByTitleandVoxsphereID(
	_ context.Context,
	customEmojis []models.CustomEmoji,
	customEmojiTitle string, voxsphereID uuid.UUID,
) (models.CustomEmoji, error) {
	for _, customEmoji := range customEmojis {
		if customEmoji.Title == customEmojiTitle && customEmoji.VoxsphereID == voxsphereID {
			return customEmoji, nil
		}
	}

	return models.CustomEmoji{}, fmt.Errorf(
		"failed to get customEmoji with title %v and voxsphere ID %v",
		customEmojiTitle,
		voxsphereID,
	)
}

func getEmojiByTitle(
	_ context.Context,
	emojis []models.Emoji,
	emojiTitle string,
) (models.Emoji, error) {
	for _, emoji := range emojis {
		if emoji.Title == emojiTitle {
			return emoji, nil
		}
	}

	return models.Emoji{}, fmt.Errorf("failed to get customEmoji with title %v", emojiTitle)
}

func insertFlairs(ctx context.Context, db *bun.DB, postsFilename string, subredditsFilename string) error {
	var subredditJson map[string][]Subreddit

	// Open the file
	jsonBytesSubreddits, err := os.ReadFile(subredditsFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytesSubreddits), &subredditJson); err != nil {
		return err
	}

	var postsJson []Post

	// Open the file
	jsonBytesPosts, err := os.ReadFile(postsFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal([]byte(jsonBytesPosts), &postsJson); err != nil {
		return err
	}

	voxspheresRepo := voxspheresrepo.NewRepo(db)
	voxspheres, err := voxspheresRepo.Voxspheres(ctx)
	if err != nil {
		return err
	}

	customEmojiRepo := customemojirepo.NewRepo(db)
	customEmojis, err := customEmojiRepo.CustomEmojis(ctx)
	if err != nil {
		return err
	}

	emojiRepo := emojirepo.NewRepo(db)
	emojis, err := emojiRepo.Emojis(ctx)
	if err != nil {
		return err
	}

	postFlairs := make([]models.PostFlair, 0)
	postFlairEmojis := make([]models.PostFlairEmoji, 0)
	postFlairCustomEmojis := make([]models.PostFlairCustomEmoji, 0)
	postFlairDescriptions := make([]models.PostFlairDescription, 0)

	userFlairs := make([]models.UserFlair, 0)
	userFlairEmojis := make([]models.UserFlairEmoji, 0)
	userFlairCustomEmojis := make([]models.UserFlairCustomEmoji, 0)
	userFlairDescriptions := make([]models.UserFlairDescription, 0)

	for _, subreddits := range subredditJson {
		for _, subreddit := range subreddits {
			voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, subreddit.Title)
			if err != nil {
				return err
			}

			// post flairs
			for _, postFlair := range subreddit.Flairs {
				postFlairID := uuid.New()
				postFlairs = append(postFlairs, models.PostFlair{
					ID:              postFlairID,
					VoxsphereID:     voxsphere.ID,
					FullText:        postFlair.FullText,
					BackgroundColor: postFlair.BackgroundColor,
				})

				orderIndex := 0

				richTexts := postFlair.RichText
				for _, richtText := range richTexts {
					switch r := richtText.(type) {
					case map[string]any:
						_, found := r["url"]
						if found {
							var flairEmoji FlairEmoji
							if err := mapToStruct(r, &flairEmoji); err != nil {
								return err
							}
							if helper.IsCustomEmoji(flairEmoji.Text) {
								customEmoji, err := getCustomEmojiByTitleandVoxsphereID(
									ctx,
									customEmojis,
									flairEmoji.Text,
									voxsphere.ID,
								)
								if err != nil {
									return err
								}
								postFlairCustomEmojis = append(postFlairCustomEmojis, models.PostFlairCustomEmoji{
									CustomEmojiID: customEmoji.ID,
									PostFlairID:   postFlairID,
									OrderIndex:    int32(orderIndex),
								})
								orderIndex += 1
							}
						} else {
							var flairText FlairText
							if err := mapToStruct(r, &flairText); err != nil {
								return err
							}

							emojiAndTexts := helper.SplitStringIntoStandardEmojisAndWords(flairText.Text)
							for _, emojiOrText := range emojiAndTexts {
								//  normal emoji
								if gomoji.ContainsEmoji(emojiOrText) {
									emojiTitle, err := helper.GetBestGuessedEmojiInfo(emojiOrText)
									if err != nil {
										return err
									}
									result := strings.ToLower(emojiTitle.CodePoint)
									joined := strings.Join(strings.Fields(result), "-")

									emoji, err := getEmojiByTitle(ctx, emojis, joined)
									if err != nil {
										return err
									}

									postFlairEmojis = append(postFlairEmojis, models.PostFlairEmoji{
										EmojiID:     emoji.ID,
										PostFlairID: postFlairID,
										OrderIndex:  int32(orderIndex),
									})
									orderIndex += 1
								} else {
									// text
									postFlairDescriptions = append(postFlairDescriptions, models.PostFlairDescription{
										PostFlairID: postFlairID,
										OrderIndex:  int32(orderIndex),
										Description: emojiOrText,
									})
									orderIndex += 1
								}
							}

						}
					default:
						return fmt.Errorf("Unknown type for rich text")
					}
				}
			}

			// user flairs
			for _, userFlair := range subreddit.UserFlairs {
				userFlairID := uuid.New()
				userFlairs = append(userFlairs, models.UserFlair{
					ID:              userFlairID,
					VoxsphereID:     voxsphere.ID,
					FullText:        userFlair.FullText,
					BackgroundColor: userFlair.BackgroundColor,
				})

				orderIndex := 0

				richTexts := userFlair.RichText
				for _, richtText := range richTexts {
					switch r := richtText.(type) {
					case map[string]any:
						_, found := r["url"]
						if found {
							var flairEmoji FlairEmoji
							if err := mapToStruct(r, &flairEmoji); err != nil {
								return err
							}
							if helper.IsCustomEmoji(flairEmoji.Text) {
								customEmoji, err := getCustomEmojiByTitleandVoxsphereID(
									ctx,
									customEmojis,
									flairEmoji.Text,
									voxsphere.ID,
								)
								if err != nil {
									return err
								}
								userFlairCustomEmojis = append(userFlairCustomEmojis, models.UserFlairCustomEmoji{
									CustomEmojiID: customEmoji.ID,
									UserFlairID:   userFlairID,
									OrderIndex:    int32(orderIndex),
								})
								orderIndex += 1
							}
						} else {
							var flairText FlairText
							if err := mapToStruct(r, &flairText); err != nil {
								return err
							}

							emojiAndTexts := helper.SplitStringIntoStandardEmojisAndWords(flairText.Text)
							for _, emojiOrText := range emojiAndTexts {
								//  normal emoji
								if gomoji.ContainsEmoji(emojiOrText) {
									emojiTitle, err := helper.GetBestGuessedEmojiInfo(emojiOrText)
									if err != nil {
										return err
									}
									result := strings.ToLower(emojiTitle.CodePoint)
									joined := strings.Join(strings.Fields(result), "-")

									emoji, err := getEmojiByTitle(ctx, emojis, joined)
									if err != nil {
										return err
									}

									userFlairEmojis = append(userFlairEmojis, models.UserFlairEmoji{
										EmojiID:     emoji.ID,
										UserFlairID: userFlairID,
										OrderIndex:  int32(orderIndex),
									})
									orderIndex += 1
								} else {
									// text
									userFlairDescriptions = append(userFlairDescriptions, models.UserFlairDescription{
										UserFlairID: userFlairID,
										OrderIndex:  int32(orderIndex),
										Description: emojiOrText,
									})
									orderIndex += 1
								}
							}

						}
					default:
						return fmt.Errorf("Unknown type for rich text")
					}
				}
			}

		}
	}

	// claim missing flairs
	for _, post := range postsJson {
		postFlair := post.LinkFlairText
		userFlair := post.AuthorFlairText

		postSubredditID := post.ID

		// find missing user flairs
		if len(userFlair) > 0 {
			foundUserFlair := false
			for _, subreddits := range subredditJson {
				for _, subreddit := range subreddits {
					if subreddit.ID == postSubredditID {
						for _, flair := range subreddit.UserFlairs {
							if userFlair == flair.FullText {
								foundUserFlair = true
							}
						}
					}
				}
			}

			if !foundUserFlair {
				voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, post.Subreddit)
				if err != nil {
					return err
				}

				wordsAndCustomEmojis := helper.SplitStringIntoCustomEmojisAndWords(userFlair)

				userFlairID := uuid.New()
				newUserFlair := models.UserFlair{
					ID:              userFlairID,
					VoxsphereID:     voxsphere.ID,
					FullText:        postFlair,
					BackgroundColor: "#000000",
				}
				userFlairs = append(userFlairs, newUserFlair)

				orderIndex := 0
				for _, wordsOrCustomEmoji := range wordsAndCustomEmojis {
					if helper.IsCustomEmoji(wordsOrCustomEmoji) {
						result := strings.ReplaceAll(wordsOrCustomEmoji, ":", "")
						newUserDescription := models.UserFlairDescription{
							UserFlairID: userFlairID,
							OrderIndex:  int32(orderIndex),
							Description: result,
						}
						userFlairDescriptions = append(userFlairDescriptions, newUserDescription)
						orderIndex += 1
					} else {
						wordsAndEmojis := helper.SplitStringIntoStandardEmojisAndWords(wordsOrCustomEmoji)
						for _, wordsOrEmoji := range wordsAndEmojis {
							if gomoji.ContainsEmoji(wordsOrEmoji) {
								emoji, err := helper.GetBestGuessedEmojiInfo(wordsOrEmoji)
								if err != nil {
									return err
								}
								result := strings.ToLower(emoji.CodePoint)
								joined := strings.Join(strings.Fields(result), "-")
								e, err := getEmojiByTitle(ctx, emojis, joined)
								if err != nil {
									return err
								}

								newUserEmoji := models.UserFlairEmoji{
									EmojiID:     e.ID,
									UserFlairID: userFlairID,
									OrderIndex:  int32(orderIndex),
								}
								userFlairEmojis = append(userFlairEmojis, newUserEmoji)
								orderIndex += 1
							} else {
								newUserDescription := models.UserFlairDescription{
									UserFlairID: userFlairID,
									OrderIndex:  int32(orderIndex),
									Description: wordsOrEmoji,
								}
								userFlairDescriptions = append(userFlairDescriptions, newUserDescription)
								orderIndex += 1
							}
						}
					}
				}
			}
		}

		// find missing post flairs
		if len(postFlair) > 0 {
			foundPostFlair := false
			for _, subreddits := range subredditJson {
				for _, subreddit := range subreddits {
					if subreddit.ID == postSubredditID {
						for _, flair := range subreddit.Flairs {
							if postFlair == flair.FullText {
								foundPostFlair = true
							}
						}
					}
				}
			}

			if !foundPostFlair {
				voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, post.Subreddit)
				if err != nil {
					return err
				}

				wordsAndCustomEmojis := helper.SplitStringIntoCustomEmojisAndWords(postFlair)

				postFlairID := uuid.New()
				newPostFlair := models.PostFlair{
					ID:              postFlairID,
					VoxsphereID:     voxsphere.ID,
					FullText:        postFlair,
					BackgroundColor: "#000000",
				}
				postFlairs = append(postFlairs, newPostFlair)

				orderIndex := 0
				for _, wordsOrCustomEmoji := range wordsAndCustomEmojis {
					if helper.IsCustomEmoji(wordsOrCustomEmoji) {
						result := strings.ReplaceAll(wordsOrCustomEmoji, ":", "")
						newPostFlairDescription := models.PostFlairDescription{
							PostFlairID: postFlairID,
							OrderIndex:  int32(orderIndex),
							Description: result,
						}
						postFlairDescriptions = append(postFlairDescriptions, newPostFlairDescription)
						orderIndex += 1
					} else {
						wordsAndEmojis := helper.SplitStringIntoStandardEmojisAndWords(wordsOrCustomEmoji)
						for _, wordsOrEmoji := range wordsAndEmojis {
							if gomoji.ContainsEmoji(wordsOrEmoji) {
								emoji, err := helper.GetBestGuessedEmojiInfo(wordsOrEmoji)
								if err != nil {
									return err
								}
								result := strings.ToLower(emoji.CodePoint)
								joined := strings.Join(strings.Fields(result), "-")
								e, err := getEmojiByTitle(ctx, emojis, joined)
								if err != nil {
									return err
								}

								newPostEmoji := models.PostFlairEmoji{
									EmojiID:     e.ID,
									PostFlairID: postFlairID,
									OrderIndex:  int32(orderIndex),
								}
								postFlairEmojis = append(postFlairEmojis, newPostEmoji)
								orderIndex += 1
							} else {
								newPostDescription := models.PostFlairDescription{
									PostFlairID: postFlairID,
									OrderIndex:  int32(orderIndex),
									Description: wordsOrEmoji,
								}
								postFlairDescriptions = append(postFlairDescriptions, newPostDescription)
								orderIndex += 1
							}
						}
					}
				}
			}
		}
	}

	postFlairRepo := postflairsrepo.NewRepo(db)

	wg := new(sync.WaitGroup)

	postFlairsCI := NewConcurrentInserter(len(postFlairs), func(postFlairs []models.PostFlair) error {
		if _, err := postFlairRepo.AddPostFlairs(ctx, postFlairs...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(postFlairs))

	go func() {
		for _, postflair := range postFlairs {
			postFlairsCI.ResC <- postflair
			wg.Done()
		}
	}()

	if err := postFlairsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	relationRepo := relationrepo.NewRepo(db)

	postFlairEmojisCI := NewConcurrentInserter(len(postFlairEmojis), func(postFlairEmojis []models.PostFlairEmoji) error {
		if _, err := relationRepo.LinkPostFlairEmojis(ctx, postFlairEmojis...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(postFlairEmojis))

	go func() {
		for _, postflairEmoji := range postFlairEmojis {
			postFlairEmojisCI.ResC <- postflairEmoji
			wg.Done()
		}
	}()

	if err := postFlairEmojisCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	postFlairCustomEmojisCI := NewConcurrentInserter(len(postFlairCustomEmojis), func(postFlairCustomEmojis []models.PostFlairCustomEmoji) error {
		if _, err := relationRepo.LinkPostFlairCustomEmojis(ctx, postFlairCustomEmojis...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(postFlairCustomEmojis))

	go func() {
		for _, postFlairCustomEmoji := range postFlairCustomEmojis {
			postFlairCustomEmojisCI.ResC <- postFlairCustomEmoji
			wg.Done()
		}
	}()

	if err := postFlairCustomEmojisCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	postFlairDescCI := NewConcurrentInserter(len(postFlairDescriptions), func(postFlairDescriptions []models.PostFlairDescription) error {
		if _, err := relationRepo.LinkPostFlairDescriptions(ctx, postFlairDescriptions...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(postFlairDescriptions))

	go func() {
		for _, postFlairDescription := range postFlairDescriptions {
			postFlairDescCI.ResC <- postFlairDescription
			wg.Done()
		}
	}()

	if err := postFlairDescCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userFlairRepo := userflairsrepo.NewRepo(db)

	userFlairsCI := NewConcurrentInserter(len(userFlairs), func(userFlairs []models.UserFlair) error {
		if _, err := userFlairRepo.AddUserFlairs(ctx, userFlairs...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(userFlairs))

	go func() {
		for _, userflair := range userFlairs {
			userFlairsCI.ResC <- userflair
			wg.Done()
		}
	}()

	if err := userFlairsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userFlairEmojisCI := NewConcurrentInserter(len(userFlairEmojis), func(userFlairEmojis []models.UserFlairEmoji) error {
		if _, err := relationRepo.LinkUserFlairEmojis(ctx, userFlairEmojis...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(userFlairEmojis))

	go func() {
		for _, userflairEmoji := range userFlairEmojis {
			userFlairEmojisCI.ResC <- userflairEmoji
			wg.Done()
		}
	}()

	if err := userFlairEmojisCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userFlairCustomEmojisCI := NewConcurrentInserter(len(userFlairCustomEmojis), func(userFlairCustomEmojis []models.UserFlairCustomEmoji) error {
		if _, err := relationRepo.LinkUserFlairCustomEmojis(ctx, userFlairCustomEmojis...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(userFlairCustomEmojis))

	go func() {
		for _, userFlairCustomEmoji := range userFlairCustomEmojis {
			userFlairCustomEmojisCI.ResC <- userFlairCustomEmoji
			wg.Done()
		}
	}()

	if err := userFlairCustomEmojisCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userFlairDescCI := NewConcurrentInserter(len(userFlairDescriptions), func(userFlairDescriptions []models.UserFlairDescription) error {
		if _, err := relationRepo.LinkUserFlairDescriptions(ctx, userFlairDescriptions...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(userFlairDescriptions))

	go func() {
		for _, userFlairDescription := range userFlairDescriptions {
			userFlairDescCI.ResC <- userFlairDescription
			wg.Done()
		}
	}()

	if err := userFlairDescCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	userRepo := user.NewRepo(db)
	users, err := userRepo.Users(ctx)
	if err != nil {
		return err
	}

	postRepo := postrepo.NewRepo(db)
	posts, err := postRepo.Posts(ctx)
	if err != nil {
		return err
	}

	user_user_flairs := make([]models.UserUserFlair, 0)
	post_post_flairs := make([]models.PostPostFlair, 0)

	for _, post := range postsJson {
		if post.LinkFlairText == "Moderator" {
			continue
		}

		userFlairText := post.AuthorFlairText
		idx1 := slices.IndexFunc(userFlairs, func(uf models.UserFlair) bool {
			return uf.FullText == userFlairText
		})

		postFlairText := post.LinkFlairText
		idx2 := slices.IndexFunc(postFlairs, func(pf models.PostFlair) bool {
			return pf.FullText == postFlairText
		})

		if idx1 != -1 {
			user, err := getuserByName(ctx, users, post.Author)
			if err != nil {
				return err
			}

			if !slices.ContainsFunc(user_user_flairs, func(uuf models.UserUserFlair) bool {
				return uuf.UserID == user.ID && uuf.UserFlairID == userFlairs[idx1].ID
			}) {
				user_user_flairs = append(user_user_flairs, models.UserUserFlair{
					UserID:      user.ID,
					UserFlairID: userFlairs[idx1].ID,
				})
			}
		}

		if idx2 != -1 {
			user, err := getuserByName(ctx, users, post.Author)
			if err != nil {
				return err
			}
			voxsphere, err := getvoxsphereByTitle(ctx, voxspheres, post.Subreddit)
			if err != nil {
				return err
			}

			post, err := getpostByTitleAndAuthorIDAndVoxsphereID(ctx, posts, post.Title, user.ID, voxsphere.ID)
			if err != nil {
				return err
			}

			if !slices.ContainsFunc(post_post_flairs, func(ppf models.PostPostFlair) bool {
				return ppf.PostID == post.ID && ppf.PostFlairID == postFlairs[idx2].ID
			}) {
				post_post_flairs = append(post_post_flairs, models.PostPostFlair{
					PostID:      post.ID,
					PostFlairID: postFlairs[idx2].ID,
				})
			}

		}

	}

	userUserFlairsCI := NewConcurrentInserter(len(user_user_flairs), func(userUserFlairs []models.UserUserFlair) error {
		if _, err := relationRepo.LinkUserUserFlairs(ctx, userUserFlairs...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(user_user_flairs))

	go func() {
		for _, userUserFlair := range user_user_flairs {
			userUserFlairsCI.ResC <- userUserFlair
			wg.Done()
		}
	}()

	if err := userUserFlairsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	postPostFlairsCI := NewConcurrentInserter(len(post_post_flairs), func(postPostFlairs []models.PostPostFlair) error {
		if _, err := relationRepo.LinkPostPostFlairs(ctx, postPostFlairs...); err != nil {
			return err
		}

		return nil
	}, 100*time.Millisecond, 500)

	wg.Add(len(post_post_flairs))

	go func() {
		for _, postPostFlair := range post_post_flairs {
			postPostFlairsCI.ResC <- postPostFlair
			wg.Done()
		}
	}()

	if err := postPostFlairsCI.Serve(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func Run() {
	// get db instance
	db, err := connectPostgres("postgres", "postgres", "127.0.0.1:5432", "voxpopuli")
	if err != nil {
		panic(err)
	}

	fileMap := make(map[string]string)

	fileMap["topics_json"] = ""
	fileMap["trophies_json"] = ""
	fileMap["awards_json"] = ""
	fileMap["subreddits_json"] = ""
	fileMap["posts_json"] = ""
	fileMap["users_json"] = ""

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// insert topics
	if err := insertTopics(ctx, db, fileMap["topics_json"]); err != nil {
		panic(err)
	}

	// insert trophies
	if err := insertTrophies(ctx, db, fileMap["trophies_json"]); err != nil {
		panic(err)
	}

	// insert awards
	if err := insertAwards(ctx, db, fileMap["awards_json"]); err != nil {
		panic(err)
	}

	// insert users
	if err := insertUsers(ctx, db, fileMap["users_json"]); err != nil {
		panic(err)
	}

	// insert voxspheres
	if err := insertVoxspheres(ctx, db, fileMap["subreddits_json"]); err != nil {
		panic(err)
	}

	// insert posts
	if err := insertPosts(ctx, db, fileMap["posts_json"]); err != nil {
		panic(err)
	}

	// insert emojis
	if err := insertEmojis(ctx, db, fileMap["subreddits_json"]); err != nil {
		panic(err)
	}

	// insert flairs
	if err := insertFlairs(ctx, db, fileMap["posts_json"], fileMap["subreddits_json"]); err != nil {
		panic(err)
	}
}
