package middleware

import (
	"fmt"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/Deiklov/diplom_backend/pkg/metrics"
	"github.com/Deiklov/diplom_backend/pkg/sanitize"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type Middleware struct {
	origins    map[string]struct{}
	serverMode string

	metr     metrics.Metrics
	uUseCase user.UseCase
}

func CreateMiddleware(uUseCase_ user.UseCase) *Middleware {
	origins_ := make(map[string]struct{})
	// TODO: вайпер
	for _, key := range viper.GetStringSlice("cors.allowed_origins") {
		origins_[key] = struct{}{}
	}
	return &Middleware{
		origins:    origins_,
		serverMode: viper.GetString("server.mode"),

		uUseCase: uUseCase_,
	}
}

func (mw *Middleware) SetMetrics(metr metrics.Metrics) {
	mw.metr = metr
}

//func (mw *Middleware) SetWebSocketPool(wsPool_ webSocketPool.WebSocketPool) {
//	mw.wsPool = wsPool_
//}

func (mw *Middleware) RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		res := next(ctx)
		logger.Infof("%s %s %d %s",
			ctx.Request().Method,
			ctx.Request().RequestURI,
			ctx.Response().Status,
			time.Since(start))
		return res
	}
}

func (mw *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		origin := ctx.Request().Header.Get("Origin")
		if _, exist := mw.origins[origin]; !exist {
			return ctx.NoContent(http.StatusForbidden)
		}
		ctx.Response().Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		if ctx.Request().Method == "OPTIONS" {
			return ctx.NoContent(http.StatusOK)
		}
		return next(ctx)
	}
}

func (mw *Middleware) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("ProcessPanic up on ", ctx.Request().Method, ctx.Request().URL.Path)
				fmt.Println("Panic statement: ", err)
				err := ctx.NoContent(http.StatusInternalServerError)
				if err != nil {
					logger.Error(err)
				}
			}
		}()
		return next(ctx)
	}
}

func (mw *Middleware) Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)
		var status int
		if err != nil {
			status = err.(*echo.HTTPError).Code
		} else {
			status = ctx.Response().Status
		}
		mw.metr.ObserveResponseTime(status, ctx.Request().Method, ctx.Path(), time.Since(start).Seconds())
		mw.metr.IncHits(status, ctx.Request().Method, ctx.Path())
		return err
	}
}

func (mw *Middleware) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()
		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("body", sanBody)
		return next(ctx)
	}
}

func (mw *Middleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.String(http.StatusUnauthorized, errors.NoCookie)
		}
		sid := cookie.Value
		uid, has := "bfa52e31-6991-48eb-8d96-1ac18fd50bd5", true
		if !has {
			// Пришла невалидная кука, стираем её из браузера
			newCookie := http.Cookie{Name: "session_id", Value: sid, Expires: time.Now().AddDate(-1, 0, 0)}
			ctx.SetCookie(&newCookie)
			return ctx.String(http.StatusUnauthorized, errors.NoCookie)
		}
		ctx.Set("uid", uid)
		ctx.Set("sid", sid)
		return next(ctx)
	}
}

func (mw *Middleware) DebugMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if mw.serverMode == "debug" {
			dump, err := httputil.DumpRequest(ctx.Request(), true)
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			logger.Debugf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump)
		}
		return next(ctx)
	}
}
