import React from "react";
import ReactDOM from "react-dom/client";
import EventList from './components/getEvents';
import 'bootstrap/dist/css/bootstrap.css'
import "./event.css"

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <EventList />
  </React.StrictMode>
);
