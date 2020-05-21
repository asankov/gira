package fixtures

//go:generate mockgen -destination gamemodelmock.go  -package fixtures -mock_names GameModel=GameModelMock github.com/asankov/gira/cmd/api/server GameModel
//go:generate mockgen -destination usermodelmock.go  -package fixtures -mock_names UserModel=UserModelMock github.com/asankov/gira/cmd/api/server UserModel
//go:generate mockgen -destination authenticatormock.go  -package fixtures -mock_names Authenticator=AuthenticatorMock github.com/asankov/gira/cmd/api/server Authenticator
