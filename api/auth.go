package api

import (
  "net/http"
  "github.com/dgrijalva/jwt-go"
  "github.com/dgrijalva/jwt-go/request"
  "os"
  "time"
)

func Authorize(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, keyLookupFunc)

  	if err == nil && token.Valid {
      next.ServeHTTP(w, r)
      return
  	}

    w.WriteHeader(http.StatusUnauthorized)
  })
}

func getSigningKey() []byte {
  return []byte(os.Getenv("APP_SECRET_KEY"))
}

 /* Handlers */
 var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
   // TODO: Check with database here.
   token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "admin": true,
    "name": "Danny",
    "exp": time.Now().Add(time.Hour * 24).Unix(),
  })

   /* Sign the token with our secret */
   tokenString, _ := token.SignedString(getSigningKey())

   /* Finally, write the token to the browser window */
   w.Write([]byte(tokenString))
 })


 func keyLookupFunc(t *jwt.Token) (interface{}, error) {
     // Look up key in database.
     // TODO

     //
     return getSigningKey(), nil
 }
