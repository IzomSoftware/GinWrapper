// package responses

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/IzomSoftware/GinWrapper/internal/configuration"
// 	"github.com/IzomSoftware/GinWrapper/internal/logger"
// 	"github.com/IzomSoftware/GinWrapper/internal/storage/redis_source"
// 	"github.com/IzomSoftware/GinWrapper/internal/storage/sql_source"
// 	"github.com/IzomSoftware/GinWrapper/internal/utils/hash_util"
// 	"github.com/IzomSoftware/GinWrapper/internal/utils/jwt_util"
// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/crypto/bcrypt"
// )

// type RateLimitProtection struct {
// 	Enabled bool
// 	Rate    int64
// 	Time    int64
// }

// /*
//  * the Protections struct, containing protections we provide on requests
//  */
// type Protections struct {
// 	/*
// 	 * Contains protections as basic as path check.
// 	 * Example: if path is /api/auth or /api/ (and is it malicious?)
// 	 */
// 	BasicProtections bool
// 	/*
// 	 * This protection only and only validates the UserAgent of the client
// 	 */
// 	UserAgent bool
// 	/*
// 	 * Wether to protect with JWTAPI or not. note that JWTAPI
// 	 * Only activates once one request uses it
// 	 */
// 	JWT bool
// 	/*
// 	 * Rate limit protection. protects the server from spam
// 	 * Attacks & API abuses
// 	 */
// 	RateLimit RateLimitProtection
// 	/*
// 	 * Wether to ban suspicious connections or not. this is
// 	 * Not a great idea if you enable protections with false-positives, obviously
// 	 */
// 	Ban bool
// }

// /*
//  * the Response struct. this struct contains everything a request would need.
//  */
// type Response struct {
// 	Handler     gin.HandlerFunc
// 	Type        string
// 	Addresses   []string
// 	Protections Protections
// }

// var (
// 	/*
// 	 * Responses map.
// 	 */
// 	Responses    = map[string]*Response{}
// 	NoRouteRoute = func(c *gin.Context) {
// 		/*
// 		 * We provide this variable for developers to set & setup custom screens (or do whatever they want)
// 		 */
// 		c.String(http.StatusNotFound, "404 Not Found")
// 	}
// 	InternalServerErrorRouteRoute = func(c *gin.Context) {
// 		/*
// 		 * We provide this variable for developers to set & setup custom screens (or do whatever they want)
// 		 */
// 	}
// 	UnexpectedTypeError = fmt.Errorf("Unexpected type for value")
// )

// /*
//  * This function executes before developer's 500 screen appear
//  * We do this to keep things modular but working
//  */
// func InternalServerErrorRoute(c *gin.Context, err error) {
// 	ip := c.ClientIP()

// 	InternalServerErrorRouteRoute(c)

// 	logger.LogError(fmt.Sprintf("%s", err))

// 	AbortConnection(ip, c, http.StatusInternalServerError)
// }

// /*
//  * This function executes before developer's 404 screen appear
//  * We do this to keep things modular but working
//  */
// func NoRoute(c *gin.Context) {
// 	if configuration.ConfigHolder.Protections.BasicProtections.Provide {
// 		path := c.Request.URL.Path
// 		// Split the path with "/" and then move to the next string (which is the first path)
// 		basePath := strings.Split(path, "/")[1]

// 		for _, response := range Responses {
// 			for _, address := range response.Addresses {

// 				// Path is incomplete. possible attack
// 				if strings.Contains(address, basePath) {
// 					ip := c.ClientIP()

// 					AbortSuspiciousConnection(ip, c)

// 					// Ban the connection if needed
// 					if configuration.ConfigHolder.Protections.BasicProtections.Aggressive {
// 						BanConnection(ip, c)
// 					}

// 					return
// 				}
// 			}
// 		}
// 	}

// 	// Finally execute the base function which developer's intend
// 	NoRouteRoute(c)
// }

// /*
//  * Adds all requests related to UserPassAPI so that we provide this API.
//  */
// func ActivateUserPassAPI() {
// 	// The Register API. Gets username & password & hashes the password & registers the user
// 	Responses["UserPassAPIRegister"] = &Response{
// 		Handler: func(c *gin.Context) {
// 			ip, username, password := c.ClientIP(), c.Query("username"), c.Query("password")

// 			err := sql_source.CreateUser(username, password)
// 			// Possible internal server error
// 			if err != nil {
// 				if err == sql_source.UserAlreadyExists {
// 					c.String(http.StatusBadRequest, "User already exists")
// 					AbortConnection(ip, c, http.StatusBadRequest)
// 					return
// 				}
// 				logger.LogError(fmt.Sprintf("UserPassAPI registration error: %s\n", err))
// 				AbortConnection(ip, c, http.StatusInternalServerError)
// 				return
// 			}

// 			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?username=%s&password=%s", "/api/auth/login", username, password))
// 		},
// 		Type:      "GET",
// 		Addresses: []string{"/api/auth/register"},
// 		Protections: Protections{
// 			BasicProtections: true,
// 			RateLimit: RateLimitProtection{
// 				Enabled: true,
// 				Rate:    5,
// 				Time:    3600,
// 			},
// 		},
// 	}
// 	// The Login API. checks if password is valid
// 	// Also is the generation API. generates the token after authorization is complete.
// 	// Requires no further actions from developers because we already provide
// 	// The UserPass boolean so we protect this request.
// 	Responses["UserPassAPILogin"] = &Response{
// 		Handler: func(c *gin.Context) {
// 			ip, username, password := c.ClientIP(), c.Query("username"), c.Query("password")

// 			result, err := sql_source.GetData("SELECT hash FROM Users WHERE username = ?", username)
// 			if err != nil {
// 				InternalServerErrorRoute(c, err)
// 				return
// 			}

// 			val := result.(string)

// 			// Checks the password (with all those bcrypt salting shit)
// 			err = hash_util.IsPasswordValid(val, password)
// 			// Error is not related to password validation so it means
// 			// we're dealing with something else.
// 			if err != nil {
// 				if err == bcrypt.ErrMismatchedHashAndPassword {
// 					c.String(http.StatusBadRequest, "Wrong Password")
// 					AbortConnection(ip, c, http.StatusBadRequest)
// 					return
// 				}

// 				InternalServerErrorRoute(c, err)
// 				return
// 			}

// 			token, err := jwt_util.GenerateJWT(username, "Bearer", time.Second*1)
// 			if err != nil {
// 				InternalServerErrorRoute(c, err)
// 				return
// 			}

// 			c.String(http.StatusOK, "%s", token)
// 		},
// 		Type:      "GET",
// 		Addresses: []string{"/api/auth/login"},
// 		Protections: Protections{
// 			BasicProtections: true,
// 			RateLimit: RateLimitProtection{
// 				Enabled: true,
// 				Rate:    5,
// 				Time:    60,
// 			},
// 		},
// 	}
// }

// /*
//  * Adds all requests related to JWTAPI so that we provide this API.
//  */
// func ActivateJWTAPI() {
// 	// The validation API. validates the token given, so we'll authorize.
// 	Responses["JWTAPIValidate"] = &Response{
// 		Handler: func(c *gin.Context) {
// 			// ip, token := c.ClientIP(), c.Query("token")

// 			// valid, err := jwt_util.ValidateJWT(token)
// 			// if err != nil {
// 			// 	AbortConnection(ip, c, http.StatusInternalServerError)
// 			// }

// 			// valid.

// 			// if valid {
// 			// 	c.String(http.StatusOK, "true")
// 			// } else {
// 			// 	c.String(http.StatusForbidden, "false")
// 			// }
// 			c.String(http.StatusOK, "")
// 		},
// 		Type:      "GET",
// 		Addresses: []string{"/api/auth/validate_token"},
// 		Protections: Protections{
// 			BasicProtections: true,
// 			RateLimit: RateLimitProtection{
// 				Enabled: true,
// 				Rate:    30,
// 				Time:    60,
// 			},
// 		},
// 	}
// }

// func (R *Response) OnProtectionFailure(c *gin.Context) {
// 	ip, protections := c.ClientIP(), R.Protections

// 	// Aborts the suspicious connection. (does not ban)
// 	AbortSuspiciousConnection(ip, c)
// 	// Bans the connection if ban required.
// 	if protections.Ban {
// 		BanConnection(ip, c)
// 	}
// }

// /*
//  * We call this function inside middleware(), where the connection begins.
//  * Does every single thing related to protections & responses
//  */
// func (R *Response) OnProtected(c *gin.Context) {
// 	response := R

// 	if !response.IsAnyProtectionEnabled() {
// 		return
// 	}

// 	ip, protections, header := c.ClientIP(), response.Protections, c.Request.Header

// 	// User agent check
// 	if userAgent, apiUserAgent := header.Get("User-Agent"),
// 		configuration.ConfigHolder.Protections.APIUserAgent;

// 	// The actual check, if user agent is valid, perform no further action
// 	protections.UserAgent && userAgent != apiUserAgent {

// 		response.OnProtectionFailure(c)
// 	}

// 	if protections.RateLimit.Enabled {
// 		lastRate, err := redis_source.GetLastRateLimit(c, ip)
// 		if err != nil {
// 			InternalServerErrorRoute(c, err)
// 			return
// 		}

// 		if time.Now().UnixMilli()-lastRate <= protections.RateLimit.Time*1000 {
// 			if err := redis_source.IncrementRateLimit(c, ip); err != nil {
// 				InternalServerErrorRoute(c, err)
// 				return
// 			}
// 		} else {
// 			if err = redis_source.UpdateHashValue(c, ip, "Rate", 1); err != nil {
// 				InternalServerErrorRoute(c, err)
// 				return
// 			}

// 			if err = redis_source.UpdateHashValue(c, ip, "LastRate", time.Now().UnixMilli()); err != nil {
// 				InternalServerErrorRoute(c, err)
// 				return
// 			}
// 		}

// 		rate, err := redis_source.GetRateLimit(c, ip)
// 		if err != nil {
// 			InternalServerErrorRoute(c, err)
// 			return
// 		}

// 		if rate > protections.RateLimit.Rate {
// 			response.OnProtectionFailure(c)
// 		}
// 	}
// }

// /*
//  * Returns true if any protection is enabled
//  */
// func (R *Response) IsAnyProtectionEnabled() bool {
// 	return R.Protections.BasicProtections || R.Protections.UserAgent || R.Protections.JWT || R.Protections.RateLimit.Enabled
// }
