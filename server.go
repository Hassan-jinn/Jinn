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
// SecretKey is loaded from the Railway environment variable (more secure)
var SecretKey = os.Getenv("SECRET_KEY") 
const DataFile = "user_data.json"

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

// Global DB
var userDatabase = make(map[string]UserData)
var dbMutex = &sync.RWMutex{}

// --- JSON PERSISTENCE FUNCTIONS ---
func loadUserData() {
    // ... (JSON loading logic remains the same as before) ...
    // For Railway, we ensure the file exists or load an empty map.
    dbMutex.Lock()
    defer dbMutex.Unlock()

    data, err := os.ReadFile(DataFile)
    if err == nil {
        json.Unmarshal(data, &userDatabase)
    }
    fmt.Printf("[DB] Loaded %d records.\n", len(userDatabase))
}

func saveUserData() {
    // ... (JSON saving logic remains the same as before) ...
    dbMutex.RLock()
    defer dbMutex.RUnlock()

    data, err := json.MarshalIndent(userDatabase, "", "  ")
    if err == nil {
        os.WriteFile(DataFile, data, 0600)
    }
}

// --- SECURITY & HANDLER FUNCTIONS ---

func generateSignature(fullKey string, timestamp int64, clientID string) string {
    if SecretKey == "" {
        fmt.Println("[SECURITY ERROR] Secret Key not set!")
        return "ERROR"
    }
    raw := fmt.Sprintf("%s|%d|%s", fullKey, timestamp, clientID)
    mac := hmac.New(sha256.New, []byte(SecretKey))
    mac.Write([]byte(raw))
    return hex.EncodeToString(mac.Sum(nil))
}

func verifyKeyHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost || SecretKey == "" {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Server not configured or method wrong."})
        return
    }

    var req ApprovalRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ApprovalResponse{Status: "fail", Message: "Invalid request format"})
        return
    }

    // HMAC & TIMESTAMP Checks (Bypass and Replay Prevention)
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
            Status:      "fail",
            Message:     "Key is not active or has expired. Please activate.",
            PaymentData: paymentData,
            TotalUsers:  getTotalApprovedKeys(),
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
        Status:        "success",
        Message:       "Key Approved!",
        RemainingTime: remainingSeconds,
        TotalUsers:    getTotalApprovedKeys(),
    })
}

func getTotalApprovedKeys() int {
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

func main() {
    loadUserData()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/verify_key", verifyKeyHandler)
    fmt.Printf("Go Approval Server running on port %s...\n", port)
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Printf("Server failed: %v\n", err)
        os.Exit(1)
    }
}
