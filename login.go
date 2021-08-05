package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// cookie handling
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionid"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, pass string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("sessionus", value); err == nil {
		cookie := &http.Cookie{
			Name:  "sessionus",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
	valuepw := map[string]string{
		"name": pass,
	}
	if encoded, err := cookieHandler.Encode("sessionpw", valuepw); err == nil {
		cookie := &http.Cookie{
			Name:  "sessionpw",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookieus := &http.Cookie{
		Name:   "sessionus",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookieus)

	cookiepw := &http.Cookie{
		Name:   "sessionpw",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookiepw)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	printLog(request)
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/main"
	if name != "" && pass != "" {
		cfgUser, cfgPass := getUserData()

		if cfgUser != "" && cfgPass != "" && cfgUser == name && cfgPass == pass {
			setSession(name, pass, response)
			redirectTarget = "/"
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	printLog(request)
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page
const indexPage = `
<!DOCTYPE html>
<html>
<head>
	<title>Login</title>
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<link rel="stylesheet" href="assets/fonts/ionicons/css/ionicons.css">
	<link rel="stylesheet" href="assets/css/bootstrap.css">
	<link rel="stylesheet" href="assets/css/style.css">
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
	fmt.Fprintf(response, indexPage)
}

func loginCheck(w http.ResponseWriter, r *http.Request) bool {

	var name string
	var pass string

	if cookie, err := r.Cookie("sessionus"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("sessionus", cookie.Value, &cookieValue); err == nil {
			name = cookieValue["name"]
		}
	}

	if cookie, err := r.Cookie("sessionpw"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("sessionpw", cookie.Value, &cookieValue); err == nil {
			pass = cookieValue["name"]
		}
	}

	cfgUser, cfgPass := getUserData()
	if name != cfgUser || pass != cfgPass {
		http.Redirect(w, r, "/main", 302)
		return false
	}
	return true
}
