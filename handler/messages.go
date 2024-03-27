package handler

const (
	ServerError   = "Internal Server Error."
	TokenNotFound = "Token not found."
	TokenInvalid  = "Invalid token."
	OnlyAdmin     = "Forbidden. Only admins are allowed to access this resource."
)

const (
	MustPOST   = "The request method must be POST."
	MustDELETE = "The request method must be DELETE."
	MustPUT    = "The request method must be PUT."
	MustGET    = "The request method must be GET."
)

const (
	InvalidBirthday           = "Invalid birthday format. Please use YYYY-MM-DD."
	InvalidReleaseDate        = "Invalid release_date format. Please use YYYY-MM-DD."
	InvalidReplaceReleaseDate = "Invalid replace_release_date format. Please use YYYY-MM-DD."
)

const (
	MissingRequiredFieldsActor = "Required fields (name_actor, gender, birthday) are missing."
	MissingRequiredFieldsFilm  = "Required fields (name, release_date) are missing."
)

const (
	MissingRequiredArgNameActor = "The required 'name_actor' argument is missing."
	MissingRequiredArgNameFilm  = "The required 'name_film' argument is missing."
)

const (
	FilmCreated = "Film successfully create."
	FilmDeleted = "Film successfully deleted."
	FilmChange  = "Film successfully change."
)

const (
	ErrFilmCreated    = "This film has already been added to the database."
	ErrFilmNotCreated = "This film has not yet been added to the database."
)

const (
	ActorCreated = "Actor successfully create."
	ActorDeleted = "Actor successfully deleted."
	ActorChange  = "Actor successfully change."
)

const (
	ErrActorCreated    = "This actor has already been added to the database."
	ErrActorNotCreated = "This actor has not yet been added to the database."
)
