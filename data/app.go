package data

import "database/sql"

func FindAccessTokenByAppId(tx *sql.Tx) (accessToken string, expiresIn int64, err error) {

	appId := "wx370126c8bcf8d00c"

	err = tx.QueryRow("SELECT access_token, expires_in FROM mini WHERE app_id = $1", appId).Scan(&accessToken, &expiresIn)

	return
}

func UpdateAccessTokenAndExpiresIn(tx *sql.Tx, accessToken string, expiresIn int64) (err error) {

	appId := "wx370126c8bcf8d00c"

	_, err = tx.Exec("UPDATE mini SET access_token = $1, expires_in = $2 WHERE app_id = $3", accessToken, expiresIn, appId)

	return
}
