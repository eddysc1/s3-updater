package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	profiles := getAWSProfiles()
	if len(profiles) == 0 {
		fmt.Println("No AWS profiles found in ~/.aws/config")
		return
	}

	fmt.Println("Available AWS Profiles:")
	for i, profile := range profiles {
		fmt.Printf("%d: %s\n", i+1, profile)
	}

	selectedProfile := selectProfile(profiles)
	if selectedProfile == "" {
		fmt.Println("Invalid selection")
		return
	}

	authenticate(selectedProfile)
	buckets := listS3Buckets(selectedProfile)
	if len(buckets) == 0 {
		fmt.Println("No S3 buckets found")
		return
	}

	fmt.Println("Available S3 Buckets:")
	for i, bucket := range buckets {
		fmt.Printf("%d: %s\n", i+1, bucket)
	}

	selectedBucket := selectBucket(buckets)
	if selectedBucket == "" {
		fmt.Println("Invalid selection")
		return
	}

	downloadBucket(selectedProfile, selectedBucket)
	if askForReupload() {
		uploadBucket(selectedProfile, selectedBucket)
	}
}

func getAWSProfiles() []string {
	homeDir, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.aws/config", homeDir)
	file, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading AWS config file:", err)
		return nil
	}

	var profiles []string
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[profile ") {
			profile := strings.TrimPrefix(line, "[profile ")
			profile = strings.TrimSuffix(profile, "]")
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

func selectProfile(profiles []string) string {
	fmt.Print("Enter the number of the profile you want to use: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selection int
	_, err := fmt.Sscanf(input, "%d", &selection)
	if err != nil || selection < 1 || selection > len(profiles) {
		return ""
	}
	return profiles[selection-1]
}

func authenticate(profile string) {
	fmt.Print("Do you need to authenticate? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "y" {
		cmd := exec.Command("aws", "sso", "login", "--profile", profile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Authentication failed:", err)
			os.Exit(1)
		}
	}
}

func listS3Buckets(profile string) []string {
	cmd := exec.Command("aws", "s3api", "list-buckets", "--profile", profile, "--query", "Buckets[].Name", "--output", "text")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error listing S3 buckets:", err)
		return nil
	}
	buckets := strings.Fields(string(output))
	return buckets
}

func selectBucket(buckets []string) string {
	fmt.Print("Enter the number of the bucket you want to use: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selection int
	_, err := fmt.Sscanf(input, "%d", &selection)
	if err != nil || selection < 1 || selection > len(buckets) {
		return ""
	}
	return buckets[selection-1]
}

func downloadBucket(profile, bucket string) {
	cmd := exec.Command("aws", "s3", "sync", fmt.Sprintf("s3://%s", bucket), "/tmp/s3-uploader", "--profile", profile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error downloading bucket contents:", err)
		os.Exit(1)
	}
	fmt.Println("Download completed successfully. Modify contents in /tmp/s3-uploader and choose to reupload in the next step.")
}

func askForReupload() bool {
	fmt.Print("Do you want to reupload the files to the bucket? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return strings.ToLower(input) == "y"
}

func uploadBucket(profile, bucket string) {
	cmd := exec.Command("aws", "s3", "sync", "/tmp/s3-uploader", fmt.Sprintf("s3://%s", bucket), "--profile", profile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error uploading bucket contents:", err)
		os.Exit(1)
	}
	fmt.Println("Upload completed successfully.")
}
