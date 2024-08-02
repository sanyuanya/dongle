package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func UpdateUser(userInfo *entity.UserInfo) error {
	baseSQL := `
		UPDATE
			users
		SET nick=$1, avatar=$2, phone=$3, updated_at=$4
		WHERE snowflake_id=$5`
	_, err := db.Exec(baseSQL, userInfo.Nick, userInfo.Avatar, userInfo.Phone, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	return nil
}

// func RegisterUser(userInfo *entity.UserInfo) error {

// 	baseSQL := `
// 		INSERT INTO
// 			users
// 			(snowflake_id, open_id, nick, avatar, phone, session_key, created_at, updated_at)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
// 	`
// 	_, err := db.Exec(baseSQL,
// 		userInfo.SnowflakeId,
// 		userInfo.OpenID,
// 		userInfo.Nick,
// 		userInfo.Avatar,
// 		userInfo.Phone,
// 		userInfo.SessionKey,
// 		time.Now(),
// 		time.Now(),
// 	)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetUserDetailBySnowflakeID(snowflakeId int64) (*entity.UserDetail, error) {
	baseSQL := `
		SELECT 
			snowflake_id, open_id, nick, avatar, phone, integral, shipments, province, city, district, id_card, company_name, job, alipay_account, created_at, updated_at, session_key
		FROM
			users
		WHERE
			snowflake_id=$1
	`
	row := db.QueryRow(baseSQL, snowflakeId)
	userDetail := &entity.UserDetail{}
	err := row.Scan(
		&userDetail.SnowflakeId,
		&userDetail.OpenID,
		&userDetail.Nick,
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
		&userDetail.CreatedAt,
		&userDetail.UpdatedAt,
		&userDetail.SessionKey,
	)
	if err != nil {
		return nil, err
	}
	return userDetail, nil
}

func UpdateUserInfo(userInfo *entity.SetUserInfoRequest) error {
	baseSQL := `
		UPDATE
			users
		SET nick=$1, avatar=$2, phone=$3, id_card=$4, province=$5, city=$6, district=$7, company_name=$8, job=$9, updated_at=$10
		WHERE snowflake_id=$11
	`
	_, err := db.Exec(baseSQL, userInfo.Nick, userInfo.Avatar, userInfo.Phone, userInfo.IDCard, userInfo.Province, userInfo.City, userInfo.District, userInfo.CompanyName, userInfo.Job, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	return nil
}

func FindPhoneNumberContext(phone string) (int64, error) {

	var snowflakeId int64

	baseQueryPhone := "SELECT snowflake_id FROM users WHERE phone = $1"

	if err := db.QueryRow(baseQueryPhone, phone).Scan(&snowflakeId); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return snowflakeId, nil
}

func UpdateUserIntegralAndShipments(snowflakeId, integral, shipments int64) error {
	baseSQL := `
		UPDATE
			users
		SET integral=integral+$1, shipments=shipments+$2, updated_at=$3
		WHERE snowflake_id=$4
	`
	_, err := db.Exec(baseSQL, integral, shipments, time.Now(), snowflakeId)

	if err != nil {
		return err
	}
	return nil
}

func ImportUserInfo(importUserInfo *entity.ImportUserInfo) error {
	baseSQL := "INSERT INTO users (nick, phone, province, city, shipments, integral, snowflake_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := db.Exec(baseSQL, importUserInfo.Nick, importUserInfo.Phone, importUserInfo.Province, importUserInfo.City, importUserInfo.Shipments, importUserInfo.Integral, importUserInfo.SnowflakeId, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func GetUserPageCount(page *entity.UserPageListRequest) (int64, error) {
	baseSQL := `
		SELECT 
			COUNT(1)
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	executeParams := []interface{}{}
	paramIndex := 1
	// 判断是否有查询条件
	if page.IsWhite != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND is_white=$%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.IsWhite)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (nick LIKE $%d OR phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	var count int64
	err := db.QueryRow(baseSQL, executeParams...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetUserPageList(page *entity.UserPageListRequest) ([]*entity.UserPageListResponse, error) {

	baseSQL := `
		SELECT 
			snowflake_id, nick, avatar, phone, integral, shipments, province, city, district, id_card, company_name, job, alipay_account, is_white
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	executeParams := []interface{}{}
	paramIndex := 1
	// 判断是否有查询条件
	if page.IsWhite != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND is_white=$%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.IsWhite)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (nick LIKE $%d OR phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := db.Query(baseSQL, executeParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userPageList := make([]*entity.UserPageListResponse, 0)
	for rows.Next() {
		user := &entity.UserPageListResponse{}
		err := rows.Scan(
			&user.SnowflakeId,
			&user.Nick,
			&user.Avatar,
			&user.Phone,
			&user.Integral,
			&user.Shipments,
			&user.Province,
			&user.City,
			&user.District,
			&user.IDCard,
			&user.CompanyName,
			&user.Job,
			&user.AlipayAccount,
			&user.IsWhite,
		)
		if err != nil {
			return nil, err
		}
		userPageList = append(userPageList, user)
	}

	return userPageList, nil
}

func AddWhite(whiteList *entity.AddWhiteRequest) error {

	baseSQL := `
		UPDATE
			users
		SET is_white=1
		WHERE snowflake_id = $1
	`

	for _, snowflakeId := range whiteList.WhiteList {
		_, err := db.Exec(baseSQL, snowflakeId)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsWhite(snowflakeId int64) error {
	baseSQL := `
		SELECT 
			is_white
		FROM
			users
		WHERE
			snowflake_id=$1 AND is_white=1
	`
	var isWhite bool
	err := db.QueryRow(baseSQL, snowflakeId).Scan(&isWhite)
	if err != nil {
		return fmt.Errorf("查询用户是否是白名单失败: %v", err)
	}
	return nil
}

func IsIntegralWithdraw(snowflakeId, integral int64) (bool, error) {
	baseSQL := `
		SELECT 
			integral
		FROM
			users
		WHERE
			snowflake_id=$1
	`
	var userIntegral int64
	err := db.QueryRow(baseSQL, snowflakeId).Scan(&userIntegral)
	if err != nil {
		return false, fmt.Errorf("查询用户积分失败: %v", err)
	}

	if userIntegral < integral {
		return false, fmt.Errorf("积分不足: 当前积分 %d", integral)
	}

	return true, nil
}

func UpdateUserAlipayAccountBySnowflakeID(snowflakeId int64, alipayAccount string) error {
	baseSQL := `
		UPDATE
			users
		SET alipay_account=$1, updated_at=$2
		WHERE snowflake_id=$3
	`
	_, err := db.Exec(baseSQL, alipayAccount, time.Now(), snowflakeId)
	if err != nil {
		return fmt.Errorf("更新支付宝账号失败: %v", err)
	}
	return nil
}

func FindOpenId(openid string) (int64, error) {

	baseSQL := `
		SELECT 
			snowflake_id
		FROM
			users
		WHERE
			open_id=$1
	`
	var snowflakeId int64
	err := db.QueryRow(baseSQL, openid).Scan(&snowflakeId)
	if err != nil {
		return 0, err
	}
	return snowflakeId, nil
}

func RegisterUser(registerUserRequest *entity.RegisterUserRequest) error {

	baseSQL := `
		INSERT INTO
			users
			(snowflake_id, open_id, session_key, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	_, err := db.Exec(baseSQL,
		registerUserRequest.SnowflakeId,
		registerUserRequest.OpenId,
		registerUserRequest.SessionKey,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateSessionKey(openid, sessionKey string) error {
	baseSQL := `
		UPDATE
			users
		SET session_key=$1, updated_at=$2
		WHERE open_id=$3
	`
	_, err := db.Exec(baseSQL, sessionKey, time.Now(), openid)
	if err != nil {
		return err
	}

	return nil
}
