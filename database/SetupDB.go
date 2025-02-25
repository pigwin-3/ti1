package database

import (
	"fmt"
	"ti1/config"
)

func SetupDB() error {
	fmt.Println("Setting up the database...")

	// Connect to PostgreSQL
	db, err := config.ConnectToPostgreSQL()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer config.DisconnectFromPostgreSQL(db)

	// Create sequences if they do not exist
	sequences := []string{
		"CREATE SEQUENCE IF NOT EXISTS public.calls_id_seq",
		"CREATE SEQUENCE IF NOT EXISTS public.estimatedvehiclejourney_id_seq",
		"CREATE SEQUENCE IF NOT EXISTS public.servicedelivery_id_seq",
	}

	for _, seq := range sequences {
		_, err := db.Exec(seq)
		if err != nil {
			return fmt.Errorf("failed to create sequence: %w", err)
		}
	}

	// Check if tables exist and have the correct structure
	tables := map[string]string{
		"calls": `CREATE TABLE IF NOT EXISTS public.calls (
			id BIGINT PRIMARY KEY DEFAULT nextval('public.calls_id_seq'),
			estimatedvehiclejourney BIGINT,
			"order" INTEGER,
			stoppointref VARCHAR,
			aimeddeparturetime TIMESTAMP,
			expecteddeparturetime TIMESTAMP,
			aimedarrivaltime TIMESTAMP,
			expectedarrivaltime TIMESTAMP,
			cancellation VARCHAR,
			actualdeparturetime TIMESTAMP,
			actualarrivaltime TIMESTAMP,
			estimated_data JSON,
			recorded_data JSON
		);`,
		"estimatedvehiclejourney": `CREATE TABLE IF NOT EXISTS public.estimatedvehiclejourney (
			id BIGINT PRIMARY KEY DEFAULT nextval('public.estimatedvehiclejourney_id_seq'),
			servicedelivery INTEGER,
			recordedattime TIMESTAMP,
			lineref VARCHAR,
			directionref VARCHAR,
			datasource VARCHAR,
			datedvehiclejourneyref VARCHAR,
			vehiclemode VARCHAR,
			dataframeref VARCHAR,
			originref VARCHAR,
			destinationref VARCHAR,
			operatorref VARCHAR,
			vehicleref VARCHAR,
			cancellation VARCHAR,
			other JSON,
			firstservicedelivery INTEGER
		);`,
		"servicedelivery": `CREATE TABLE IF NOT EXISTS public.servicedelivery (
			id INTEGER PRIMARY KEY DEFAULT nextval('public.servicedelivery_id_seq'),
			responsetimestamp TIMESTAMPTZ,
			recordedattime TIMESTAMPTZ,
			data JSON
		);`,
	}

	for table, createStmt := range tables {
		var exists bool
		err := db.QueryRow(fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", table)).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}

		if !exists {
			_, err := db.Exec(createStmt)
			if err != nil {
				return fmt.Errorf("failed to create table %s: %w", table, err)
			}
			fmt.Printf("Table %s created successfully!\n", table)
		} else {
			fmt.Printf("Table %s already exists.\n", table)
		}
	}

	// Check if the unique constraint exists before adding it to calls table
	var constraintExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = 'unique_estimatedvehiclejourney_order'
		);
	`).Scan(&constraintExists)
	if err != nil {
		return fmt.Errorf("failed to check if unique constraint exists: %w", err)
	}

	if !constraintExists {
		_, err = db.Exec(`ALTER TABLE calls ADD CONSTRAINT unique_estimatedvehiclejourney_order UNIQUE (estimatedvehiclejourney, "order");`)
		if err != nil {
			return fmt.Errorf("failed to add unique constraint to calls table: %w", err)
		}
		fmt.Println("Unique constraint added to calls table.")
	} else {
		fmt.Println("Unique constraint already exists on calls table.")
	}

	// Check if the unique constraint exists before adding it to estimatedvehiclejourney table
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = 'unique_lineref_directionref_datasource_datedvehiclejourneyref'
		);
	`).Scan(&constraintExists)
	if err != nil {
		return fmt.Errorf("failed to check if unique constraint exists: %w", err)
	}

	if !constraintExists {
		_, err = db.Exec(`ALTER TABLE estimatedvehiclejourney ADD CONSTRAINT unique_lineref_directionref_datasource_datedvehiclejourneyref UNIQUE (lineref, directionref, datasource, datedvehiclejourneyref);`)
		if err != nil {
			return fmt.Errorf("failed to add unique constraint to estimatedvehiclejourney table: %w", err)
		}
		fmt.Println("Unique constraint added to estimatedvehiclejourney table.")
	} else {
		fmt.Println("Unique constraint already exists on estimatedvehiclejourney table.")
	}

	fmt.Println("Database setup is good!")
	return nil
}
