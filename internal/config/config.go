package config

import (
	"bufio"
	"os"
	"strings"
)

// Config encapsulates runtime configuration for the Hirelly backend.
type Config struct {
	Port         string
	Env          string
	PostgresURL  string
	SupabaseURL  string
	SupabaseKey  string
	BrevoAPIKey  string
	AdminEmail   string
	SenderEmail  string
	SenderName   string
	CORSOrigins  []string
}

// Load reads configuration from environment variables, with fallback to local .env if present.
func Load() *Config {
	loadDotEnv(".env")

	port := getEnv("PORT", "8080")
	env := getEnv("ENV", getEnv("VERCEL_ENV", "development"))

	postgresURL := getEnv("POSTGRES_URL", 
		getEnv("POSTGRES_PRISMA_URL", 
		getEnv("DATABASE_URL", "")))

	supabaseURL := getEnv("SUPABASE_URL", 
		getEnv("NEXT_PUBLIC_SUPABASE_URL", 
		getEnv("STORAGE_SUPABASE_URL", "")))

	supabaseKey := getEnv("SUPABASE_SERVICE_ROLE_KEY", 
		getEnv("SUPABASE_SECRET_KEY", 
		getEnv("SUPABASE_ANON_KEY", "")))

	brevoAPIKey := getEnv("BREVO_API_KEY", 
		getEnv("BREVO_KEY", 
		getEnv("SENDINBLUE_API_KEY", "")))

	adminEmail := getEnv("ADMIN_EMAIL", "connect@hirelly.in")
	senderEmail := getEnv("SENDER_EMAIL", "connect@hirelly.in")
	senderName := getEnv("SENDER_NAME", "Hirelly Executive Advisory")

	corsRaw := getEnv("CORS_ORIGINS", "*")
	corsOrigins := strings.Split(corsRaw, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	return &Config{
		Port:        port,
		Env:         env,
		PostgresURL: postgresURL,
		SupabaseURL: strings.TrimRight(supabaseURL, "/"),
		SupabaseKey: supabaseKey,
		BrevoAPIKey: brevoAPIKey,
		AdminEmail:  adminEmail,
		SenderEmail: senderEmail,
		SenderName:  senderName,
		CORSOrigins: corsOrigins,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}
