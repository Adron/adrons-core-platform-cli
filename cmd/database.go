package cmd

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Display database connection information",
	Long: `Display detailed information about the database connection including
database name, connection properties, and other relevant details.`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresURL := viper.GetString("POSTGRES_URL")
		if postgresURL == "" {
			fmt.Println("Error: POSTGRES_URL is not set in configuration")
			return
		}

		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			fmt.Printf("Error parsing connection details: %v\n", err)
			return
		}
		defer db.Close()

		fmt.Println("Database Connection Information:")
		fmt.Println("-------------------------------")

		// Get database name
		var dbName string
		err = db.QueryRow("SELECT current_database()").Scan(&dbName)
		if err == nil {
			fmt.Printf("Database Name: %s\n", dbName)
		}

		// Get database version
		var version string
		err = db.QueryRow("SELECT version()").Scan(&version)
		if err == nil {
			fmt.Printf("Database Version: %s\n", version)
		}

		// Get current user
		var user string
		err = db.QueryRow("SELECT current_user").Scan(&user)
		if err == nil {
			fmt.Printf("Connected User: %s\n", user)
		}

		// Get server encoding
		var encoding string
		err = db.QueryRow("SHOW server_encoding").Scan(&encoding)
		if err == nil {
			fmt.Printf("Server Encoding: %s\n", encoding)
		}

		// Get timezone
		var timezone string
		err = db.QueryRow("SHOW timezone").Scan(&timezone)
		if err == nil {
			fmt.Printf("Timezone: %s\n", timezone)
		}

		// Parse and display connection string parts (safely)
		fmt.Println("\nConnection String Properties:")
		fmt.Println("----------------------------")
		fmt.Printf("SSL Mode: %s\n", getSSLMode(postgresURL))
		fmt.Printf("Host: %s\n", getHost(postgresURL))
		fmt.Printf("Port: %s\n", getPort(postgresURL))
	},
}

// dbTablesCmd represents the tables subcommand
var dbTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "List all tables in the database",
	Long: `Display a list of all tables in the connected database.
This includes both public and system tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresURL := viper.GetString("POSTGRES_URL")
		if postgresURL == "" {
			fmt.Println("Error: POSTGRES_URL is not set in configuration")
			return
		}

		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT 
				table_schema,
				table_name,
				(SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
			FROM information_schema.tables t
			WHERE table_schema = 'public'
			ORDER BY table_schema, table_name;
		`)
		if err != nil {
			fmt.Printf("Error querying tables: %v\n", err)
			return
		}
		defer rows.Close()

		fmt.Println("\nDatabase Tables:")
		fmt.Println("----------------")
		fmt.Printf("%-20s %-30s %s\n", "SCHEMA", "TABLE NAME", "COLUMNS")
		fmt.Println(strings.Repeat("-", 60))

		var count int
		for rows.Next() {
			var schema, name string
			var columnCount int
			if err := rows.Scan(&schema, &name, &columnCount); err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			fmt.Printf("%-20s %-30s %d\n", schema, name, columnCount)
			count++
		}

		if count == 0 {
			fmt.Println("No tables found in the public schema.")
		} else {
			fmt.Printf("\nTotal tables found: %d\n", count)
		}
	},
}

// tenantsCmd represents the tenants command
var tenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "List all tenants in the database",
	Long:  `Display a list of all tenants stored in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresURL := viper.GetString("POSTGRES_URL")
		if postgresURL == "" {
			fmt.Println("Error: POSTGRES_URL is not set in configuration")
			return
		}

		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT id, name, created_at
			FROM tenants
			ORDER BY name;
		`)
		if err != nil {
			fmt.Printf("Error querying tenants: %v\n", err)
			return
		}
		defer rows.Close()

		fmt.Println("\nTenants:")
		fmt.Println("---------")
		fmt.Printf("%-36s %-30s %-25s\n", "ID", "NAME", "CREATED AT")
		fmt.Println(strings.Repeat("-", 91))

		var count int
		for rows.Next() {
			var id, name string
			var createdAt time.Time
			if err := rows.Scan(&id, &name, &createdAt); err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			fmt.Printf("%-36s %-30s %-25s\n", id, name, createdAt.Format("2006-01-02 15:04:05"))
			count++
		}

		if count == 0 {
			fmt.Println("No tenants found.")
		} else {
			fmt.Printf("\nTotal tenants: %d\n", count)
		}
	},
}

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List all roles in the database",
	Long:  `Display a list of all roles stored in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresURL := viper.GetString("POSTGRES_URL")
		if postgresURL == "" {
			fmt.Println("Error: POSTGRES_URL is not set in configuration")
			return
		}

		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT id, name, description, created_at
			FROM roles
			ORDER BY name;
		`)
		if err != nil {
			fmt.Printf("Error querying roles: %v\n", err)
			return
		}
		defer rows.Close()

		fmt.Println("\nRoles:")
		fmt.Println("------")
		fmt.Printf("%-36s %-20s %-30s %-25s\n", "ID", "NAME", "DESCRIPTION", "CREATED AT")
		fmt.Println(strings.Repeat("-", 111))

		var count int
		for rows.Next() {
			var id, name, description string
			var createdAt time.Time
			if err := rows.Scan(&id, &name, &description, &createdAt); err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			fmt.Printf("%-36s %-20s %-30s %-25s\n", id, name, description, createdAt.Format("2006-01-02 15:04:05"))
			count++
		}

		if count == 0 {
			fmt.Println("No roles found.")
		} else {
			fmt.Printf("\nTotal roles: %d\n", count)
		}
	},
}

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users in the database",
	Long:  `Display a list of all users stored in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		postgresURL := viper.GetString("POSTGRES_URL")
		if postgresURL == "" {
			fmt.Println("Error: POSTGRES_URL is not set in configuration")
			return
		}

		db, err := sql.Open("postgres", postgresURL)
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT id, username, email, created_at
			FROM users
			ORDER BY username;
		`)
		if err != nil {
			fmt.Printf("Error querying users: %v\n", err)
			return
		}
		defer rows.Close()

		fmt.Println("\nUsers:")
		fmt.Println("------")
		fmt.Printf("%-36s %-20s %-30s %-25s\n", "ID", "USERNAME", "EMAIL", "CREATED AT")
		fmt.Println(strings.Repeat("-", 111))

		var count int
		for rows.Next() {
			var id, username, email string
			var createdAt time.Time
			if err := rows.Scan(&id, &username, &email, &createdAt); err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			fmt.Printf("%-36s %-20s %-30s %-25s\n", id, username, email, createdAt.Format("2006-01-02 15:04:05"))
			count++
		}

		if count == 0 {
			fmt.Println("No users found.")
		} else {
			fmt.Printf("\nTotal users: %d\n", count)
		}
	},
}

// Helper functions for parsing connection strings
func getSSLMode(connStr string) string {
	if strings.Contains(connStr, "sslmode=") {
		parts := strings.Split(connStr, "sslmode=")
		if len(parts) > 1 {
			return strings.Split(parts[1], "&")[0]
		}
	}
	return "not specified"
}

func getHost(connStr string) string {
	if strings.Contains(connStr, "@") {
		parts := strings.Split(connStr, "@")
		if len(parts) > 1 {
			hostPort := strings.Split(parts[1], "/")[0]
			return strings.Split(hostPort, ":")[0]
		}
	}
	return "not specified"
}

func getPort(connStr string) string {
	if strings.Contains(connStr, ":") {
		parts := strings.Split(connStr, ":")
		if len(parts) > 2 {
			return strings.Split(parts[2], "/")[0]
		}
	}
	return "5432 (default)"
}

func init() {
	// Add all database-related commands to root
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(tenantsCmd)
	rootCmd.AddCommand(rolesCmd)
	rootCmd.AddCommand(usersCmd)

	// Add tables as a subcommand of db
	dbCmd.AddCommand(dbTablesCmd)
}
