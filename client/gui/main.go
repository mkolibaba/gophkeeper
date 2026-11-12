package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mkolibaba/gophkeeper/client/core"
	coremodule "github.com/mkolibaba/gophkeeper/client/core/module"
	"github.com/mkolibaba/gophkeeper/client/gui/components"
	"github.com/starfederation/datastar-go/datastar"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/fx"
)

//go:embed all:frontend/dist
var assets embed.FS

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Repeat   string `json:"repeat"`
}

func main() {
	var wailsApp *application.App

	fxApp := fx.New(
		coremodule.Module,
		fx.Provide(newRouter),
		fx.Provide(newApplication),
		fx.Populate(&wailsApp),
	)

	ctx := context.Background()
	if err := fxApp.Start(ctx); err != nil {
		log.Fatal(err)
	}

	runErr := wailsApp.Run()

	if err := fxApp.Stop(ctx); err != nil {
		log.Println(err)
	}
	if runErr != nil {
		log.Fatal(runErr)
	}
}

func EchoMiddleware(e *echo.Echo) application.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/wails") {
				next.ServeHTTP(w, r)
				return
			}

			// Wails (WebView2) не всегда проставляет Content-Length,
			// из-за чего echo.Bind пропускает тело запроса.
			if r.ContentLength == 0 && r.Body != nil && r.Body != http.NoBody {
				r.ContentLength = -1
			}

			e.ServeHTTP(w, r)
		})
	}
}

func newApplication(router *echo.Echo) *application.App {
	app := application.New(application.Options{
		Name:        "Gophkeeper",
		Description: "Gophkeeper client",
		Assets: application.AssetOptions{
			Handler:    router,
			Middleware: EchoMiddleware(router),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Gophkeeper",
		Width:  1000,
		Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	return app
}

func newRouter(
	authService core.AuthorizationService,
	userService core.UserService,
	loginService core.LoginService,
	binaryService core.BinaryService,
	noteService core.NoteService,
	cardService core.CardService,
) *echo.Echo {
	e := echo.New()
	e.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	e.Use(middleware.RequestLogger())

	e.Filesystem = echo.MustSubFS(assets, "frontend/dist")
	e.File("/", "index.html")
	e.Static("/static", ".")

	e.GET("/form/authorization", datastarTemplWrapper(components.AuthorizationForm()))
	e.GET("/form/registration", datastarTemplWrapper(components.RegistrationForm()))

	e.POST("/login", func(c *echo.Context) error {
		sse := datastar.NewSSE(c.Response(), c.Request())

		var form LoginForm
		if err := datastar.ReadSignals(c.Request(), &form); err != nil {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Invalid request",
			})
		}

		// Сбрасываем пароль на форме
		sse.MarshalAndPatchSignals(map[string]any{
			"password": "",
		})

		token, err := authService.Authorize(c.Request().Context(), form.Login, form.Password)
		if err != nil {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Invalid login and/or password",
			})
		}

		userService.SetInfo(form.Login, token)

		// TODO: поправить после исправления апи
		logins, err := loginService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		binaries, err := binaryService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		notes, err := noteService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		cards, err := cardService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		return sse.PatchElementTempl(components.Vault(components.VaultData{
			Logins: logins,
			Notes:  notes,
			Binary: binaries,
			Cards:  cards,
		}))
	})
	e.POST("/register", func(c *echo.Context) error {
		sse := datastar.NewSSE(c.Response(), c.Request())

		var form RegisterForm
		if err := datastar.ReadSignals(c.Request(), &form); err != nil {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Invalid request",
			})
		}

		sse.MarshalAndPatchSignals(map[string]any{
			"password": "",
			"repeat":   "",
		})

		if form.Password == "" {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Password should not be empty",
			})
		}
		if form.Password != form.Repeat {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Passwords should match",
			})
		}

		token, err := authService.Register(c.Request().Context(), form.Login, form.Password)
		if err != nil {
			return sse.MarshalAndPatchSignals(map[string]any{
				"errMsg": "Can not register: "+err.Error(),
			})
		}

		userService.SetInfo(form.Login, token)

		// TODO: поправить после исправления апи
		logins, err := loginService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		binaries, err := binaryService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		notes, err := noteService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		cards, err := cardService.GetAll(c.Request().Context())
		if err != nil {
			// return render(c, components.Vault(components.VaultData{}))
		}

		return sse.PatchElementTempl(components.Vault(components.VaultData{
			Logins: logins,
			Notes:  notes,
			Binary: binaries,
			Cards:  cards,
		}))
	})

	return e
}

func datastarTemplWrapper(cmp templ.Component) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sse := datastar.NewSSE(c.Response(), c.Request())
		return sse.PatchElementTempl(cmp)
	}
}
