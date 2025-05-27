package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"supernom/api"
)

func main() {
	log.Println("🚀 Starting SuperNom - Decentralized VPN on Binomena Blockchain")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get port from environment or use default
	port := os.Getenv("SUPERNOM_PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize SuperNom Gateway
	gateway := api.NewSuperNomGateway(port)

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 SuperNom Gateway starting on port %s", port)
		log.Printf("📊 Dashboard: http://localhost:%s/status", port)
		log.Printf("🔐 Auth API: http://localhost:%s/auth/check", port)
		log.Printf("💡 Health: http://localhost:%s/health", port)

		if err := gateway.Start(); err != nil {
			log.Fatalf("❌ Failed to start gateway: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("✅ SuperNom is running! Press Ctrl+C to shutdown")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Block until we receive a signal
	<-sigChan

	log.Println("🛑 Shutting down SuperNom gracefully...")
	log.Println("💫 Thanks for using SuperNom - The Future of Decentralized Internet!")
}
