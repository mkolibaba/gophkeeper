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
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/fx"
)

//go:embed all:frontend/dist
var assets embed.FS

type LoginForm struct {
	Login    string `form:"login"`
	Password string `form:"password"`
}

type RegisterForm struct {
	Login    string `form:"login"`
	Password string `form:"password"`
	Repeat   string `form:"repeat"`
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

	e.GET("/form/authorization", echo.WrapHandler(templ.Handler(components.AuthorizationForm("", ""))))
	e.GET("/form/registration", echo.WrapHandler(templ.Handler(components.RegistrationForm("", ""))))
	e.GET("/vault", func(c *echo.Context) error {
		login := userService.GetUserLogin()

		logins, err := loginService.GetAll(c.Request().Context())
		if err != nil {
			return render(c, components.Vault(login, components.VaultData{}, "Error while loading data: "+err.Error()))
		}

		binaries, err := binaryService.GetAll(c.Request().Context())
		if err != nil {
			return render(c, components.Vault(login, components.VaultData{}, "Error while loading data: "+err.Error()))
		}

		notes, err := noteService.GetAll(c.Request().Context())
		if err != nil {
			return render(c, components.Vault(login, components.VaultData{}, "Error while loading data: "+err.Error()))
		}

		cards, err := cardService.GetAll(c.Request().Context())
		if err != nil {
			return render(c, components.Vault(login, components.VaultData{}, "Error while loading data: "+err.Error()))
		}

		return render(c, components.Vault(login, components.VaultData{
			Logins: logins,
			Notes:  notes,
			Binary: binaries,
			Cards:  cards,
		}, ""))
	})
	e.POST("/login", func(c *echo.Context) error {
		var form LoginForm
		if err := c.Bind(&form); err != nil {
			return render(c, components.AuthorizationForm("", "Invalid request"))
		}

		token, err := authService.Authorize(c.Request().Context(), form.Login, form.Password)
		if err != nil {
			return render(c, components.AuthorizationForm(form.Login, "Invalid login and/or password"))
		}

		userService.SetInfo(form.Login, token)

		return c.Redirect(http.StatusFound, "/vault")
	})
	e.POST("/register", func(c *echo.Context) error {
		var form RegisterForm
		if err := c.Bind(&form); err != nil {
			return render(c, components.RegistrationForm("", "Invalid request"))
		}

		if form.Password == "" {
			return render(c, components.RegistrationForm(form.Login, "Password should not be empty"))
		}
		if form.Password != form.Repeat {
			return render(c, components.RegistrationForm(form.Login, "Passwords should match"))
		}

		token, err := authService.Register(c.Request().Context(), form.Login, form.Password)
		if err != nil {
			return render(c, components.RegistrationForm(form.Login, "Can not register: "+err.Error()))
		}

		userService.SetInfo(form.Login, token)

		return c.Redirect(http.StatusFound, "/vault")
	})

	return e
}

func render(c *echo.Context, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
}
