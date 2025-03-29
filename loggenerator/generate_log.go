package loggenerator

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func randomLogLevel() string {
	logLevels := []string{"INFO", "ERROR", "DEBUG", "WARN", "TRACE"}
	return logLevels[rand.Intn(len(logLevels))]
}

// Generate a random user ID
func randomUser() string {
	users := []string{"user123", "admin22", "guest001", "devOps", "root"}
	return users[rand.Intn(len(users))]
}

// Generate a random HTTP method
func randomMethod() string {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	return methods[rand.Intn(len(methods))]
}

// Generate a random HTTP status code
func randomStatus() int {
	statusCodes := []int{200, 301, 403, 404, 500, 422}
	return statusCodes[rand.Intn(len(statusCodes))]
}

func randomURL() string {
	domains := []string{"server1-us-central.com", "server2-us-central.com", "server3-us-central.com"}
	paths := []string{"/home", "/v1/shoppers", "/v1/candidates", "/api/data"}
	return fmt.Sprintf("https://%s%s", domains[rand.Intn(len(domains))], paths[rand.Intn(len(paths))])
}

// Generate a random application name
func randomApp() string {
	apps := []string{"shopper-onboarding", "gig-candidates", "gig-identity-verification"}
	return apps[rand.Intn(len(apps))]
}

// Generate a random datacenter location
func randomDatacenter() string {
	locations := []string{"us-east-1", "us-central-1", "us-west-1"}
	return locations[rand.Intn(len(locations))]
}

// Generate a random error event
func randomError() string {
	errors := []string{
		"CONNECTION_TIMEOUT", "RATE_LIMIT_EXCEEDED",
		"INTERNAL_SERVER_ERROR", "BAD_REQUEST", "KAFKA_ERROR",
	}
	return errors[rand.Intn(len(errors))]
}

// Generate a random metadata source
func randomSource() string {
	sources := []string{"proxy", "firewall", "vpn_gateway", "endpoint_security", "dns_filter"}
	return sources[rand.Intn(len(sources))]
}

// Generate a random log filename
func randomFilename() string {
	files := []string{"applog_2024_03_22.log", "applog_2024_03_23.log", "security_alerts.log"}
	return files[rand.Intn(len(files))]
}

func generateLogEntry() string {
	timestamp := time.Now().Format(time.RFC3339Nano)
	loglevel := randomLogLevel()
	method := randomMethod()
	url := randomURL()
	status := randomStatus()
	size := rand.Intn(5000) + 500 // Random response size (500-5000 bytes)
	app := randomApp()
	source := randomSource()
	filename := randomFilename()
	datacenter := randomDatacenter()
	user := randomUser()
	ip := randomIP()

	if rand.Float32() < 0.1 {
		errorMsg := randomError()
		return fmt.Sprintf(
			"%s ERROR %s %s %s | app=%s | source=%s | file=%s | datacenter=%s | ip=%s | user=%s\n",
			timestamp, method, url, errorMsg, app, source, filename, datacenter, ip, user)
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d bytes | app=%s | source=%s | file=%s | datacenter=%s, | ip=%s | user=%s\n",
		timestamp, loglevel, method, url, status, size, app, source, filename, datacenter, ip, user)
}

func GenerateLog(filename string) error {

	if _, err := os.Stat("/var/log"); os.IsNotExist(err) {
		fmt.Println("Creating /var/log directory...")
		os.Mkdir("/var/log", 0755)
	}
	tmpFile := fmt.Sprintf("/tmp/%s", filename)
	finalFile := fmt.Sprintf("/var/log/%s", filename)
	file, err := os.Create(tmpFile)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return err
	}
	defer file.Close()

	var totalSize int64 = 0
	targetSize := int64(1024 * 1024 * 1024) // 1 GB
	for totalSize < targetSize {
		logEntry := generateLogEntry()
		n, err := file.WriteString(logEntry)
		if err != nil {
			fmt.Println("Error writing to log file:", err)
			return err
		}
		totalSize += int64(n)
	}
	fmt.Println("1GB+ log file with error lines generated successfully.")

	err = os.Rename(tmpFile, finalFile)
	if err == nil {
		return nil
	}

	if os.IsPermission(err) {
		fmt.Println("permission denied to move file with current user: trying with elevated privilages")
		cmd := exec.Command("sudo", "mv", tmpFile, finalFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error moving file:", err)
			return err
		}

		cmd = exec.Command("chmod", "644", finalFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error setting permissions:", err)
			return err
		}
	}
	return nil
}
