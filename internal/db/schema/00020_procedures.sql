-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_auction_offer_duration()
  RETURNS TRIGGER 
  LANGUAGE PLPGSQL
  AS
$$
  BEGIN
    UPDATE auction_offers
    SET duration = duration + 300000
    WHERE id = NEW.auction_offer_id AND 
    extract(epoch from (created_at + duration * interval '1 millisecond' - Now() at TIME zone 'UTC')) between 0 and 300;

    RETURN NEW;
  END;
$$
-- +goose StatementEnd

-- +goose Down
DROP PROCEDURE IF EXISTS update_auction_offer_duration;