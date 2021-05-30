package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SearchDir = "./testdata/format_test"
	Excludes  = "./testdata/format_test/web"
	MainFile  = "main.go"
)

func testFormat(t *testing.T, filename, contents, want string) {
	got, err := NewFormatter().Format(filename, []byte(contents))
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func Test_FormatMain(t *testing.T) {
	contents := `package main
	// @title Swagger Example API
	// @version 1.0
	// @description This is a sample server Petstore server.
	// @termsOfService http://swagger.io/terms/

	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @contact.email support@swagger.io

	// @license.name Apache 2.0
	// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

	// @host petstore.swagger.io
	// @BasePath /v2

	// @securityDefinitions.basic BasicAuth

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	// @securitydefinitions.oauth2.application OAuth2Application
	// @tokenUrl https://example.com/oauth/token
	// @scope.write Grants write access
	// @scope.admin Grants read and write access to administrative information

	// @securitydefinitions.oauth2.implicit OAuth2Implicit
	// @authorizationurl https://example.com/oauth/authorize
	// @scope.write Grants write access
	// @scope.admin Grants read and write access to administrative information

	// @securitydefinitions.oauth2.password OAuth2Password
	// @tokenUrl https://example.com/oauth/token
	// @scope.read Grants read access
	// @scope.write Grants write access
	// @scope.admin Grants read and write access to administrative information

	// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
	// @tokenUrl https://example.com/oauth/token
	// @authorizationurl https://example.com/oauth/authorize
	// @scope.admin Grants read and write access to administrative information
	func main() {}`

	want := `package main
	//	@title			Swagger Example API
	//	@version		1.0
	//	@description	This is a sample server Petstore server.
	//	@termsOfService	http://swagger.io/terms/

	//	@contact.name	API Support
	//	@contact.url	http://www.swagger.io/support
	//	@contact.email	support@swagger.io

	//	@license.name	Apache 2.0
	//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

	//	@host		petstore.swagger.io
	//	@BasePath	/v2

	//	@securityDefinitions.basic	BasicAuth

	//	@securityDefinitions.apikey	ApiKeyAuth
	//	@in							header
	//	@name						Authorization

	//	@securitydefinitions.oauth2.application	OAuth2Application
	//	@tokenUrl								https://example.com/oauth/token
	//	@scope.write							Grants write access
	//	@scope.admin							Grants read and write access to administrative information

	//	@securitydefinitions.oauth2.implicit	OAuth2Implicit
	//	@authorizationurl						https://example.com/oauth/authorize
	//	@scope.write							Grants write access
	//	@scope.admin							Grants read and write access to administrative information

	//	@securitydefinitions.oauth2.password	OAuth2Password
	//	@tokenUrl								https://example.com/oauth/token
	//	@scope.read								Grants read access
	//	@scope.write							Grants write access
	//	@scope.admin							Grants read and write access to administrative information

	//	@securitydefinitions.oauth2.accessCode	OAuth2AccessCode
	//	@tokenUrl								https://example.com/oauth/token
	//	@authorizationurl						https://example.com/oauth/authorize
	//	@scope.admin							Grants read and write access to administrative information
	func main() {}`
	testFormat