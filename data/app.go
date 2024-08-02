package data

func FindAccessTokenByAppId() (accessToken string, expiresIn int64, err error) {

	appId := "wx370126c8bcf8d00c"

	err = db.QueryRow("SELECT access_token, expires_in FROM mini WHERE app_id = $1", appId).Scan(&accessToken, &expiresIn)

	return
}

func UpdateAccessTokenAndExpiresIn(accessToken string, expiresIn int64) (err error) {

	appId := "wx370126c8bcf8d00c"

	_, err = db.Exec("UPDATE mini SET access_token = $1, expires_in = $2 WHERE app_id = $3", accessToken, expiresIn, appId)

	return
}
