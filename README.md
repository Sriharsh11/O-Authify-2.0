# OAuthify

OAuth 2.0 in GoLang

## APIs

- **/addUser** - Adds new users to the postgresql database locally
- **/oauth** - Authenticates users via ROPC(their email and password) and return access tokens to verified users
- **/home** - A route accessible only to authorised users who have proper access tokens

## Details

- **Database Used** - PostgreSQL
- **Models** - Two models 'Users3' and 'login'
- **Tables In DB** - Users3
- **Columns In DB** - name, email and password
- **ORM Used** - GORM
- **HTTP Framework Used** - Gin (gin-gonic/gin)
- **JWT Library Used** - jwt-go

## Functions Used {}

- **init** - Loads environment variables.

- **HashPassword** - Hashes password using bcrypt package before storing in the database.

- **CheckPasswordHash** - Checks a plain-text password against a hash and returns true if the match is successfull.

- **EnterIntoDB** - Calls 'HashPassword' to hash passwords and then store it in the database.

- **GenerateAccessToken** - Generates an access token(JWT).

- **CheckForExistingUser** - Checks the entered user details against existing users in database.

## Handlers Used (/)

- **AddUsers** - Handler for /addUser which posts name, email, password as multipart/form-data. The form detail is stored as a new entry in the database if none of the fields(name, email, password) are empty.

- **AuthenticateUsers** - Handler for /oauth which posts email and password as multipart/form-data. The data is checked against existing entries in database using the email. If an entry exists, the function generates an access token(JWT) using HS256 algorithm. According to this algorithm a shared key is sent along with the access token which is known by both the user and the authenticator. For generating and verifying JWTs, [https://github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) was used.

- **HomeAccess** - Handler for /home which is accessible only if the requesting user sends the proper access token as a request header. This token is then decoded and verified using the shared key and the user is granted access if the match is successfull.

## Testing </>

- **TestEnterIntoDB** - Tests 'EnterIntoDB' function with a dummy name, email and password. It uses [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) (go-sqlmock) for testing database interactions.(creates a fake database server for testing purpose)

- **TestCheckForExistingUser** - Tests 'CheckForExistingUser' function with a dummy email and password. It uses [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) (go-sqlmock) for testing database interactions.(creates a fake database server for testing purpose)

- **TestAddUsers** - Tests 'AddUsers' handler by creating a fake HTTP server(Mock Server).

- **TestAuthenticateUsers** - Tests 'AuthenticateUsers' handler by creating a fake HTTP server(Mock Server).

- **TestHomeAccess** - Tests 'HomeAccess' handler by creating a fake HTTP server(Mock Server).
