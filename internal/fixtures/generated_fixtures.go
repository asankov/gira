package fixtures

//go:generate mockgen -destination gamemodelmock.go  -package fixtures -mock_names GameModel=GameModelMock github.com/gira-games/api/cmd/api/server GameModel
//go:generate mockgen -destination usermodelmock.go  -package fixtures -mock_names UserModel=UserModelMock github.com/gira-games/api/cmd/api/server UserModel
//go:generate mockgen -destination user_games_model_mock.go  -package fixtures -mock_names UserGamesModel=UserGamesModelMock github.com/gira-games/api/cmd/api/server UserGamesModel
//go:generate mockgen -destination franchises_model_mock.go  -package fixtures -mock_names FranchiseModel=FranchiseModelMock github.com/gira-games/api/cmd/api/server FranchiseModel
//go:generate mockgen -destination authenticatormock.go  -package fixtures -mock_names Authenticator=AuthenticatorMock github.com/gira-games/api/cmd/api/server Authenticator
