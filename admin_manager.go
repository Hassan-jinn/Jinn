package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

// --- CONFIGURATION ---
const AdminPassword = "H4ss4nKh4n4588" // !!! Final Admin Password !!!
const DataFile = "user_data.json"

// Struct for user data (MUST match server.go)
type UserData struct {
    Status        string
    ExpiryTime    time.Time
    TotalUsage    int
}

// Global variable for the in-memory database
var userDatabase = make(map[string]UserData)

// --- JSON PERSISTENCE FUNCTIONS (Local File) ---

func loadUserDataLocal() {
    data, err := os.ReadFile(DataFile)
    if err != nil {
        if os.IsNotExist(err) {
            return // File doesn't exist, start with empty map
        }
        fmt.Printf("\033[1;31m[DB Error] Reading data: %v\033[0m\n", err)
        return
    }
    if err := json.Unmarshal(data, &userDatabase); err != nil {
        fmt.Printf("\033[1;31m[DB Error] Unmarshalling data: %v\033[0m\n", err)
    }
}

func saveUserDataLocal() {
    data, err := json.MarshalIndent(userDatabase, "", "  ")
    if err != nil {
        fmt.Printf("\033[1;31m[DB Error] Marshalling data: %v\033[0m\n", err)
        return
    }
    if err := os.WriteFile(DataFile, data, 0600); err != nil {
        fmt.Printf("\033[1;31m[DB Error] Writing data: %v\033[0m\n", err)
    } else {
        fmt.Println("\033[1;32m[DB] Data saved successfully.\033[0m")
    }
}

// --- ADMIN FUNCTIONS (Same as before) ---

func loginAdmin() bool {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("\n\033[1;33mEnter Admin Password: \033[0m")
    password, _ := reader.ReadString('\n')
    password = strings.TrimSpace(password)

    if password == AdminPassword {
        fmt.Println("\033[1;32mAccess Granted! Welcome, JINN.\033[0m")
        return true
    }
    fmt.Println("\033[1;31mAccess Denied. Incorrect Password.\033[0m")
    return false
}

func approveKey(key string, days int) {
    if days <= 0 {
        fmt.Println("\033[1;31mError: Days must be a positive number.\033[0m")
        return
    }

    expiry := time.Now().Add(time.Hour * 24 * time.Duration(days))
    
    currentUsage := 0
    if data, exists := userDatabase[key]; exists {
        currentUsage = data.TotalUsage
    }

    userDatabase[key] = UserData{
        Status:     "Approved",
        ExpiryTime: expiry,
        TotalUsage: currentUsage,
    }

    saveUserDataLocal()

    fmt.Printf("\033[1;32m\nKey Approved Successfully!\033[0m\n")
    fmt.Printf("🔑 Key: \033[1;36m%s\033[0m\n", key)
    fmt.Printf("⏳ New Expiry: \033[1;36m%s\033[0m\n", expiry.Format("2006-01-02 15:04:05 MST"))
    fmt.Printf("📅 Duration: \033[1;36m%d Days\033[0m\n", days)
}

func showStats() {
    // ... (ShowStats function remains the same as before) ...
    totalApproved := 0
    totalExpired := 0
    
    fmt.Println("\n\033[1;37m--- CURRENT KEY STATS ---")
    
    for key, data := range userDatabase {
        remaining := data.ExpiryTime.Sub(time.Now())
        
        if data.Status == "Approved" && remaining > 0 {
            totalApproved++
            fmt.Printf("  ✅ \033[1;32mActive:\033[0m %s (\033[1;36m%s remaining\033[0m) | Usage: %d\n", key, remaining.Round(time.Hour).String(), data.TotalUsage)
        } else {
            totalExpired++
            fmt.Printf("  ❌ \033[1;31mExpired:\033[0m %s | Usage: %d\n", key, data.TotalUsage)
        }
    }
    
    fmt.Println("\n\033[1;33m--- SUMMARY ---")
    fmt.Printf("🔑 Total Active Keys: \033[1;32m%d\033[0m\n", totalApproved)
    fmt.Printf("💀 Total Expired Keys: \033[1;31m%d\033[0m\n", totalExpired)
    fmt.Println("\033[1;37m--------------------------")
}

func adminMenu() {
    // ... (Admin Menu remains the same) ...
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n\033[1;34m--- ADMIN PANEL ---")
		fmt.Println("\033[1;33m[1] Approve New Key")
		fmt.Println("[2] Show All Key Stats")
		fmt.Println("[0] Exit")
		fmt.Print("Select Option: \033[0m")
		
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "1" {
			fmt.Print("\033[1;32mEnter User Key (e.g., JINN-RAND-...): \033[0m")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
			
			fmt.Print("\033[1;32mEnter Days to Approve (e.g., 7, 30): \033[0m")
			daysStr, _ := reader.ReadString('\n')
			daysStr = strings.TrimSpace(daysStr)
			days, err := strconv.Atoi(daysStr)
			
			if err != nil {
				fmt.Println("\033[1;31mInvalid day input.\033[0m")
				continue
			}
			approveKey(key, days)
			
		} else if input == "2" {
			showStats()
		} else if input == "0" {
			fmt.Println("\033[1;36mExiting Admin Panel. Goodbye.\033[0m")
			break
		} else {
			fmt.Println("\033[1;31mInvalid Option.\033[0m")
		}
	}
}

func main() {
    loadUserDataLocal()

    fmt.Println("\n\033[1;96m===================================")
    fmt.Println("  JINN Key Management System (Admin)")
    fmt.Println("===================================\033[0m")
    
    if loginAdmin() {
        adminMenu()
    }
}
