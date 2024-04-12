-- name: GetShopPacks :many
SELECT 
  p.*, -- pack
  pf.id AS pack_file_id, pf.name AS pack_file_name, pf.path AS pack_file_path, -- pack file
  json_agg(pr.*) as pack_rarities
FROM packs AS p
LEFT JOIN files pf ON p.file_id = pf.id
INNER JOIN pack_rarities pr ON p.id = pr.pack_id
GROUP BY p.id, pf.id;