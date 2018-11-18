function placeMarker(address){
    // Create the parameters for the geocoding request:
    var geocodingParams = {
        searchText: address
    };

    // Define a callback function to process the geocoding response:
    var onResult = function(result) {
        var locations = result.Response.View[0].Result,
            position,
            marker;
        // Add a marker for each location found
        for (i = 0;  i < locations.length; i++) {
        position = {
            lat: locations[i].Location.DisplayPosition.Latitude,
            lng: locations[i].Location.DisplayPosition.Longitude
        };
        marker = new H.map.Marker(position);
        map.addObject(marker);
        }
    };

    // Get an instance of the geocoding service:
    var geocoder = platform.getGeocodingService();

    // Call the geocode method with the geocoding parameters,
    // the callback and an error callback function (called if a
    // communication error occurs):
    geocoder.geocode(geocodingParams, onResult, function(e) {
    alert(e);
    });
}

function grabJSON(){
    fetch('http://35.227.91.78/service')
  .then(function(response) {
    return JSON.parse(response.json());
  })
  .then(function(myJson) {
    console.log(JSON.stringify(myJson));
  });

}