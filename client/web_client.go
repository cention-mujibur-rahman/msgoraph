package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	authURL = "https://login.microsoftonline.com/%v/oauth2/v2.0/authorize"
)

// Web is used to authenticate requests in the context of an online/user-facing app, such
// as a website. This type of client is mostly useful for debugging or for command line apps where
// the user configures their own app on the Microsoft Graph portal. In a normal web app, the
// InitializeCredentials()->setAuthorizationCode() part of this would be called on the client,
// then the code would be sent to the backend for the setAccessToken() part, given that that part
// does require an ApplicationSecret. Be sure to specify DelegatedOfflineAccess as a scope if you
// want refreshing to work.
type Web struct {
	TenantID           string
	ApplicationID      string
	ApplicationSecret  string
	AuthorizationCode  string
	Error              error
	LocalhostPort      int
	RedirectURI        string
	RefreshToken       string
	RequestCredentials *RequestCredentials
	Scopes             string
}

// NewWeb creates a new client.Web connection. To initialize the authentication on this, call
// web.InitializeCredentials()
func NewWeb(tenantID, applicationID, applicationSecret string, redirectURI string, scopes []string) *Web {
	scp := strings.Join(scopes, ",")
	return &Web{
		TenantID:           tenantID,
		ApplicationID:      applicationID,
		ApplicationSecret:  applicationSecret,
		RedirectURI:        redirectURI,
		RequestCredentials: &RequestCredentials{},
		Scopes:             scp,
	}
}

// Credentials returns back the set of request credentials in this client. Conforms to the
// client.Client interface.
func (w *Web) Credentials() *RequestCredentials {
	return w.RequestCredentials
}

// InitializeCredentials starts an oauth login flow to retrieve an authorization code, then exchange
// that authorization code for an access token and (if offline access is enabled) a refresh token.
func (w *Web) InitializeCredentials() error {
	err := w.setAccessToken()
	return err
}

// InitializeAuth starts an oauth login flow to retrieve an authorization code, then exchange
// that authorization code for an access token and (if offline access is enabled) a refresh token.
func (w *Web) Auth() string {
	return w.setAuthorizationCode()
}

func (w *Web) localServer() *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%v", w.LocalhostPort)}
	http.HandleFunc("/login", func(wr http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.Error = fmt.Errorf("Error while parsing form from response %s", err)
			return
		}
		if v, ok := r.Form["error"]; ok && len(v) > 0 {
			errorDescription, ok := r.Form["error_description"]
			if ok && len(errorDescription) > 0 {
				err = fmt.Errorf("%v: %v", strings.Join(v, ""), errorDescription)
				fmt.Fprintf(wr, "%v", err)
				w.Error = err
			} else {
				err = fmt.Errorf("%v", strings.Join(v, ""))
				fmt.Fprintf(wr, "%v", err)
				w.Error = err
			}
			return
		}
		code, codeOk := r.Form["code"]
		if len(code) > 0 && codeOk {
			fmt.Fprintf(wr, "authorization done. you may close this window now")
			w.AuthorizationCode = strings.Join(code, "")
			return
		}
		err = fmt.Errorf("error getting authorization code from login response")
		fmt.Fprintf(wr, "%v", err)
		w.Error = err
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// This will throw an error when we shutdown the server during the normal authorization flow
			// So we try to catch that error, and only return the real error if it isn't the expected
			// error.
			if !strings.Contains(err.Error(), "Server closed") {
				w.Error = fmt.Errorf("error on ListenAndServe: %v", err)
			}
		}
	}()
	return srv
}

//Authorization Request an authorization code
//The authorization code flow begins with the client directing the user to the /authorize endpoint.
//
//https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow
func (w *Web) Authorization() string {
	return fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/authorize?response_mode=query&state=12345&response_type=code&client_id=%v&scope=%v&redirect_uri=%v", w.TenantID, w.ApplicationID, strings.ReplaceAll(w.Scopes, ",", "%20"), w.RedirectURI)
}

func (w *Web) redirectURI() string {
	return fmt.Sprintf("http://localhost:%v/login", w.LocalhostPort)
}

// RefreshCredentials will attempt to refresh the access token if it is expired. This call will fail
// if the original authorization was not made with a Offline scope provided.
func (w *Web) RefreshCredentials() error {
	if w.RefreshToken == "" {
		return fmt.Errorf("client.Web: no refresh token found in web client. call client.InitializeCredentials to fill this")
	}
	w.RequestCredentials.AccessTokenUpdating.Lock()
	defer w.RequestCredentials.AccessTokenUpdating.Unlock()

	turl := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token", w.TenantID)
	tokenURI, err := url.Parse(turl)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(tokenURI.String(), url.Values{
		"client_id":     {w.ApplicationID},
		"client_secret": {w.ApplicationSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {w.RefreshToken},
	})
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	serverErrCode, ok := data["error"].(string)
	if ok {
		serverErr, ok := data["error_description"].(string)
		if ok {
			return fmt.Errorf("%v: %v", serverErrCode, serverErr)
		}
		return fmt.Errorf(serverErrCode)
	}
	accessToken, ok := data["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("no access token found in response")
	}
	durationSecs, ok := data["expires_in"].(float64)
	if !ok || durationSecs == 0 {
		return fmt.Errorf("no token duration found in response")
	}
	refreshToken, ok := data["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return fmt.Errorf("no refresh token found in response")
	}
	expiresAt := time.Now().Add(time.Duration(durationSecs) * time.Second)
	w.RequestCredentials.AccessToken = accessToken
	w.RequestCredentials.AccessTokenExpiresAt = expiresAt
	w.RefreshToken = refreshToken
	return nil
}

func (w *Web) setAccessToken() error {
	if w.AuthorizationCode == "" {
		return fmt.Errorf("client.Web: no access code found in web client")
	}
	w.RequestCredentials.AccessTokenUpdating.Lock()
	defer w.RequestCredentials.AccessTokenUpdating.Unlock()
	if w.RequestCredentials.AccessToken != "" && w.RequestCredentials.AccessTokenExpiresAt.After(time.Now()) {
		return nil
	}
	turl := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token", w.TenantID)
	tokenURI, err := url.Parse(turl)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(tokenURI.String(), url.Values{
		"client_id":     {w.ApplicationID},
		"client_secret": {w.ApplicationSecret},
		"code":          {w.AuthorizationCode},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {w.RedirectURI},
	})
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	serverErrCode, ok := data["error"].(string)
	if ok {
		serverErr, ok := data["error_description"].(string)
		if ok {
			return fmt.Errorf("%v: %v", serverErrCode, serverErr)
		}
		return fmt.Errorf(serverErrCode)
	}
	accessToken, ok := data["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("no access token found in response")
	}
	durationSecs, ok := data["expires_in"].(float64)
	if !ok || durationSecs == 0 {
		return fmt.Errorf("no token duration found in response")
	}
	//if w.Scopes.HasScope(scopes.DelegatedOfflineAccess) {
	refreshToken, ok := data["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return fmt.Errorf("no refresh token found in response")
	}
	w.RefreshToken = refreshToken
	//}
	expiresAt := time.Now().Add(time.Duration(durationSecs) * time.Second)
	w.RequestCredentials.AccessToken = accessToken
	w.RequestCredentials.AccessTokenExpiresAt = expiresAt
	return nil
}
func (w *Web) setAuthorizationCode() string {
	formVals := url.Values{}
	formVals.Set("client_id", w.ApplicationID)
	formVals.Set("grant_type", "authorization_code")
	formVals.Set("redirect_uri", w.RedirectURI)
	formVals.Set("response_mode", "query")
	formVals.Set("response_type", "code")
	formVals.Set("scope", w.Scopes)
	uri, err := url.Parse("https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	if err != nil {
		return "something wrong"
	}
	uri.RawQuery = formVals.Encode()

	return uri.String()
}
