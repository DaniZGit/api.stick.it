package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/DaniZGit/api.stick.it/internal/ws"
	"github.com/go-co-op/gocron/v2"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func markCompletedAuctionsTask(queries *database.Queries, hubs *ws.HubModels) (gocron.JobDefinition, gocron.Task) {
	cronDuration := gocron.DurationJob(
		1*time.Minute,
	)

	cronTask := gocron.NewTask(
		func() {
			// do things
			auctionOffers, err := queries.MarkCompletedAuctionOffers(context.Background())
			if err != nil {
				fmt.Println("Error while marking completed auction offers", err)
			}

			for _, auctionOffer := range auctionOffers {
				// get won sticker
				if err != nil {
					fmt.Println("Error while getting user_sticker from DB", err)
				}

				// get auctioned user_sticker data
				userSticker, err := queries.GetUserSticker(context.Background(), auctionOffer.UserStickerID)
				if err != nil && err != pgx.ErrNoRows {
					fmt.Println("Error while getting user_sticker from DB", err)
				}

				// get last auction bid
				auctionBid, auctionBidErr := queries.GetLatestAuctionBid(context.Background(), auctionOffer.ID)
				if err != nil && err != pgx.ErrNoRows {
					fmt.Println("Error while getting last auction bid from DB", err)
				}

				userID := uuid.UUID{}
				stickerID := uuid.UUID{}
				bid := 0
				if auctionBidErr == pgx.ErrNoRows {
					// no bids happened, return the sticker to the auctioneer
					userID = userSticker.UserID
					stickerID = userSticker.StickerID
					bid = int(auctionOffer.StartingBid)
				} else {
					// add sticker to the winner
					userID = auctionBid.UserID
					stickerID = auctionBid.StickerID
					bid = int(auctionBid.Bid)
				}

				// add the sticker to the user
				_, err = queries.CreateUserSticker(context.Background(), database.CreateUserStickerParams{
					ID: uuid.Must(uuid.NewV4()),
					UserID: userID,
					StickerID: stickerID,
					Amount: 1,
				})
				if err != nil {
					fmt.Println("Error while adding sticker to the user in DB", err)
				}

				// add tokens to the auctioneer
				_, err = queries.IncrementUserTokens(context.Background(), database.IncrementUserTokensParams{
					ID: userSticker.UserID,
					Tokens: int64(bid),
				})
				if err != nil {
					fmt.Println("Error while trying to increment auctioneer tokens in DB", err)
				}

				// broadcast the completed auction offer to all clients
				event := ws.AuctionEvent{
					Type: ws.AuctionEventTypeCompleted,
					Payload: data.AuctionOffer{
						ID: auctionOffer.ID,
					},
				}
				data, err := json.Marshal(event)
				if err != nil {
					fmt.Println("Failed to broadcast the auction bid event", err)
				} else {
					hubs.AuctionHub.Broadcast <- data
				}
			}
		},
	)
	return cronDuration, cronTask
}
