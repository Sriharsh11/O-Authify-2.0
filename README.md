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

<!-- ## Code

... -->

## Functions Used

- **HashPassword** - Hashes password using bcrypt package before storing in the database.

- **CheckPasswordHash** - Checks a plain-text password against a hash and return true if the match is successfull.

- **AddUsers** - Handler for /addUser which posts name, email, password as multipart/form-data. The form detail is stored as a new entry in the database if none of the fields(name, email, password) are empty.

- **AuthenticateUsers** - Handler for /oauth which posts email and password as multipart/form-data. The data is checked against existing entries in database using the email. If an entry exists, the function generates an access token(JWT) using HS256 algorithm. According to this algorithm a shared key is sent along with the access token which known by both the user and the authenticator.

- **HomeAccess** - Handler for /home which is accessible only if the requesting user sends the proper access token as a request header. This token is then decoded and verified using the shared key and the user is granted access if the match is successfull.

<!-- ### Prerequisites

What things you need to install the software and how to install them

```
Give examples
```

### Installing

A step by step series of examples that tell you how to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo -->

<!-- ## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

- [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
- [Maven](https://maven.apache.org/) - Dependency Management
- [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).

## Authors

- **Billie Thompson** - _Initial work_ - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

- Hat tip to anyone whose code was used
- Inspiration
- etc -->
