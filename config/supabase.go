package config

import (
	"github.com/nedpals/supabase-go"
)

func InitSupabase() *supabase.Client {

	supaUrl := "https://txuubteajinfbretuxob.supabase.co"
	supaKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InR4dXVidGVhamluZmJyZXR1eG9iIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDE3MDAxNjMsImV4cCI6MjA1NzI3NjE2M30.Nff7OKkEEcISpcDaz6zA10aAd2Etw7YrFOFuC4EYosQ"

	return supabase.CreateClient(supaUrl, supaKey)

}
