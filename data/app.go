package data

func FindAppId() (accessToken string, expiresIn int64, err error) {

	appId := "wx370126c8bcf8d00c"

	err = db.QueryRow("SELECT access_token, expires_in FROM mini WHERE app_id = ?", appId).Scan(&accessToken, &expiresIn)

	return
}

func UpdateAccessTokenAndExpiresIn(accessToken string, expiresIn int64) (err error) {

	appId := "wx370126c8bcf8d00c"

	_, err = db.Exec("UPDATE mini SET access_token = ?, expires_in = ? WHERE app_id = ?", accessToken, expiresIn, appId)

	return
}
