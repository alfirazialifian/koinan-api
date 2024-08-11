package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

// Define the EmailRequest struct
type EmailRequest struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

// Declare environment variable variables
var (
	FROM      string
	TO        string
	PASSWORD  string
	SMTP_HOST string
	SMTP_PORT string
	APP_PORT  string
)

// Function to send an email
func sendEmail(subject, name, instance, message string) error {
	body := fmt.Sprintf("Name: %s\nInstance: %s\n\n%s", name, instance, message)

	msg := "From: " + FROM + "\n" +
		"To: " + TO + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", FROM, PASSWORD, SMTP_HOST)

	address := SMTP_HOST + ":" + SMTP_PORT
	err := smtp.SendMail(address, auth, FROM, []string{TO}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// Handler function for sending email
func emailHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var emailReq EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := sendEmail(emailReq.Subject, emailReq.Name, emailReq.Instance, emailReq.Message); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		log.Println("Failed to send email:", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email sent successfully")
}

// Handler function for the Hello World endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!")
}

// Initialize environment variables
func init() {
	FROM = os.Getenv("FROM_EMAIL")
	PASSWORD = os.Getenv("PASSWORD_EMAIL")
	TO = os.Getenv("TO_EMAIL")
	SMTP_HOST = os.Getenv("SMTP_HOST")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	APP_PORT = os.Getenv("APP_PORT")
}

// Main function to start the server
func main() {
	http.HandleFunc("/send-email", emailHandler)
	http.HandleFunc("/hello", helloHandler) // Add the new Hello World endpoint

	fmt.Printf("Server is running on port %s\n", APP_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", APP_PORT), nil))
}
