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

type UserDetail struct {
	SnowflakeId   int64  `json:"snowflake_id"`
	OpenID        string `json:"open_id"`
	Nike          string `json:"nike"`
	Avatar        string `json:"avatar"`
	Phone         string `json:"phone"`
	Integral      int    `json:"integral"`
	Shipments     int    `json:"shipments"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	IDCard        string `json:"id_card"`
	CompanyName   string `json:"company_name"`
	Job           string `json:"job"`
	AlipayAccount string `json:"alipay_account"`
}

func GetUserDetailBySnowflakeID(snowflakeID int64) (*UserDetail, error) {
	baseSQL := `
		SELECT 
			snowflake_id, open_id, nike, avatar, phone, integral, shipments, province, city, district, id_card, company_name, job, alipay_account
		FROM
			users
		WHERE
			snowflake_id=$1
	`
	row := db.QueryRow(baseSQL, snowflakeID)
	userDetail := &UserDetail{}
	err := row.Scan(
		&userDetail.SnowflakeId,
		&userDetail.OpenID,
		&userDetail.Nike,
		&userDetail.Avatar,
		&userDetail.Phone,
		&userDetail.Integral,
		&userDetail.Shipments,
		&userDetail.Province,
		&userDetail.City,
		&userDetail.District,
		&userDetail.IDCard,
		&userDetail.CompanyName,
		&userDetail.Job,
		&userDetail.AlipayAccount,
	)
	if err != nil {
		return nil, err
	}
	return userDetail, nil
}

func UpdateUserInfo(userInfo *SetUserInfoRequest) error {
	baseSQL := `
		UPDATE
			users
		SET nike=$1, avatar=$2, phone=$3, id_card=$4, province=$5, city=$6, district=$7, company_name=$8, job=$9, updated_at=$10
		WHERE snowflake_id=$11
	`
	_, err := db.Exec(baseSQL, userInfo.Nike, userInfo.Avatar, userInfo.Phone, userInfo.IDCard, userInfo.Province, userInfo.City, userInfo.District, userInfo.CompanyName, userInfo.Job, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	return nil
}
