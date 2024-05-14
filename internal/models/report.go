package models

type Report struct {
	AppId      string `json:"appId"`
	DeviceId   string `json:"deviceId"`
	AppVersion string `json:"appVersion"`
}

func CreateReport(payload Report, ip string) error {
	var exists bool

	existQuery := `SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1)`

	if err := db.Get(&exists, existQuery, payload.DeviceId); err != nil {
		return err
	}

	tx := db.MustBegin()

	if !exists {
		statementDevice := `INSERT INTO devices (id, ip) VALUES ($1, $2)`

		if _, err := tx.Exec(statementDevice, payload.DeviceId, ip); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	var existsRelation bool

	existRealtionQuery := `SELECT EXISTS(SELECT 1 FROM apps_devices WHERE appid = $1 AND deviceid = $2 AND ver = $3)`

	if err := db.Get(&existsRelation, existRealtionQuery, payload.AppId, payload.DeviceId, payload.AppVersion); err != nil {
		return err
	}

	if existsRelation {
		statementAppLinkUpdate := `UPDATE apps_devices SET usage = usage + 1 WHERE appid = $1 AND deviceid = $2 AND ver = $3`

		if _, err := tx.Exec(statementAppLinkUpdate, payload.AppId, payload.DeviceId, payload.AppVersion); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	} else {
		statementAppLink := `INSERT INTO apps_devices (appid, deviceid, ver, usage) VALUES ($1, $2, $3, 1)`

		if _, err := tx.Exec(statementAppLink, payload.AppId, payload.DeviceId, payload.AppVersion); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}
