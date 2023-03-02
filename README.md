# OHOS Neptune api

OHOS is running a Neptune database; the code in this repo is the API running in front of Neptune handling requests.

There are two components here: Kong and a GoLang based API. 

Kong is exposed to the outside world and forwards requests to the API. It is included in this setup to allow for rate limiting, control of methods, etc. 

The Go API has two endpoints:
- `-/`
- `-/sparql`

The blank endpoint acts as a test, to ensure that the AP I is running and available. 

The /sparql endpoint takes two arguments: 
- `sparqlquery=` This requires a sparql query, and will only accept read-only ones (e.g. select, sort, etc)
- `limit=` This is a required variable, and takes a number between 1 and 10,000.

A valid query, for example, is `curl -d "sparqlquery= select ?s ?p ?o where {?s ?p ?o}" -d "limit=10" ec2-13-40-156-226.eu-west-2.compute.amazonaws.com:5000/sparql` 

The API then returns the data from Neptune, with the paramters supplied in the command as a sanity-check. 



Also included in this repo are: 
- The dockercompose file used to initiate this API on the EC2 server.
- The config for Kong including the required curl command to send it to the EC2.
- The commands to load an unload data.

Note that the commands to send Kong the config, and to load and unload data need the user to be SSH’d into the EC2. The instructions for this are not here for security. 