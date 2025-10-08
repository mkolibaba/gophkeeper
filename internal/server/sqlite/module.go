package sqlite

import (
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"sqlite",
	fx.Provide(
		NewConfig,
		NewDB,
		NewQueries,
		fx.Annotate(NewUserService, fx.As(new(server.UserService))),
		fx.Annotate(NewLoginService, fx.As(new(server.LoginService))),
		fx.Annotate(NewNoteService, fx.As(new(server.NoteService))),
		fx.Annotate(NewBinaryService, fx.As(new(server.BinaryService))),
		fx.Annotate(NewCardService, fx.As(new(server.CardService))),
	),
	fx.Invoke(
		OpenDB,
	),
)

func NewQueries(db *DB) *sqlc.Queries {
	return sqlc.New(db.db)
}

func OpenDB(db *DB) error {
	return db.Open()
}
