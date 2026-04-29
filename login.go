package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

const sessionCookieName = "imagecloud_session"

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func setSession(userName string, response http.ResponseWriter, request *http.Request) {
	value := map[string]string{
		"name": userName,
	}

	encoded, err := cookieHandler.Encode(sessionCookieName, value)
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   request.TLS != nil,
	}
	http.SetCookie(response, cookie)
}

func clearSession(response http.ResponseWriter, request *http.Request) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   request.TLS != nil,
	}
	http.SetCookie(response, cookie)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	printLog(request)
	if request.Method != http.MethodPost {
		http.Redirect(response, request, "/main", http.StatusFound)
		return
	}

	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/main"
	if name != "" && pass != "" {
		cfgUser, cfgPass := getUserData()
		if cfgUser != "" && cfgPass != "" && cfgUser == name && cfgPass == pass {
			setSession(name, response, request)
			redirectTarget = "/"
		}
	}
	http.Redirect(response, request, redirectTarget, http.StatusFound)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	printLog(request)
	clearSession(response, request)
	http.Redirect(response, request, "/main", http.StatusFound)
}

const indexPage = `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <link rel="stylesheet" href="assets/fonts/ionicons/css/ionicons.css">
    <link rel="stylesheet" href="assets/css/bootstrap.css?v=1.0.0">
    <link rel="stylesheet" href="assets/css/style.css?v=1.0.1">
</head>
<body>
        <section class="cd-section index7 visible">
            <div class="cd-content style7">
                <div class="login-box">
                    <div class="login-form-slider">
                        <div class="login-slide slide">
                            <div class="login-header">
                                <div class="sign-up-txt text-right">
                                    ImageServer
                                </div>
                            </div>
                            <form class="padding-40px" method="post" action="/login">
                                <div class="form-group user-name-field">
                                    <input type="text" id="name" name="name" class="form-control" placeholder="User name">
                                    <div class="field-icon"><i class="ion-person"></i></div>
                                </div>
                                <div class="form-group margin-bottom-30px forgot-password-field">
                                    <input type="password" id="password" name="password" class="form-control" placeholder="Password">
                                    <div class="field-icon"><i class="ion-locked"></i></div>
                                </div>
                                <div class="form-group sign-in-btn">
                                    <input type="submit" class="submit" value="Login">
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </section>
</body>
</html>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	fmt.Fprint(response, indexPage)
}

func loginCheck(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/main", http.StatusFound)
		return false
	}

	cookieValue := make(map[string]string)
	if err := cookieHandler.Decode(sessionCookieName, cookie.Value, &cookieValue); err != nil {
		http.Redirect(w, r, "/main", http.StatusFound)
		return false
	}

	name := cookieValue["name"]
	cfgUser, cfgPass := getUserData()
	if name == "" || cfgUser == "" || cfgPass == "" || name != cfgUser {
		http.Redirect(w, r, "/main", http.StatusFound)
		return false
	}

	return true
}
