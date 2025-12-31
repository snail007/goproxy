package services

import (
	"log"
	"net"
	"time"

	"github.com/snail007/goproxy/game"
	"github.com/snail007/goproxy/stealth"
)

// GameHTTP extends HTTP service with game optimizations and stealth features
type GameHTTP struct {
	HTTP
	gameProxy      *game.GameProxy
	headerRandom   *stealth.HeaderRandomizer
	trafficObfusc  *stealth.TrafficObfuscator
	stealthEnabled bool
	gameMode       bool
}

// NewGameHTTP creates a new game-optimized HTTP proxy
func NewGameHTTP() Service {
	gh := &GameHTTP{
		HTTP:           HTTP{},
		gameProxy:      game.NewGameProxy(),
		headerRandom:   stealth.NewHeaderRandomizer(),
		trafficObfusc:  stealth.NewTrafficObfuscator(true),
		stealthEnabled: true,
		gameMode:       true,
	}
	return gh
}

// InitService initializes the game HTTP service
func (gh *GameHTTP) InitService() {
	gh.HTTP.InitService()
	
	// Start game proxy cleanup
	gh.gameProxy.StartCleanup(30*time.Second, 5*time.Minute)
	
	log.Printf("Game HTTP proxy initialized (stealth: %v, game mode: %v)", 
		gh.stealthEnabled, gh.gameMode)
}

// Start starts the game HTTP proxy
func (gh *GameHTTP) Start(args interface{}) (err error) {
	// Use parent HTTP Start method
	err = gh.HTTP.Start(args)
	if err != nil {
		return err
	}

	log.Printf("Game HTTP proxy started on %s", *gh.cfg.Local)
	return nil
}

// callback handles incoming connections with stealth and game optimizations
func (gh *GameHTTP) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("game http conn handler crashed: %s", err)
		}
	}()

	// Optimize connection for low latency if game mode enabled
	if gh.gameMode {
		if tcpConn, ok := inConn.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)
			tcpConn.SetKeepAlive(true)
		}
	}

	// Use parent callback but with stealth modifications
	gh.HTTP.callback(inConn)
}

// EnableStealth enables stealth mode
func (gh *GameHTTP) EnableStealth() {
	gh.stealthEnabled = true
	log.Println("Stealth mode enabled")
}

// DisableStealth disables stealth mode
func (gh *GameHTTP) DisableStealth() {
	gh.stealthEnabled = false
	log.Println("Stealth mode disabled")
}

// EnableGameMode enables game mode optimizations
func (gh *GameHTTP) EnableGameMode() {
	gh.gameMode = true
	gh.gameProxy.Enable()
	log.Println("Game mode enabled")
}

// DisableGameMode disables game mode optimizations
func (gh *GameHTTP) DisableGameMode() {
	gh.gameMode = false
	gh.gameProxy.Disable()
	log.Println("Game mode disabled")
}

// GetStats returns proxy statistics
func (gh *GameHTTP) GetStats() map[string]interface{} {
	stats := gh.gameProxy.GetStats()
	stats["stealth_enabled"] = gh.stealthEnabled
	stats["game_mode"] = gh.gameMode
	return stats
}
