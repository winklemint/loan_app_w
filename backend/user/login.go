package user

import (
	con "backend/Config"
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	//"image"
	"image/color"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

var jwtKey = []byte("secret key")
var alreadyExecuted bool

func init() {
	var err error
	_, err = con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
}

type Credential struct {
	//UserID     int    `json:"user_id"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password"`
	ContactNum int    `json:"contact_num"`
	Captcha    string `json:"captcha"`
	jwt.StandardClaims
}

type Session struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	_, err := con.GetDB()

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check request method
	if r.Method != http.MethodPost {
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credential Credential
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		//http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// fmt.Println(credential)
	user, err := GetUserFromDatabase(credential.Email)

	if err != nil || user == nil || !ComparePasswords(credential.Password, user.Password) || user.Contact_Num != credential.ContactNum {
		fmt.Println(credential.Password)
		fmt.Println(user.Password)

		logrus.Warnln("Invalid credentials for user:", credential.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Credential{
		Email: credential.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	session := &Session{
		UserID:    user.Id,
		Token:     tokenString,
		ExpiresAt: expirationTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = SaveSession(session)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	}

	w.WriteHeader(http.StatusOK)
	logrus.Infoln("Loggin Successfully", user.Email)

	go removeExpiredSessions()

}

func removeExpiredSessions() {
	for {
		// Check for expired sessions and delete them
		err := DeleteExpiredSessions()
		if err != nil {
			logrus.Errorln("Error deleting expired sessions:", err)
		}

		fmt.Println("grtn")

		// Sleep for a certain interval before checking again
		// Adjust the interval based on your application needs
		time.Sleep(30 * time.Minute)
	}
}

// DeleteExpiredSessions deletes expired sessions from the database
func DeleteExpiredSessions() error {
	// Connect to the database
	db, err := con.GetDB()
	if err != nil {
		logrus.Errorln("failed to connect to the database:")
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Retrieve the expired session tokens from the database
	rows, err := db.Query("SELECT token FROM sessions WHERE  DATE_ADD(created_at, INTERVAL 30 MINUTE) < NOW()")
	if err != nil {
		logrus.Errorln("failed to fetch expired sessions from the database:")
		return fmt.Errorf("failed to fetch expired sessions from the database: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and delete the corresponding sessions
	for rows.Next() {
		var sessionToken string
		err := rows.Scan(&sessionToken)
		if err != nil {
			logrus.Errorln("failed to scan session token from the database::")
			return fmt.Errorf("failed to scan session token from the database: %v", err)
		}

		// Delete the expired session from the database
		_, err = db.Exec("DELETE FROM sessions WHERE token = ?", sessionToken)
		if err != nil {
			logrus.Errorln("failed to delete expired session from the database")
			return fmt.Errorf("failed to delete expired session from the database: %v", err)
		} else {
			logrus.Errorln("Delete successful")
		}
	}

	return nil
}

func GetUserIDFromSession(session *Session) (int, error) {
	// Validate and parse the token
	token, err := jwt.ParseWithClaims(session.Token[6:], &Credential{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	// Extract the userID from the token claims
	_, ok := token.Claims.(*Credential)

	if !ok {
		return 0, errors.New("Failed to extract claims from token")
	}

	// Connect to the database
	db, err := con.GetDB()
	if err != nil {
		return 0, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Query the database for the user ID

	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", session.Token[6:]).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch user ID from the database: %v", err)
	}

	return userID, nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := con.GetDB()

	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &Credential{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*Credential)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userEmail := claims.Email

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome, %s!", userEmail)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	_, err := con.GetDB()
	cookie, err := r.Cookie("token")
	// w.Write([]byte(cookie.String()))
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &Credential{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*Credential)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 5*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	newClaims := &Credential{
		Email: claims.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   newTokenString,
		Expires: expirationTime,
	})

	w.WriteHeader(http.StatusOK)

}

func SaveSession(session *Session) error {

	db, err := con.GetDB()
	query := "INSERT INTO sessions (user_id, token, expires_at, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.UserID, session.Token, session.ExpiresAt)

	if err != nil {
		return fmt.Errorf("failed to save session in the database: %v", err)
	}

	return err
}

func GetUserIDFromDatabase(useremail string) (int, error) {
	db, err := con.GetDB()

	query := "SELECT id FROM user_info WHERE user_email = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare the query: %v", err)
	}
	defer stmt.Close()

	var userID int
	err = stmt.QueryRow(useremail).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to fetch user ID from the database: %v", err)
	}

	return userID, nil
}

// Modify the GetUserFromDatabase function to retrieve the hashed password from the database
func GetUserFromDatabase(email string) (*User_info, error) {
	db, err := con.GetDB()
	user := &User_info{}

	query := "SELECT id, user_email, user_contact_num, user_password FROM user_info WHERE user_email = ? LIMIT 1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&user.Id, &user.Email, &user.Contact_Num, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user from the database: %v", err)
	}
	fmt.Println(user.Password)

	return user, nil

}

func GetContactFromDatabase(contact int) (*User_info, error) {
	db, err := con.GetDB()
	user := &User_info{}

	query := "SELECT id, user_email, user_contact_num, user_password FROM user_info WHERE user_contact_num = ? LIMIT 1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(contact).Scan(&user.Id, &user.Email, &user.Contact_Num, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user from the database: %v", err)
	}
	//fmt.Println(user.Password)

	return user, nil
}

// Add the ComparePasswords function for comparing hashed passwords
func ComparePasswords(passwordHash, hashedPassword string) bool {
	hashBytes, err := hex.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, []byte(passwordHash))
	fmt.Println(hashBytes)
	fmt.Println([]byte(passwordHash))
	return err == nil
}

type OTPData struct {
	Email      string
	Contact    string
	OTP        string
	Expiration time.Time
}

var otpStore = make(map[string]OTPData)

type Captcha struct {
	Contact    int
	Captcha    string
	Expiration time.Time
}

var captchaStore = make(map[string]Captcha)

func generateOTP() string {

	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func storeOTP(email, otp string) {
	expiration := time.Now().Add(15 * time.Minute)
	otpStore[email] = OTPData{Email: email, OTP: otp, Expiration: expiration}
}
func StoreCaptcha(captcha string) {
	expiration := time.Now().Add(5 * time.Minute)
	captchaStore[captcha] = Captcha{Captcha: captcha, Expiration: expiration}
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var requestData struct {
			Email string `json:"email"`
		}
		var user User_info

		// Decode the request body into the requestData struct
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		db, err := con.GetDB()
		// Check if user exists in database
		rows, err := db.Query("SELECT user_email FROM user_info WHERE user_email = ?", requestData.Email)

		for rows.Next() {
			err = rows.Scan(&user.Email)

		}
		if requestData.Email != user.Email {
			http.Error(w, "No user found", http.StatusBadRequest)
		} else {

			email := requestData.Email
			//fmt.Println(email)

			otp := generateOTP()

			storeOTP(email, otp)

			err = sendOTPEmail(email, otp)
			if err != nil {
				http.Error(w, "Failed to send OTP via email", http.StatusInternalServerError)
				return

			}
			fmt.Fprint(w, "OTP sent successfully!")
		}

	}
}

func VerifyOTPhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		otp := r.FormValue("otp")

		otpData, ok := otpStore[email]
		fmt.Println(otpData)
		if !ok {
			http.Error(w, "Invalid Email or OTP", http.StatusUnauthorized)
			return
		}

		if otp == otpData.OTP && time.Now().Before(otpData.Expiration) {

			fmt.Fprint(w, "OTP verification successful!")

		} else {

			http.Error(w, "Invalid Email or OTP", http.StatusUnauthorized)

			return
		}
	}

}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		if err != nil {
			logrus.Error("Failed to load config:", err)
			// http.Error(w, "Failed to load config", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "POST" {
		email := r.FormValue("email")
		newPassword := r.FormValue("newpassword")
		conformpassword := r.FormValue("conformpassword")
		if newPassword != conformpassword {

			http.Error(w, "Invalid credentials", http.StatusUnauthorized)

		}
		fmt.Printf("Password reset successful for email: %s, new password: %s\n", email, newPassword)
		maxPasswordLength := 120 // Modify this value based on your database column's maximum length

		if len(newPassword) > maxPasswordLength {
			newPassword = newPassword[:maxPasswordLength]
		}

		// Encrypt the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedPasswordStr := hex.EncodeToString(hashedPassword)
		fmt.Println(hashedPasswordStr)
		_, err = db.Exec("UPDATE user_info SET user_password=? WHERE user_email=?", hashedPasswordStr, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete(otpStore, email)

		fmt.Fprint(w, "Password reset successful!")
	}
}

func sendOTPEmail(email, otp string) error {

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "raksha.rajput98@gmail.com"
	smtpPassword := "gmrzimywbfdtxslt"

	m := gomail.NewMessage()
	m.SetHeader("From", "raksha.rajput98@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "OTP for Password Reset")
	m.SetBody("text/plain", "Your OTP is: "+otp)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)

	}

	fmt.Println("Email sent successfully")

	return nil
}

// LogOut handles user logout by deleting the session from the database
func LogOut(w http.ResponseWriter, r *http.Request) {
	db, _ := con.GetDB()

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent, Authorization")

	// Retrieve the JWT token from the request header
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Unauthorized"}})
		return
	}

	// Extract the token from the Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	session := &Session{
		Token: tokenString,
	}

	// Use the GetUserIDFromSession function to retrieve the userID
	sessionToken, err := GetSessionFromSession1(session)

	if err != nil {
		logrus.Errorf("Error retrieving session token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		//json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Unauthorized"}})
		return
	}

	// Delete the session from the database
	_, err = db.Exec("DELETE FROM sessions WHERE token=?", sessionToken)
	if err != nil {
		logrus.Errorf("Error deleting session from the database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	} else {
		response := Response{
			Data: "Successfully logged out",
		}

		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
		// w.Write("jdjkded")

	}

}
func GetSessionFromSession1(session *Session) (string, error) {
	// Validate and parse the token

	// Validate and parse the token
	token, err := jwt.ParseWithClaims(session.Token, &Credential{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Extract the userID from the token claims
	_, ok := token.Claims.(*Credential)
	if !ok {
		return "", errors.New("Failed to extract claims from token")
	}
	db, err := con.GetDB()
	if err != nil {
		return "", fmt.Errorf("failed to connect to the database: %v", err)
	}

	var session1 string

	err = db.QueryRow("SELECT token FROM sessions WHERE token = ?", session.Token).Scan(&session1)

	if err != nil {

		fmt.Println(err)
		return "", fmt.Errorf("failed to fetch user ID from the database: %v", err)
	}

	return session1, nil
}

const (
	width      = 200
	height     = 80
	characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	charLen    = len(characters)
)

func generateCaptcha() (string, []byte, error) {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		code += string(characters[rand.Intn(charLen)])
	}

	dc := gg.NewContext(width, height)
	dc.SetRGB(67, 255, 100)
	dc.Clear()
	dc.SetRGB(255, 0, 0)
	dc.SetLineWidth(1)
	fontPath := "D:/hello/Mango_Sticky.ttf" // Replace with the path to your TrueType font file
	if err := dc.LoadFontFace(fontPath, 60); err != nil {
		fmt.Println("Failed to load font:", err)
		return "", nil, err
	}
	dc.DrawStringAnchored(code, float64(width)/2, float64(height)/2, 0.5, 0.5)

	// Distort the image a bit
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			noise := float64(rand.Intn(100)) / 130.0
			noise = noise * noise * 1.0
			px := float64(dc.Image().At(x, y).(color.RGBA).B) / 255.0
			px = px + noise
			dc.SetRGB(px, px, px)
			dc.SetPixel(x, y)
		}
	}

	// Encode the image as PNG
	var err error
	if err = dc.EncodePNG(ioutil.Discard); err != nil {
		fmt.Println("Failed to encode PNG:", err)
		return "", nil, err
	}

	// Encode the image as PNG and get the PNG image bytes
	var imageBuffer bytes.Buffer
	if err := dc.EncodePNG(&imageBuffer); err != nil {
		fmt.Println("Failed to encode PNG:", err)
		return "", nil, err
	}

	return code, imageBuffer.Bytes(), nil
}

func CaptchaHandler(w http.ResponseWriter, r *http.Request) {
	code, imageBytes, err := generateCaptcha()
	if err != nil {
		http.Error(w, "Failed to generate Captcha", http.StatusInternalServerError)
		return
	}

	StoreCaptcha(code)
	fmt.Println("Generated Captcha code:", code)

	w.Header().Set("Content-Type", "image/png")
	w.Write(imageBytes)
	return
}

func LoginWithDevice(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		logrus.Error("Failed to load config:", err)
		//http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent")

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credential Credential
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := GetContactFromDatabase(credential.ContactNum)
	if err != nil || user == nil || user.Contact_Num != credential.ContactNum {
		logrus.Warnln("Invalid credentials for user:", credential.ContactNum)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get the stored captcha data
	captchaData, ok := captchaStore[credential.Captcha]
	if !ok || time.Now().After(captchaData.Expiration) {
		fmt.Fprint(w, "Invalid captcha or expired captcha")
		return
	}

	// Validate the inputted captcha against the stored captcha
	if credential.Captcha != captchaData.Captcha {
		fmt.Fprint(w, "Invalid captcha")
		return
	}

	// At this point, the captcha is valid
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Credential{
		Email: credential.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	session := &Session{
		UserID:    user.Id,
		Token:     tokenString,
		ExpiresAt: expirationTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = SaveSession(session)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logrus.Infoln("Login Successfully", user.Email)

}
