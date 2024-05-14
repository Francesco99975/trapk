-- Create a trigger function to update the updated column
CREATE OR REPLACE FUNCTION update_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Create a macro to apply the trigger to a table
CREATE OR REPLACE FUNCTION apply_update_trigger(table_name TEXT)
RETURNS VOID AS $$
BEGIN
 IF NOT EXISTS (
    SELECT 1
    FROM information_schema.triggers
    WHERE trigger_schema = 'public'
      AND trigger_name = format('trigger_update_updated_%I', table_name)
  ) THEN
    EXECUTE format('
        CREATE TRIGGER trigger_update_updated_%I
        BEFORE UPDATE ON %I
        FOR EACH ROW
        EXECUTE FUNCTION update_updated()
    ', table_name, table_name);
  END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS apps (
  id TEXT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(appid)
);

SELECT apply_update_trigger('apps');

CREATE TABLE IF NOT EXISTS devices (
  id TEXT NOT NULL,
  ip TEXT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(deviceid)
);

SELECT apply_update_trigger('devices');

CREATE TABLE IF NOT EXISTS apps_devices (
  appid TEXT NOT NULL,
  deviceid TEXT NOT NULL,
  usage INT NOT NULL DEFAULT 0,
  version_TEXT TEXT NOT NULL,
  CONSTRAINT fk_ad
  FOREIGN KEY (appid)
  REFERENCES apps(id),
  CONSTRAINT fk_ru
  FOREIGN KEY (deviceid)
  REFERENCES devices(id),
  PRIMARY KEY(appid, deviceid)
);
