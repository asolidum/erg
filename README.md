# ERG - Evaluator for Reverse Geocoders

This tool pings top reverse geocoder providers (currently OpenCage, FourSquare, Yelp, Mapbox, OSM) and returns the results

## Setup
Obtain API keys from each of the providers and add them to a `.env` file. If you don't have one, replace the keys you get from each of the providers in the `.env.sample` file and save as your new `.env` file

### Here are some links for obtaining new reverse geocoder developer API keys
[OpenCage](https://opencagedata.com/)

[Foursquare](https://developer.foursquare.com/)

[Yelp](https://www.yelp.com/developers)

[MapBox](https://www.mapbox.com/search/)


## Program Usage
To get results from a single latitude and longitude simply use the _`--latlon`_ flag with the values comma separated (see below)
```bash
./erg --latlon 47.6124759,-122.3395052
```

`erg` can also use a CSV file as its input. 
```bash
./erg --filename latlon.csv
```
The following flags are useful when running in `file` mode

_`--cols`_ : use to specify the col number for the lat, lon (respectively). Currently defaults to first and second column

_`--append`_ : output in `file` mode is typically of the form
```bash
latitude, longitude, OpenCage, FourSquare, Yelp, MapBox, OSM
```
If you would like these results to be appened to the end of your input file then set this flag to true. Typically this will be used when you obtain a 3rd party POI file and want to compare their data against the reverse geocoders

_`--sleep`_ : Almost all reverse geocoders are rate limited. `erg` can run very quickly but you will likely run into rate limits set by the providers. This value should be set to the strictest API rate limit