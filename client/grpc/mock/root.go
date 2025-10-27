//go:generate moq -stub -pkg mock -out card.go . CardServiceClientForMocking
//go:generate moq -stub -pkg mock -out login.go . LoginServiceClientForMocking
//go:generate moq -stub -pkg mock -out note.go . NoteServiceClientForMocking

package mock

import "github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"

// TODO: разобраться в настройках moq и убрать эти декларации
type CardServiceClientForMocking = gophkeeperv1.CardServiceClient
type LoginServiceClientForMocking = gophkeeperv1.LoginServiceClient
type NoteServiceClientForMocking = gophkeeperv1.NoteServiceClient
