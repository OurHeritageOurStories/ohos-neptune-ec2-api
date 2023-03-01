# OHOS Neptune api

OHOS is running a Neptune database; the code in this repo is the API running in front of Neptune handling requests.

There are two components here: Kong and a GoLang based API. 

Kong is exposed to the outside world and forwards requests to the API. It is included in this setup to allow for rate limiting, control of methods, etc. 

The Go API has two endpoints:
- `-/`
- `-/sparql`

The blank endpoint acts as a test, to ensure that the AP I is running and available. 

The /sparql endpoint takes two arguments: 
- `sparqlstring=` This requires a partial sparql query, starting from the variables, and ending after the final “}”. 
- `limit=` This is a required variable, and takes a number between 1 and 10,000.
- `prefixString=` Also required, takes a string with the format "PREFIX ex: <http://example.com/exampleOntology#>" (See [this](https://en.wikipedia.org/wiki/SPARQL) page for examples of prefixes)

A valid query, for example, is `curl -d "sparqlstring= ?s ?p ?o where {?s ?p ?o}" -d "limit=10" -d "prefixString=PREFIX ex: <http://example.com/exampleOntology#>" ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/sparql` 

The API then returns the data from Neptune, with the paramters supplied in the command as a sanity-check. 



Also included in this repo are: 
- The dockercompose file used to initiate this API on the EC2 server.
- The config for Kong including the required curl command to send it to the EC2.
- The commands to load an unload data.

Note that the commands to send Kong the config, and to load and unload data need the user to be SSH’d into the EC2. The instructions for this are not here for security. 