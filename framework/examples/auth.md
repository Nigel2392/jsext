# In this package you can find documentation on authentication with the JSExt package.

We will provide an example file, showing the login, registration and logout process.
```go
//go:build js && wasm
// +build js,wasm

package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/app"
	"github.com/Nigel2392/jsext/app/tokens"
	"github.com/Nigel2392/jsext/components/forms"
	"github.com/Nigel2392/jsext/components/forms/validators"
	"github.com/Nigel2392/jsext/components/loaders"
	"github.com/Nigel2392/jsext/elements"
	"github.com/Nigel2392/jsext/messages"
	"github.com/Nigel2392/jsext/router"
)
```
This is an example to show token authentication with jsext.  
Because this is  just an example, we do not provide a server to handle the authentication.
You can use this example to implement your own authentication system.
```go
var Application = app.App("#app") // Leave empty to append to body!

func init() {

	// Initialize a new token, with a 90 day refresh token, and a 15 minute access token.
	// When the token is gained (Logging in, or registering), we will automatically run a manager, which will automatically refresh the token.
	// This refresh period is set to 90% of the token's lifetime, so that we can refresh the token before it expires.
	tokens.AuthToken = tokens.NewToken(90*24*time.Hour, 15*time.Minute, "access", "refresh", "detail")

	// Define some basic callbacks, to give an impression of what's possible.
	tokens.AuthToken.OnUpdate(func(t *tokens.Token) {
		// We need to manually set the token cookie.
		// This is to provide the developer options on how to store the token.
		tokens.SetTokenCookie(t)
		fmt.Println("Token updated!")
	})
	tokens.AuthToken.OnInit(func(t *tokens.Token) {
		println("HIDING URLS! ONINIT")
		tokens.SetTokenCookie(t)
	})
	tokens.AuthToken.OnUpdateError(func(err error) {
		println("SHOWING URLS! ONUPDATEERROR")
		tokens.DeleteTokenCookie()
	})
	tokens.AuthToken.OnReset(func() {
		println("SHOWING URLS! ONRESET")
		tokens.DeleteTokenCookie()
	})

	// Set the authtoken inside of the application.
	// This means we can access it from anywhere, when we register url's using the application.
    // Otherwise, the application defines a `WrapURL` method, which willwrap the router callback functions to provide the application, andaccess to the application data.
	Application.Data.Set("AuthToken", tokens.AuthToken)
	// We define an onload function, which will be called when the application is loaded.
	// In here, we fetch the old token, if one exists.
	// If it does, we will run the token manager manually, which will automatically refresh the token.
	Application.OnLoad(func() {
		token, err := tokens.GetTokenCookie(Application.Data.Get("AuthToken").(*tokens.Token))
		if err != nil {
			println(err.Error())
		}
		if token != nil {
			Application.Data.Set("AuthToken", token)
			Application.Data.Get("AuthToken").(*tokens.Token).RunManager()
			println("HIDING URLS! ONLOAD")
		}
	})
}

// Define a login form struct, which will be used to validate the form.
type LoginForm struct {
	Email    string 
	Password string 
}

func (l *LoginForm) Validate() error {
	if validators.EmptyCheck(l.Email, l.Password) {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Email and password required")
	}
	if !strings.Contains(l.Email, "@") {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Invalid email")
	}
	if validators.IsValidLength(8, 32, l.Password) {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Invalid email or password")
	}
	return nil
}

// Define a register form struct, which will be used to validate the form.
type RegisterForm struct {
	Email            string 
	Phone_Number     string 
	Username         string 
	First_Name       string 
	Last_Name        string 
	Password         string 
	Password_Confirm string 
}

func (r *RegisterForm) Validate() error {
	if validators.EmptyCheck(r.Email, r.Username, r.Password, r.Password_Confirm) {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Email, username, password and password confirmation required")
	}
	if r.Password != r.Password_Confirm {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Passwords do not match")
	}
	if !strings.Contains(r.Email, "@") {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Invalid email")
	}
	if validators.IsValidLength(8, 32, r.Password) {
		//lint:ignore ST1005 error strings should not be capitalized
		return errors.New("Password must be between 8 and 32 characters.")
	}
	return nil
}

func LoginView(a *app.Application, v router.Vars, u *url.URL) {
	div := elements.Div().AttrClass("container")
	div.Append(elements.H1("Login"))
	var form = forms.NewForm("/", "POST")
	form.Inner.Label("Email", "email")
	form.Inner.Input("email", "email", "Email").AttrClass("login-form-input form-control")
	form.Inner.Label("Password", "password")
	form.Inner.Input("password", "password", "Password").AttrClass("login-form-input form-control")
	form.OnSubmit(func(data map[string]string, elements []jsext.Element) {
		var loginForm = &LoginForm{Email: data["email"], Password: data["password"]}
		if err := loginForm.Validate(); err != nil {
			// messages.ActiveMessages.NewError(err.Error(), 3000, 150)
			messages.SendError(err.Error())
			a.Render(div)
			return
		}
		// When using a.Load, the application will render the loader element, and then call the callback function.
		// This is useful for when you want to show a loader while the application is doing something, such as making a request.
		a.Load(func() {
			var err = a.Data.Get("AuthToken").(*tokens.Token).Login(map[string]string{
				"email":    data["email"],
				"password": data["password"],
			})
			if err != nil {
				messages.SendError(err.Error())
				// messages.ActiveMessages.NewError(err.Error(), 3000, 150)
				a.Render(div)
				return
			}
			a.Redirect("/")
			messages.SendSuccess("Successfully logged in.")
			// messages.ActiveMessages.NewSuccess("Successfully logged in.", 3000, 150)
		})
	})
	form.Inner.Button("Login").AttrType("submit").AttrClass("btn btn-primary")
	div.Append(form.Element())
	a.Render(div)
}

func RegisterView(a *app.Application, v router.Vars, u *url.URL) {
	div := elements.Div().AttrClass("container")
	div.Append(elements.H1("Register"))
	var form = forms.NewForm("/", "POST")
	form.Inner.Label("Email", "email")
	form.Inner.Input("email", "email", "Email").AttrClass("login-form-input form-control")
	form.Inner.Label("Phone Number", "phone_number")
	form.Inner.Input("phone_number", "tel", "Phone Number").AttrClass("login-form input form-control")
	form.Inner.Label("Username", "username")
	form.Inner.Input("username", "text", "Username").AttrClass("login-form-input form-control")
	form.Inner.Label("First Name", "first_name")
	form.Inner.Input("first_name", "text", "First Name").AttrClass("login-form-input form-control")
	form.Inner.Label("Last Name", "last_name")
	form.Inner.Input("last_name", "text", "Last Name").AttrClass("login-form-input form-control")
	form.Inner.Label("Password", "password")
	form.Inner.Input("password", "password", "Password").AttrClass("login-form-input form-control")
	form.Inner.Label("Password Confirmation", "password_confirm")
	form.Inner.Input("password_confirm", "password", "Password Confirmation").AttrClass("login-form-input form-control")
	form.Inner.Button("Register").AttrType("submit").AttrClass("btn btn-primary")
	form.OnSubmit(func(data map[string]string, elements []jsext.Element) {
		var registerForm = &RegisterForm{
			Email:            data["email"],
			Phone_Number:     data["phone_number"],
			Username:         data["username"],
			First_Name:       data["first_name"],
			Last_Name:        data["last_name"],
			Password:         data["password"],
			Password_Confirm: data["password_confirm"],
		}
		if err := registerForm.Validate(); err != nil {
			// messages.ActiveMessages.NewError(err.Error(), 3000, 150)
			messages.SendError(err.Error())
			a.Render(div)
			return
		}
		// When using a.Load, the application will render the loader element, and then call the callback function.
		a.Load(func() {
			var err = a.Data.Get("AuthToken").(*tokens.Token).Register(map[string]string{
				"email":            registerForm.Email,
				"phone_number":     registerForm.Phone_Number,
				"username":         registerForm.Username,
				"first_name":       registerForm.First_Name,
				"last_name":        registerForm.Last_Name,
				"password":         registerForm.Password,
				"password_confirm": registerForm.Password_Confirm,
			})
			if err != nil {
				messages.SendError(err.Error())
				a.Render(div)
				return
			}
			err = tokens.SetTokenCookie(a.Data.Get("AuthToken").(*tokens.Token))
			if err != nil {
				messages.SendError(err.Error())
				a.Render(div)
				return
			}
			a.Redirect("/")
			messages.SendSuccess("Successfully registered.")
		})
	})
	form.Inner.Button("Register").AttrType("submit").AttrClass("btn btn-primary")
	div.Append(form.Element())
	a.Render(div)
}

func LogoutView(a *app.Application, v router.Vars, u *url.URL) {
	// This is useful for when you want to show a loader while the application is doing something, such as making a request.
	a.Load(func() {
		var err = a.Data.Get("AuthToken").(*tokens.Token).Logout()
		if err != nil {
			messages.SendError(err.Error())
			a.Render(elements.Div())
			return
		}
		a.Redirect("/")
		messages.SendSuccess("Logged out")
	})
}

func index(a *app.Application, v router.Vars, u *url.URL) {
	a.Render(elements.Div().Append(elements.H1("Index")))
}

func main() {
	// Nagivate around this application using the URLs below.
	Application.SetLoader(loaders.NewLoader("#app", "Loading...", true, loaders.LoaderQuadSquares))

	Application.Register("index", "/", index)
	var auth = Application.Register("auth", "/auth/", nil)
	// URLs for the auth package will be prefixed with /auth/
	// Example: /auth/login
	auth.Register("login", "login/", Application.WrapURL(LoginView))
	auth.Register("register", "register/", Application.WrapURL(RegisterView))
	auth.Register("logout", "logout/", Application.WrapURL(LogoutView))

	os.Exit(Application.Run())
}
```