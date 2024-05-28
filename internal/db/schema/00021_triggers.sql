-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER on_auction_bid_insert_trigger
AFTER INSERT ON auction_bids FOR EACH ROW
EXECUTE PROCEDURE update_auction_offer_duration();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS on_auction_bid_insert_trigger ON auction_bids;