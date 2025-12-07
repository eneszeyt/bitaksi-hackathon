import { useEffect, useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { getNearbyDrivers } from '../services/api';
import 'leaflet/dist/leaflet.css';

// fix for map icons (solving a known bug between leaflet and react)
import L from 'leaflet';
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';

let DefaultIcon = L.icon({
    iconUrl: icon,
    shadowUrl: iconShadow,
    iconSize: [25, 41],
    iconAnchor: [12, 41]
});
L.Marker.prototype.options.icon = DefaultIcon;

const CENTER = [41.0, 29.0];

export default function MapPage() {
  const [drivers, setDrivers] = useState([]);
  // 1. new feature: filter state
  const [filterType, setFilterType] = useState('all'); // 'all', 'yellow', 'black'

  // extracted data fetching function to reuse it
  const loadData = async () => {
    try {
      // if 'all' is selected, send empty string for type, otherwise send the selected type
      const typeParam = filterType === 'all' ? '' : filterType;
      
      // request to backend
      const res = await getNearbyDrivers(41.0, 29.0, typeParam);
      setDrivers(res.data || []);
    } catch (error) {
      console.error("Failed to fetch data", error);
      if (error.response && error.response.status === 401) {
          window.location.href = '/';
      }
    }
  };

  // 2. new feature: re-fetch data when filterType changes
  useEffect(() => {
    loadData();
  }, [filterType]);

  // 3. new feature: call taxi function
  const handleCallTaxi = (driverName) => {
    alert(`🎉 ${driverName} is on the way! Arriving in 3 minutes.`);
  };

  return (
    <div className="h-screen flex flex-col">
      <header className="bg-gray-800 p-4 flex justify-between items-center text-white shadow-md z-10">
        <div className="flex items-center gap-2">
            <span className="text-2xl">🚕</span>
            <h1 className="text-xl font-bold text-taxi">TaxiHub Monitor</h1>
        </div>

        {/* 4. new feature: filter buttons */}
        <div className="flex gap-2 bg-gray-700 p-1 rounded-lg">
            <button 
                onClick={() => setFilterType('all')}
                className={`px-3 py-1 rounded text-sm transition ${filterType === 'all' ? 'bg-gray-500 text-white' : 'text-gray-300 hover:bg-gray-600'}`}
            >
                All
            </button>
            <button 
                onClick={() => setFilterType('yellow')}
                className={`px-3 py-1 rounded text-sm transition ${filterType === 'yellow' ? 'bg-taxi text-black font-bold' : 'text-gray-300 hover:bg-gray-600'}`}
            >
                Yellow
            </button>
            <button 
                onClick={() => setFilterType('black')}
                className={`px-3 py-1 rounded text-sm transition ${filterType === 'black' ? 'bg-black text-white border border-gray-500' : 'text-gray-300 hover:bg-gray-600'}`}
            >
                Black
            </button>
        </div>

        <div className="flex items-center gap-4">
            <span className="bg-gray-700 px-3 py-1 rounded text-sm text-white">
                {drivers.length} Vehicles
            </span>
            <button 
                onClick={() => { localStorage.clear(); window.location.href = '/' }}
                className="text-sm text-red-400 hover:text-red-300 underline"
            >
                Logout
            </button>
        </div>
      </header>

      <div className="flex-1 relative z-0">
        <MapContainer center={CENTER} zoom={13} scrollWheelZoom={true}>
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          
          {drivers.map((driver) => (
            <Marker 
              key={driver.id} 
              position={[driver.location.lat, driver.location.lon]}
            >
              <Popup>
                <div className="text-center min-w-[150px]">
                  <h3 className="font-bold text-lg">{driver.firstName} {driver.lastName}</h3>
                  <p className="text-gray-600 font-mono text-sm mb-1">{driver.plate}</p>
                  
                  <div className="flex justify-center gap-2 mb-3">
                    <span className={`px-2 py-0.5 rounded text-xs font-bold inline-block border ${driver.taxiType === 'yellow' ? 'bg-taxi text-black border-yellow-500' : 'bg-black text-white border-gray-600'}`}>
                        {driver.taxiType === 'yellow' ? 'Yellow Taxi' : 'Black Taxi'}
                    </span>
                  </div>

                  {/* 5. new feature: call button */}
                  <button 
                    onClick={() => handleCallTaxi(driver.firstName)}
                    className="w-full bg-blue-600 text-white text-sm py-1.5 rounded hover:bg-blue-700 transition"
                  >
                    CALL TAXI
                  </button>
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      </div>
    </div>
  );
}