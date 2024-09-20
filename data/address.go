package data

import (
	"database/sql"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddAddress(tx *sql.Tx, addAddress *entity.AddAddress) error {
	_, err := tx.Exec("INSERT INTO address (snowflake_id, location, phone_number, is_default, consignee, longitude, latitude, detailed_address, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		addAddress.SnowflakeId,
		addAddress.Location,
		addAddress.PhoneNumber,
		addAddress.IsDefault,
		addAddress.Consignee,
		addAddress.Longitude,
		addAddress.Latitude,
		addAddress.DetailedAddress,
		addAddress.UserId,
	)

	return err
}

func GetAddressList(tx *sql.Tx, userId string) ([]entity.AddressList, error) {
	rows, err := tx.Query("SELECT snowflake_id, location, phone_number, is_default, consignee, longitude, latitude, detailed_address FROM address WHERE user_id=$1 AND deleted_at IS NULL ORDER BY is_default DESC, created_at DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addressList := make([]entity.AddressList, 0)

	for rows.Next() {
		var address entity.AddressList
		err := rows.Scan(&address.SnowflakeId, &address.Location, &address.PhoneNumber, &address.IsDefault, &address.Consignee, &address.Longitude, &address.Latitude, &address.DetailedAddress)
		if err != nil {
			return nil, err
		}
		addressList = append(addressList, address)
	}
	return addressList, nil
}

func DeleteAddress(tx *sql.Tx, addressId string, userId string) error {
	_, err := tx.Exec("UPDATE address SET deleted_at=$1 WHERE snowflake_id=$2 AND user_id = $3", time.Now(), addressId, userId)
	return err
}

func UpdateAddress(tx *sql.Tx, updateAddress *entity.UpdateAddress) error {
	_, err := tx.Exec("UPDATE address SET location=$1, phone_number=$2, is_default=$3, consignee=$4, longitude=$5, latitude=$6, detailed_address=$7 WHERE snowflake_id=$8 AND user_id=$9",
		updateAddress.Location,
		updateAddress.PhoneNumber,
		updateAddress.IsDefault,
		updateAddress.Consignee,
		updateAddress.Longitude,
		updateAddress.Latitude,
		updateAddress.DetailedAddress,
		updateAddress.SnowflakeId,
		updateAddress.UserId,
	)
	return err
}

func UpdateAddressIsDefault(tx *sql.Tx, userId string) error {
	_, err := tx.Exec("UPDATE address SET is_default=0 WHERE user_id=$1", userId)
	return err
}

func FindByAddressSnowflakeId(tx *sql.Tx, snowflakeId string, userId string) (*entity.AddressList, error) {
	row := tx.QueryRow("SELECT snowflake_id, location, phone_number, is_default, consignee, longitude, latitude, detailed_address FROM address WHERE snowflake_id=$1 AND user_id=$2 AND deleted_at IS NULL", snowflakeId, userId)
	var address entity.AddressList
	err := row.Scan(&address.SnowflakeId, &address.Location, &address.PhoneNumber, &address.IsDefault, &address.Consignee, &address.Longitude, &address.Latitude, &address.DetailedAddress)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
