package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Server struct {
	IP       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
	LogPath  string `json:"logPath"`
	Discord  string `json:"discord,omitempty"`
}

type Servers struct {
	Servers []Server `json:"servers"`
}

func Load(path string) (*Servers, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	var cfg Servers
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Servers) == 0 {
		return nil, errors.New("no servers defined in config")
	}
	return &cfg, nil
}

func save(path string, cfg *Servers) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func InitConfig() error {
	cfg, err := Load("config.json")
	if err != nil {
		return err
	}

	if cfg != nil && len(cfg.Servers) > 0 {
		return nil
	}

	if cfg == nil {
		cfg = &Servers{}
		log.Println("No config found, creating default config.json")
	}

	read := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Enter server IP:")
		ip, _ := read.ReadString('\n')
		ip = strings.TrimSpace(ip)

		fmt.Println("Enter server Port:")
		port, _ := read.ReadString('\n')
		port = strings.TrimSpace(port)

		fmt.Println("Enter RCON Password:")
		password, _ := read.ReadString('\n')
		password = strings.TrimSpace(password)

		fmt.Println("Enter log path:")
		logPath, _ := read.ReadString('\n')
		logPath = strings.TrimSpace(logPath)

		fmt.Println("Enter Discord invite (optional):")
		discordLink, _ := read.ReadString('\n')
		discordLink = strings.TrimSpace(discordLink)

		cfg.Servers = append(cfg.Servers, Server{
			IP:       ip,
			Port:     port,
			Password: password,
			LogPath:  logPath,
			Discord:  discordLink,
		})

		fmt.Print("Add another server? (y/n): ")
		another, _ := read.ReadString('\n')
		another = strings.TrimSpace(strings.ToLower(another))

		if another != "y" && another != "yes" {
			break
		}
	}

	if err := save("config.json", cfg); err != nil {
		log.Fatalf("Failed to save config.json: %v", err)
	}

	return nil
}
