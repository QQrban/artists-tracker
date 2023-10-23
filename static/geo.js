window.onload = async () => {
  const locHTML = document.querySelectorAll('.location-dates-city');
  const geoData = [];
  let updatedCoords = [];

  locHTML.forEach((loc) => {
    geoData.push({
      location: loc.textContent.slice(0, -1),
      coordinates: {
        x: +loc.getAttribute('data-x'),
        y: +loc.getAttribute('data-y'),
      },
    });
  });

  const map = L.map('my-map').setView([0, 0], 2);

  locHTML.forEach((loc) => {
    loc.addEventListener('click', () => {
      const x = +loc.getAttribute('data-x');
      const y = +loc.getAttribute('data-y');
      map.setView([y, x], 7.5);
    });
  });

  const southWest = L.latLng(-71, -160),
    northEast = L.latLng(85, 180);
  const bounds = L.latLngBounds(southWest, northEast);

  map.setMaxBounds(bounds);

  const myAPIKey = '87acceb9227d49a4ae37201d5673afb6';

  const isRetina = L.Browser.retina;

  const baseUrl = `https://maps.geoapify.com/v1/tile/osm-bright/{z}/{x}/{y}.png?apiKey=${myAPIKey}`;
  const retinaUrl = `https://maps.geoapify.com/v1/tile/osm-bright/{z}/{x}/{y}@2x.png?apiKey=${myAPIKey}`;

  L.tileLayer(isRetina ? retinaUrl : baseUrl, {
    attribution:
      'Powered by <a href="https://www.geoapify.com/" target="_blank">Geoapify</a> | © OpenStreetMap <a href="https://www.openstreetmap.org/copyright" target="_blank">contributors</a>',
    apiKey: myAPIKey,
    maxZoom: 20,
    id: 'osm-bright',
  }).addTo(map);

  geoData.forEach((geo) => {
    const popupContent = geo.location ? geo.location : 'Unknown Location';
    const markerPopup = L.popup().setContent(popupContent);

    const zooMarker = L.marker([geo.coordinates.y, geo.coordinates.x])
      .bindPopup(markerPopup)
      .addTo(map);
  });
};
