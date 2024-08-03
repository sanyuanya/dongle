package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func UpdateUserBySnowflakeId(userInfo *entity.UserInfo) error {
	baseSQL := `
		UPDATE
			users
		SET nick=$1, avatar=$2, phone=$3, updated_at=$4
		WHERE snowflake_id=$5 AND deleted_at IS NULL`
	_, err := db.Exec(baseSQL, userInfo.Nick, userInfo.Avatar, userInfo.Phone, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	return nil
}

func GetUserDetailBySnowflakeID(snowflakeId string) (*entity.UserDetail, error) {
	baseSQL := `
		SELECT 
			snowflake_id, openid, nick, avatar, phone, integral, shipments, province, city, district, id_card, company_name, job, alipay_account, created_at, updated_at, session_key, is_white, withdrawable_points
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
		&userDetail.IsWhite,
		&userDetail.WithdrawablePoints,
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
		WHERE snowflake_id=$11 AND deleted_at IS NULL
	`
	result, err := db.Exec(baseSQL, userInfo.Nick, userInfo.Avatar, userInfo.Phone, userInfo.IDCard, userInfo.Province, userInfo.City, userInfo.District, userInfo.CompanyName, userInfo.Job, time.Now(), userInfo.SnowflakeId)

	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}

	if row == 0 {
		return fmt.Errorf("更新用户信息失败: %v", "未找到用户")
	}

	return nil
}

func FindPhoneNumberContext(phone string) (int64, error) {

	var snowflakeId int64

	baseQueryPhone := "SELECT snowflake_id FROM users WHERE phone = $1 AND deleted_at IS NULL"

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
		WHERE snowflake_id=$4 AND deleted_at IS NULL
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
			snowflake_id, nick, avatar, phone, integral, shipments, province, city, district, id_card, company_name, job, alipay_account, is_white, withdrawable_points
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
			&user.WithdrawablePoints,
		)
		if err != nil {
			return nil, err
		}
		userPageList = append(userPageList, user)
	}

	return userPageList, nil
}

func SetUpWhiteRequest(whiteList *entity.SetUpWhiteRequest) error {

	baseSQL := `
		UPDATE
			users
		SET is_white=$1, updated_at=$2
		WHERE snowflake_id = $3 AND deleted_at IS NULL
	`
	for _, snowflakeId := range whiteList.WhiteList {
		_, err := db.Exec(baseSQL, whiteList.Status, time.Now(), snowflakeId)
		if err != nil {
			return err
		}
	}

	return nil
}

func IsWhite(snowflakeId string) error {
	baseSQL := `
		SELECT 
			is_white
		FROM
			users
		WHERE
			snowflake_id=$1 AND is_white=1 AND deleted_at IS NULL
	`
	var isWhite bool
	err := db.QueryRow(baseSQL, snowflakeId).Scan(&isWhite)
	if err != nil {
		return fmt.Errorf("查询用户是否是白名单失败: %v", err)
	}
	return nil
}

func IsWithdrawablePoints(snowflakeId, integral int64) (bool, error) {
	baseSQL := `
		SELECT 
			integral
		FROM
			users
		WHERE
			snowflake_id=$1 AND deleted_at IS NULL
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

func IsIntegralWithdraw(snowflakeId string, integral int64) error {
	baseSQL := `
		SELECT 
			integral
		FROM
			users
		WHERE
			snowflake_id=$1 AND deleted_at IS NULL
	`
	var userIntegral int64
	err := db.QueryRow(baseSQL, snowflakeId).Scan(&userIntegral)
	if err != nil {
		return fmt.Errorf("查询用户积分失败: %v", err)
	}

	if userIntegral < integral {
		return fmt.Errorf("积分不足: 当前积分 %d", integral)
	}

	return nil
}

func UpdateUserAlipayAccountBySnowflakeId(snowflakeId string, alipayAccount string) error {
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

func FindOpenId(openid string) (string, error) {

	baseSQL := `
		SELECT 
			snowflake_id
		FROM
			users
		WHERE
			openid=$1 AND deleted_at IS NULL
	`
	var snowflakeId string
	err := db.QueryRow(baseSQL, openid).Scan(&snowflakeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return snowflakeId, nil
}

func RegisterUser(registerUserRequest *entity.RegisterUserRequest) error {

	baseSQL := `
		INSERT INTO
			users
			(snowflake_id, openid, session_key, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
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
		WHERE openid=$3 AND deleted_at IS NULL
	`
	_, err := db.Exec(baseSQL, sessionKey, time.Now(), openid)
	if err != nil {
		return err
	}

	return nil
}

func FindUserByPhone(phone string) (*entity.UserInfoReplace, error) {
	baseSQL := `
		SELECT 
			snowflake_id, nick, phone, province, city, shipments, integral, is_white
		FROM
			users
		WHERE
			phone=$1 AND deleted_at IS NULL
	`

	userInfoReplace := &entity.UserInfoReplace{}

	err := db.QueryRow(baseSQL, phone).Scan(
		&userInfoReplace.SnowflakeId,
		&userInfoReplace.Nick,
		&userInfoReplace.Phone,
		&userInfoReplace.Province,
		&userInfoReplace.City,
		&userInfoReplace.Shipments,
		&userInfoReplace.Integral,
		&userInfoReplace.IsWhite,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return userInfoReplace, nil
}

func UserInfoReplace(userInfoReplace *entity.UserInfoReplace, snowflakeId string) error {
	baseSQL := `
		UPDATE
			users
		SET nick=$1, phone=$2, province=$3, city=$4, shipments=$5, integral=$6, is_white=$7, updated_at=$8
		WHERE snowflake_id=$9
	`
	_, err := db.Exec(baseSQL,
		userInfoReplace.Nick,
		userInfoReplace.Phone,
		userInfoReplace.Province,
		userInfoReplace.City,
		userInfoReplace.Shipments,
		userInfoReplace.Integral,
		userInfoReplace.IsWhite,
		time.Now(),
		snowflakeId,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(snowflakeId string) error {
	baseSQL := `
		UPDATE
			users
		SET deleted_at=$1
		WHERE snowflake_id=$2
	`
	_, err := db.Exec(baseSQL, time.Now(), snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserApiToken(snowflakeId string, apiToken string) error {
	baseSQL := `
		UPDATE
			users
		SET api_token=$1, updated_at=$2
		WHERE snowflake_id=$3
	`
	_, err := db.Exec(baseSQL, apiToken, time.Now(), snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func DeductUserIntegralAndWithdrawablePointsBySnowflakeId(snowflakeId string, integral int64) error {
	baseSQL := `
		UPDATE
			users
		SET integral=integral-$1, withdrawable_points=withdrawable_points-$1, updated_at=$2
		WHERE snowflake_id=$3
	`
	_, err := db.Exec(baseSQL, integral, time.Now(), snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func AddIntegralAndWithdrawablePointsBySnowflakeId(snowflakeId string, integral int64) error {
	baseSQL := `
		UPDATE
			users
		SET integral=integral+$1, withdrawable_points=withdrawable_points+$1, updated_at=$2
		WHERE snowflake_id=$3
	`
	_, err := db.Exec(baseSQL, integral, time.Now(), snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserDetail(payload *entity.UpdateUserDetailRequest) error {
	baseSQL := `
		UPDATE
			users
		SET nick=$1, phone=$2, province=$3, city=$4, district=$5, company_name=$8, job=$9, is_white=$10, shipments=$11, integral=$12, withdrawable_points=$13, updated_at=$14
		WHERE snowflake_id=$15
	`
	result, err := db.Exec(baseSQL, payload.Nick, payload.Phone, payload.Province, payload.City, payload.District, payload.CompanyName, payload.Job, payload.IsWhite, payload.Shipments, payload.Integral, payload.WithdrawablePoints, time.Now(), payload.SnowflakeId)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}

	if row == 0 {
		return fmt.Errorf("更新用户信息失败: %v", "未找到用户")
	}

	return nil
}
