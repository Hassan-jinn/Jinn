package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// --- CONFIGURATION ---
var SecretKey = os.Getenv("SECRET_KEY") 
const DataFile = "user_data.json" 

// (Baaqi Structs aur Global DB, loadUserData, saveUserData functions same rahengy)

// Structs
type UserData struct {
	Status        string 
	ExpiryTime    time.Time
	TotalUsage    int
}
type ApprovalRequest struct {
	ClientID  string `json:"client_id"`
	FullKey   string `json:"full_key"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}
type ApprovalResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	RemainingTime int64  `json:"remaining_time,omitempty"`
	TotalUsers    int    `json:"total_jinn_rand,omitempty"`
	PaymentData   string `json:"payment_data,omitempty"`
}
var userDatabase = make(map[string]UserData)
var dbMutex = &sync.RWMutex{}

// --- JSON PERSISTENCE (SAME AS BEFORE) ---
func loadUserData() {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	data, err := os.ReadFile(DataFile)
	if err == nil {
		json.Unmarshal(data, &userDatabase)
	}
    // Added safety check for debugging
    fmt.Printf("[DB] Loaded %d records on startup.\n", len(userDatabase))
}

func saveUserData() {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	data, err := json.MarshalIndent(userDatabase, "", "  ")
	if err == nil {
		os.WriteFile(DataFile, data, 0600)
	}
}

// (Baaqi security aur handler functions same rahengy)

func generateSignature(fullKey string, timestamp int64, clientID string) string {
    // ... (same as before) ...
	if SecretKey == "" { return "ERROR" }
	raw := fmt.Sprintf("%s|%d|%s", fullKey, timestamp, clientID)
	mac := hmac.New(sha256.New, []byte(SecretKey))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}

func getTotalApprovedKeys() int {
    // ... (same as before) ...
	totalApproved := 0
	dbMutex.RLock()
	for _, data := range userDatabase {
		if data.Status == "Approved" && data.ExpiryTime.After(time.Now()) {
			totalApproved++
		}
	}
	dbMutex.RUnlock()
	return totalApproved
}

func verifyKeyHandler(w http.ResponseWriter, r *http.Request) {
    // ... (same as before) ...
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost || SecretKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Server not configured or method wrong."})
		return
	}
    // ... (HMAC and Approval Logic is the same) ...

	var req ApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Invalid request format"})
		return
	}

	// HMAC & TIMESTAMP Checks
	expectedSignature := generateSignature(req.FullKey, req.Timestamp, req.ClientID)
	if req.Signature != expectedSignature || expectedSignature == "ERROR" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Signature mismatch. Request tampered."})
		return
	}
	if time.Now().Unix()-req.Timestamp > 60 || time.Now().Unix()-req.Timestamp < -5 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Timestamp too old."})
		return
	}

	// APPROVAL LOGIC
	dbMutex.RLock()
	userData, exists := userDatabase[req.FullKey]
	dbMutex.RUnlock()

	if !exists || userData.Status != "Approved" || userData.ExpiryTime.Before(time.Now()) {
		w.WriteHeader(http.StatusForbidden)
		paymentData := "PKR:150/7D-300/15D-550/30D|USD:1.5/7D-3.0/15D-5.0/30D"
		json.NewEncoder(w).Encode(ApprovalResponse{
			Status:      "fail", Message: "Key is not active or has expired. Please activate.",
			PaymentData: paymentData, TotalUsers: getTotalApprovedKeys(),
		})
		return
	}

	// Valid Key: Update usage count and save
	dbMutex.Lock()
	userData.TotalUsage++
	userDatabase[req.FullKey] = userData
	saveUserData() 
	dbMutex.Unlock()
	
	remainingSeconds := int64(userData.ExpiryTime.Sub(time.Now()).Seconds())
	json.NewEncoder(w).Encode(ApprovalResponse{
		Status:        "success", Message: "Key Approved!",
		RemainingTime: remainingSeconds, TotalUsers: getTotalApprovedKeys(),
	})
}


// --- MAIN FUNCTION (FIXED) ---

func main() {
	loadUserData() // Load data on startup

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
    
    // Yahan zaroori sudhaar kiya gaya hai: 0.0.0.0 par listen karna
	bindAddress := "0.0.0.0:" + port 

	http.HandleFunc("/verify_key", verifyKeyHandler)
	fmt.Printf("[INFO] Go Approval Server binding to %s...\n", bindAddress)
    
	if err := http.ListenAndServe(bindAddress, nil); err != nil {
		fmt.Printf("[FATAL] Server failed: %v\n", err)
		os.Exit(1)
	}
}
