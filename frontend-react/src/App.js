import React, { useEffect, useState } from "react";
import Profile from "./Profile";
import "./App.css";

function App() {
  const [data, setData] = useState(null);

  useEffect(() => {
    fetch("http://localhost:8080")
      .then((res) => res.json())
      .then(setData)
      .catch(console.error);
  }, []);

  if (!data) return <div>Loading...</div>;
  return <Profile {...data} />;
}

export default App;
