package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Carrier  string `json:"carrier"`
}

func createConfig(configPath string) (*Config, error) {
	// If the file doesn't exist, create a default config
	defaultConfig := Config{}
	fmt.Println("欢迎使用南邮校园网自动登录程序, 请按照提示输入信息")
	fmt.Println("请输入校园网账号: ")
	_, err := fmt.Scanln(&defaultConfig.Username)
	if err != nil {
		print(err)
		os.Exit(1)
	}
	fmt.Println("请输入校园网密码: ")
	_, err = fmt.Scanln(&defaultConfig.Password)
	if err != nil {
		print(err)
		os.Exit(1)
	}
	fmt.Println("请选择运营商（输入数字）: 1. CMCC; 2. CHINANET; 3. NJUPT;(默认为 NJUPT)")
	var carrierIndex int
	_, err = fmt.Scanln(&carrierIndex)
	if err != nil {
		print(err)
		os.Exit(1)
	}
	switch carrierIndex {
	case 1:
		defaultConfig.Carrier = "cmcc"
	case 2:
		defaultConfig.Carrier = "njxy"
	default:
		defaultConfig.Carrier = ""
	}

	// Serialize the default config to JSON
	defaultConfigJSON, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return nil, err
	}

	// Create the directories leading to the config file
	_, err = os.Stat(filepath.Dir(configPath))
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(configPath), 0700)
		if err != nil {
			return nil, err
		}
	}

	// Write the default config to the file
	err = os.WriteFile(configPath, defaultConfigJSON, 0600)
	if err != nil {
		return nil, err
	}

	return &defaultConfig, nil
}

func readConfig() (*Config, error) {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Define the path to the config file
	configPath := filepath.Join(currentUser.HomeDir, ".config", "njupt-connect", "config.json")

	// Check if the config file exists
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		config, err := createConfig(configPath)
		if err != nil {
			return nil, err
		}
		return config, nil
	} else if err != nil {
		// If there was an error other than "file does not exist," return the error
		return nil, err
	}

	// If the file exists, read and parse the config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getIP() (string, error) {
	// Get the current user's IPs
	adds, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", errors.New("system IP not found")
	}
	for _, address := range adds {
		if aspnet, ok := address.(*net.IPNet); !(!ok || aspnet.IP.IsLoopback()) {
			// Get the first IPv4 starting with 10
			if aspnet.IP.To4() != nil && aspnet.IP.To4()[0] == 10 {
				return aspnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("school network IP not found")
}

func tryConnect(ip string, config *Config) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://p.njupt.edu.cn:802/eportal/portal/login?callback=dr1003&login_method=1&user_account=%2C0%2C"+config.Username+"%40"+config.Carrier+"&user_password="+config.Password+"&wlan_user_ip="+ip+"&wlan_user_ipv6=&wlan_user_mac=000000000000&wlan_ac_ip=&wlan_ac_name=&jsVersion=4.1.3&terminal_type=1&lang=zh-cn&v=2383&lang=zh",
		nil,
	)
	if err != nil {
		println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println(err)
			return
		}
	}(resp.Body)
}

func main() {
	ip, err := getIP()
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := readConfig()
	if err != nil {
		println(err)
		return
	}
	tryConnect(ip, config)
}
