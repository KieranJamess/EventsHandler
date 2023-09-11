import { useState, useEffect } from "react";

interface ListEventsProps {
  events: Event[] | null;
  onSave: (editedEvent: Event) => void;
}

interface Event {
  id: number;
  name: string;
  startDateTime: string;
  endDateTime: string;
}

function GetEvents() {
  const [events, setEvents] = useState<Event[] | null>(null);
  useEffect(() => {
    fetch("http://127.0.0.1:3000/events")
      .then((response) => response.json())
      .then((data) => {
        setEvents(data);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
      });
  }, []);

  const handleSave = async (editedEvent: Event) => {
    try {
      const response = await fetch(`http://127.0.0.1:3000/events/${editedEvent.id}`, {
        method: "PATCH", 
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: editedEvent.name }),
      });

      if (response.ok) {
        const updatedEvents = events?.map((event) =>
          event.id === editedEvent.id ? editedEvent : event
        );
        setEvents(updatedEvents || []);
      } else {
        console.error("Failed to update event:", response.statusText);
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  return <ListEvents events={events} onSave={handleSave} />;
}

function ListEvents({ events, onSave }: ListEventsProps) {
  const [editableCell, setEditableCell] = useState<Event | null>(null);

  const handleEdit = (event: Event) => {
    setEditableCell(event);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (editableCell) {
      setEditableCell({
        ...editableCell,
        name: e.target.value,
      });
    }
  };

  const handleBlur = () => {
    if (editableCell) {
      onSave(editableCell); 
      setEditableCell(null); 
    }
  };

  return (
    <>
      <div className="title-container">
        <h1 className="title-box">
          <span className="title">Event List</span>
        </h1>
      </div>
      <div className="events-containers">
        {events !== null ? (
          <table className="table table-hover">
            <thead>
              <tr>
                <th>Name</th>
                <th>Start Date</th>
                <th>End Date</th>
              </tr>
            </thead>
            <tbody>
              {events.map((event) => (
                <tr key={event.id}>
                  <td>
                    {editableCell?.id === event.id ? (
                      <input
                        type="text"
                        value={editableCell.name}
                        onChange={handleInputChange}
                        onBlur={handleBlur}
                        autoFocus
                      />
                    ) : (
                      <span
                        onClick={() => handleEdit(event)}
                        style={{ cursor: "pointer" }}
                      >
                        {event.name}
                      </span>
                    )}
                  </td>
                  <td>{event.startDateTime}</td>
                  <td>{event.endDateTime}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p>No events found</p>
        )}
      </div>
    </>
  );
}



export default GetEvents;
