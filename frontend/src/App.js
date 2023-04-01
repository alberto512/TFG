import './App.css';

const auxFunc = () => {
  window.location.href =
    'https://apis-sandbox.bancosantander.es/canales-digitales/sb/prestep-authorize?redirect_uri=https://tfg-app.netlify.app/&response_type=code&client_id=bc75ee49-9924-4160-904e-6b246d751e2c';
  return null;
};

function App() {
  return (
    <div className='App'>
      <button onClick={auxFunc}>Hola</button>
    </div>
  );
}

export default App;
