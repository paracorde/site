package utils

import (
   "net/http"

   "github.com/labstack/echo/v4"
)

func SetCookie(context echo.Context, name string, value string){
   cookie := &http.Cookie{
      Name: name,
      Value: value,
      Path: "/",
   }
   context.SetCookie(cookie)
}

func Flash(context echo.Context, message string) {
   SetCookie(context, "flash", message)
}

func GetCookieOrEmpty(context echo.Context, name string) string {
   cookie, err := context.Cookie(name)
   if err != nil {
      return ""
   }
   return cookie.Value
}

func GetFlash(context echo.Context) string {
   return GetCookieOrEmpty(context, "flash")
}
