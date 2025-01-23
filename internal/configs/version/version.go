package version

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	logo = `
___________.__              .__       .___
\__    ___/|__| ____   ____ |  |    __| _/
  |    |   |  |/ __ \ /    \|  |   / __ |
  |    |   |  \  ___/|   |  \  |__/ /_/ |
  |____|   |__|\___  >___|  /____/\____ |
				   \/     \/           \/
`
)

func GetInfoVersion() string {
	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	version := os.Getenv("APP_VERSION")
	concurrency := runtime.NumCPU()

	fmt.Printf("%s", logo)

	fmt.Printf("\033[32m🚀 Server Info:\033[0m\n")
	fmt.Printf("  🌐 \033[34mVersion:\033[0m %s\n", version)
	fmt.Printf("  🌍 \033[34mPort:\033[0m %s\n", port)
	fmt.Printf("  🛠️  \033[34mConcurrency:\033[0m %d CPUs\n", concurrency)
	fmt.Printf("  🔄 \033[34mFork:\033[0m Enabled\n")
	fmt.Println()
	log.Printf("\033[32m🚀 TienLD Framework App %s started at %s\n", version, time.Now().Format(time.RFC1123))
	log.Println("\033[33m🌐 Visit: http://localhost" + port + "\033[0m")
	fmt.Println("\033[32m🚀 Server is running on port 🚀", port)
	fmt.Println()

	return port
}
