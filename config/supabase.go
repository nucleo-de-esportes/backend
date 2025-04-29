package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

func InitSupabase(env string) *supabase.Client {
	if env == "file" {
		err := godotenv.Load(filepath.Join("..", "supabase.env"))

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	supaUrl := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_KEY")

	return supabase.CreateClient(supaUrl, supaKey)

}
