package main

import "time"

func UpdateUser(userInfo *UserInfo) error {
	baseSQL := `
		UPDATE
			users
		SET nike=$1, avatar=$2, phone=$3, updated_at=$4
		WHERE snowflake_id=$5`
	_, err := db.Exec(baseSQL, userInfo.Nike, userInfo.Avatar, userInfo.Phone, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	return nil
}

func RegisterUser(userInfo *UserInfo) error {

	baseSQL := `
		INSERT INTO 
			users 
			(snowflake_id, open_id, nike, avatar, phone, session_key, access_token, expires_in, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
	`
	_, err := db.Exec(baseSQL,
		userInfo.SnowflakeId,
		userInfo.OpenID,
		userInfo.Nike,
		userInfo.Avatar,
		userInfo.Phone,
		userInfo.SessionKey,
		userInfo.AccessToken,
		userInfo.ExpiresIn,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}
