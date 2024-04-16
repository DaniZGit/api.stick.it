-- +goose Up
-- users
ALTER TABLE users ADD CONSTRAINT users_file_id_files_id FOREIGN KEY (file_id) REFERENCES files(id);
ALTER TABLE users ADD CONSTRAINT users_role_id_roles_id FOREIGN KEY (role_id) REFERENCES roles(id);
ALTER TABLE users ADD CONSTRAINT users_amount_positive CHECK (tokens >= 0);

-- albums
ALTER TABLE albums ADD CONSTRAINT albums_file_id_files_id FOREIGN KEY (file_id) REFERENCES files(id);

-- pages
ALTER TABLE pages ADD CONSTRAINT pages_album_id_albums_id FOREIGN KEY (album_id) REFERENCES albums(id);
ALTER TABLE pages ADD CONSTRAINT pages_file_id_files_id FOREIGN KEY (file_id) REFERENCES files(id);

-- stickers
ALTER TABLE stickers ADD CONSTRAINT stickers_rarity_id_rarities_id FOREIGN KEY (rarity_id) REFERENCES rarities(id);
ALTER TABLE stickers ADD CONSTRAINT stickers_file_id_files_id FOREIGN KEY (file_id) REFERENCES files(id);
ALTER TABLE stickers ADD CONSTRAINT stickers_page_id_pages_id FOREIGN KEY (page_id) REFERENCES pages(id);
ALTER TABLE stickers ADD CONSTRAINT stickers_sticker_id_stickers_id FOREIGN KEY (sticker_id) REFERENCES stickers(id);

-- user_stickers
ALTER TABLE user_stickers ADD CONSTRAINT user_stickers_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_stickers ADD CONSTRAINT user_stickers_sticker_id_stickers_id FOREIGN KEY (sticker_id) REFERENCES stickers(id);

-- packs
ALTER TABLE packs ADD CONSTRAINT packs_album_id_albums_id FOREIGN KEY (album_id) REFERENCES albums(id);
ALTER TABLE packs ADD CONSTRAINT packs_file_id_files_id FOREIGN KEY (file_id) REFERENCES files(id);
ALTER TABLE packs ADD CONSTRAINT packs_rarity_id_rarities_id FOREIGN KEY (rarity_id) REFERENCES rarities(id);

-- pack_rarities
ALTER TABLE pack_rarities ADD CONSTRAINT pack_rarities_pack_id_packs_id FOREIGN KEY (pack_id) REFERENCES packs(id) ON DELETE CASCADE;
ALTER TABLE pack_rarities ADD CONSTRAINT pack_rarities_rarity_id_rarities_id FOREIGN KEY (rarity_id) REFERENCES rarities(id);

-- user_packs
ALTER TABLE user_packs ADD CONSTRAINT user_packs_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_packs ADD CONSTRAINT user_packs_pack_id_packs_id FOREIGN KEY (pack_id) REFERENCES packs(id);
ALTER TABLE user_packs ADD CONSTRAINT user_packs_unique UNIQUE(user_id, pack_id);

-- opened_packs
ALTER TABLE opened_packs ADD CONSTRAINT opened_packs_pack_id_packs_id FOREIGN KEY (pack_id) REFERENCES packs(id);
ALTER TABLE opened_packs ADD CONSTRAINT opened_packs_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

-- opened_pack_stickers
ALTER TABLE opened_pack_stickers ADD CONSTRAINT opened_pack_stickers_sticker_id_stickers_id FOREIGN KEY (sticker_id) REFERENCES stickers(id);
ALTER TABLE opened_pack_stickers ADD CONSTRAINT opened_pack_stickers_opened_pack_id_opened_packs_id FOREIGN KEY (sticker_id) REFERENCES opened_packs(id);

-- auction_offers
ALTER TABLE auction_offers ADD CONSTRAINT auction_offers_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE auction_offers ADD CONSTRAINT auction_offers_sticker_id_stickers_id FOREIGN KEY (sticker_id) REFERENCES stickers(id);

-- auction_bids
ALTER TABLE auction_bids ADD CONSTRAINT auction_bids_auction_offer_id_auction_offers_id FOREIGN KEY (auction_offer_id) REFERENCES auction_bids(id);
ALTER TABLE auction_bids ADD CONSTRAINT auction_bids_user_id_users_id FOREIGN KEY (user_id) REFERENCES auction_bids(id);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_file_id_files_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_id_roles_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_amount_positive;
ALTER TABLE albums DROP CONSTRAINT IF EXISTS albums_file_id_files_id;
ALTER TABLE pages DROP CONSTRAINT IF EXISTS pages_album_id_albums_id;
ALTER TABLE pages DROP CONSTRAINT IF EXISTS pages_file_id_files_id;
ALTER TABLE stickers DROP CONSTRAINT IF EXISTS stickers_rarity_id_rarities_id;
ALTER TABLE stickers DROP CONSTRAINT IF EXISTS stickers_file_id_files_id;
ALTER TABLE stickers DROP CONSTRAINT IF EXISTS stickers_page_id_pages_id;
ALTER TABLE stickers DROP CONSTRAINT IF EXISTS stickers_sticker_id_stickers_id;
ALTER TABLE user_stickers DROP CONSTRAINT IF EXISTS user_stickers_user_id_users_id;
ALTER TABLE user_stickers DROP CONSTRAINT IF EXISTS user_stickers_sticker_id_stickers_id;
ALTER TABLE packs DROP CONSTRAINT IF EXISTS packs_album_id_albums_id;
ALTER TABLE packs DROP CONSTRAINT IF EXISTS packs_file_id_files_id;
ALTER TABLE packs DROP CONSTRAINT IF EXISTS packs_rarity_id_rarities_id;
ALTER TABLE pack_rarities DROP CONSTRAINT IF EXISTS pack_rarities_pack_id_packs_id;
ALTER TABLE pack_rarities DROP CONSTRAINT IF EXISTS pack_rarities_rarity_id_rarities_id;
ALTER TABLE user_packs DROP CONSTRAINT IF EXISTS user_packs_user_id_users_id;
ALTER TABLE user_packs DROP CONSTRAINT IF EXISTS user_packs_pack_id_packs_id;
ALTER TABLE user_packs DROP CONSTRAINT IF EXISTS user_packs_unique;
ALTER TABLE opened_packs DROP CONSTRAINT IF EXISTS opened_packs_pack_id_packs_id;
ALTER TABLE opened_packs DROP CONSTRAINT IF EXISTS opened_packs_user_id_users_id;
ALTER TABLE opened_pack_stickers DROP CONSTRAINT IF EXISTS opened_pack_stickers_sticker_id_stickers_id;
ALTER TABLE opened_pack_stickers DROP CONSTRAINT IF EXISTS opened_pack_stickers_opened_pack_id_opened_packs_id;
ALTER TABLE auction_offers DROP CONSTRAINT IF EXISTS auction_offers_user_id_users_id;
ALTER TABLE auction_offers DROP CONSTRAINT IF EXISTS auction_offers_sticker_id_stickers_id;
ALTER TABLE auction_bids DROP CONSTRAINT IF EXISTS auction_bids_auction_offer_id_auction_offers_id;
ALTER TABLE auction_bids DROP CONSTRAINT IF EXISTS auction_bids_user_id_users_id;
