package models

type Report struct {
	AppId      string `json:"appId"`
	DeviceId   string `json:"deviceId"`
	AppVersion string `json:"appVersion"`
}

func CreateReport(payload Report, ip string) error {
	statementDevice := `IF NOT EXISTS (SELECT 1 FROM devices WHERE id = $1)
									THEN INSERT INTO devices (id, ip) VALUES ($1, $2);
								END IF;`

	tx := db.MustBegin()

	if _, err := tx.Exec(statementDevice, payload.DeviceId, ip); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	statementAppLink := `INSERT INTO apps_devices (appid, deviceid, version, usage)
												VALUES ($1, $2, $3, 1)
												ON CONFLICT (appid, deviceid, version)
												DO UPDATE SET usage_count = apps_devices.usage + 1;`

	if _, err := tx.Exec(statementAppLink, payload.AppId, payload.DeviceId, payload.AppVersion); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}
