create a simple gin api that users can signup and login after signup.

the user can create a url target list for scan by sending a url list and the service name to the api with a post request server then will save this url list to the database with unique id trackable by the user.

the user can get the url list from the database with a get request to the api.

the user can delete the url list from the database with a delete request to the api.

the user can update the url list from the database with a put request to the api.

the user then can select scanner between burp suite and nuclei and zaproxy.

the user can start scan the url list by the selected scanner.

api should be able to poll the scan from the scanner and save the result to the database.

the user can get the scan result from the database with a get request to the api.

the user can delete the scan result from the database with a delete request to the api.

use a api file structure to organize the api.

use a database file structure to organize the database.

use a model file structure to organize the model.

and use a controller file structure to organize the controller.

use a router file structure to organize the router.

use a middleware file structure to organize the middleware.

and use a config file structure to organize the config.

and use a log file structure to organize the log.

and put each function in the file and folder that is relevant to it.

and use a .env file to organize the environment variable.

and use a .gitignore file to organize the file that should not be pushed to the repository.


