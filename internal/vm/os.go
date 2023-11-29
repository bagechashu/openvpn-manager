package vm

import (
	"os"
	"strings"
)

func GetOsType() (osType string) {
	// Read the contents of the respective files to detect the OS
	fileContents, err := os.ReadFile("/etc/os-release")
	if err == nil && strings.Contains(string(fileContents), "ubuntu") {
		osType = "ubuntu"
	} else if _, err := os.ReadFile("/etc/debian_version"); err == nil {
		osType = "debian"
	} else if checkIfFileExists("/etc/almalinux-release") || checkIfFileExists("/etc/rocky-release") || checkIfFileExists("/etc/centos-release") {
		osType = "centos"
	} else if checkIfFileExists("/etc/fedora-release") {
		osType = "fedora"
	} else {
		osType = "unknown"
	}
	return osType
}
func GetNobodyGroup() (groupName string) {
	osType := GetOsType()
	switch osType {
	case "ubuntu":
		groupName = "nogroup"
	case "debian":
		groupName = "nogroup"
	case "centos":
		groupName = "nobody"
	case "fedora":
		groupName = "nobody"
	default:
		groupName = "nobody"
	}
	return groupName
}

// Function to check if a file exists
func checkIfFileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
