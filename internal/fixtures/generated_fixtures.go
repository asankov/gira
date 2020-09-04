package fixtures

//go:generate mockgen -destination gamemodelmock.go  -package fixtures -mock_names GameModel=GameModelMock github.com/asankov/gira/cmd/api/server GameModel
//go:generate mockgen -destination usermodelmock.go  -package fixtures -mock_names UserModel=UserModelMock github.com/asankov/gira/cmd/api/server UserModel
//go:generate mockgen -destination user_games_model_mock.go  -package fixtures -mock_names UserGamesModel=UserGamesModelMock github.com/asankov/gira/cmd/api/server UserGamesModel
//go:generate mockgen -destination authenticatormock.go  -package fixtures -mock_names Authenticator=AuthenticatorMock github.com/asankov/gira/cmd/api/server Authenticator
//go:generate mockgen -destination renderer_mock.go  -package fixtures -mock_names Renderer=RendererMock github.com/asankov/gira/cmd/front-end/server Renderer
//go:generate mockgen -destination api_client_mock.go  -package fixtures -mock_names APIClient=APIClientMock github.com/asankov/gira/cmd/front-end/server APIClient
