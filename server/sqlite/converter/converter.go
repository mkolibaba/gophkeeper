package converter

import (
	"context"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

// goverter:converter
// goverter:output:file ./gen/converter.go
type DataConverter interface {
	// -- Login --

	// goverter:context ctx
	// goverter:map User | UserFromContext
	ConvertToInsertLogin(ctx context.Context, source server.LoginData) sqlc.InsertLoginParams

	ConvertToLoginDataSlice(source []sqlc.Login) []server.LoginData

	// goverter:useZeroValueOnPointerInconsistency
	ConvertToLoginData(source sqlc.Login) server.LoginData

	ConvertToUpdateLogin(source sqlc.Login) sqlc.UpdateLoginParams

	// goverter:update target
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignore ID
	ConvertToUpdateLoginUpdate(source server.LoginDataUpdate, target *sqlc.UpdateLoginParams)

	// -- Note --

	// goverter:context ctx
	// goverter:map User | UserFromContext
	ConvertToInsertNote(ctx context.Context, source server.NoteData) sqlc.InsertNoteParams

	ConvertToNoteDataSlice(source []sqlc.Note) []server.NoteData

	// goverter:useZeroValueOnPointerInconsistency
	ConvertToNoteData(source sqlc.Note) server.NoteData

	ConvertToUpdateNote(source sqlc.Note) sqlc.UpdateNoteParams

	// goverter:update target
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignore ID
	ConvertToUpdateNoteUpdate(source server.NoteDataUpdate, target *sqlc.UpdateNoteParams)

	// -- Binary --

	// goverter:context ctx
	// goverter:map User | UserFromContext
	// goverter:autoMap BinaryData
	ConvertToInsertBinary(ctx context.Context, source server.ReadableBinaryData) sqlc.InsertBinaryParams

	ConvertToBinaryDataSlice(source []sqlc.Binary) []server.BinaryData

	// goverter:useZeroValueOnPointerInconsistency
	ConvertToBinaryData(source sqlc.Binary) server.BinaryData

	ConvertToUpdateBinary(source sqlc.Binary) sqlc.UpdateBinaryParams

	// goverter:update target
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignore ID
	ConvertToUpdateBinaryUpdate(source server.BinaryDataUpdate, target *sqlc.UpdateBinaryParams)

	// -- Card --

	// goverter:context ctx
	// goverter:map User | UserFromContext
	// goverter:map CVV Cvv
	ConvertToInsertCard(ctx context.Context, source server.CardData) sqlc.InsertCardParams

	ConvertToCardDataSlice(source []sqlc.Card) []server.CardData

	// goverter:useZeroValueOnPointerInconsistency
	// goverter:map Cvv CVV
	ConvertToCardData(source sqlc.Card) server.CardData

	ConvertToUpdateCard(source sqlc.Card) sqlc.UpdateCardParams

	// goverter:update target
	// goverter:useZeroValueOnPointerInconsistency
	// goverter:ignore ID
	// goverter:map CVV Cvv
	ConvertToUpdateCardUpdate(source server.CardDataUpdate, target *sqlc.UpdateCardParams)
}

// goverter:context ctx
func UserFromContext(ctx context.Context) string {
	return server.UserFromContext(ctx)
}
