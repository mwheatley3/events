import React, { useEffect, useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import './App.css';
import ListEvents from './list_events.js';

const baseURL = "http://localhost:8080"

function App() {

  // manage local state for setting params for en event
  const [type, setType] = useState('INCREMENT');
  const [value, setValue] = useState(0);
  // const [refresh, setRefresh] = useState(0);

  // state for system value and refreshing that value
  const [systemValue, setSystemValue] = useState(0);
  const [t, setT] = useState(-1);

  // state for listing events
  const [r, getEvents] = useState(0);
  const [events, setEvents] = useState([]);

  const getValue = () => {
    console.log("getValue t:", t);
    const url = (t >= 0) ? `${baseURL}/value/${t}` : `${baseURL}/value`;
    fetch(url)
      .then(r => r.json())
      .then(r => {console.log('r', r); setSystemValue(r);})
      .catch(e => console.error(e))
  }

  useEffect(getValue, [t])

  useEffect(() => {
    fetch(`${baseURL}/events`)
      .then(r => r.json())
      .then(r => setEvents(r))
      .catch(e => console.error(e))
  }, [r])

  const submitEvent = (e) => {
    e.preventDefault()
    if (value <= 0) {
      return
    }

    const url = `${baseURL}/event`
    const body = {type, value}

    fetch(url, {
      method: 'POST',
      body: JSON.stringify(body)
    })
    .then(res => res.json())
    .then(res => console.log(res))
    .then(x => getValue())
    .catch(e => console.error(e));

  }

  //TODO add some paging mechanism for events
  // const truncatedEvents = events.slice(Math.max(events.length - 20, 0))

  return (
    <div className="card" style={{ width: '30rem' }}>
      <div className="card-body">
        <h4 className="card-title">curbFlow Event Sourcing</h4>
        <p className="card-text">Bringing order to the chaos of city streets, starting at the curb.</p>
      </div>
      <div className="card-body">
        <Form onSubmit={ submitEvent }>
          <Form.Label>Event Type</Form.Label>
          <Form.Control as="select" placeholder="hello" onChange={ e => setType(e.target.value)}>
            <option value="INCREMENT">Increment</option>
            <option value="DECREMENT">Decrement</option>
          </Form.Control>
          <Form.Text className="text-muted">
            choose Increment to increase the value.
            choose Decrement to decrease the value
          </Form.Text>
          <Form.Label>Value</Form.Label>
          <Form.Control placeholder="0" type="number" onChange={ e => setValue(parseInt(e.target.value))}/>
          <Button variant="primary" type="submit">
            Submit
          </Button>
        </Form>
      </div>
      <div className="card-body">
        <h5 className="card-title">Refresh The Current Value of the System</h5>
        <h6 className="card-subtitle mb-2 text-muted">Current Value: {systemValue}</h6>
        <p className="card-text">
        A producer is constantly emitting events,
        type a value to get the system value at that point
        </p>
        <Form.Control placeholder="0" type="number" onChange={ e => setT(parseInt(e.target.value))}/>
      </div>
      <div className="card-body">
        <Button variant="primary" onClick={() => getEvents(r+1) }>Refresh List of Events</Button>
        <ListEvents events={events} md="auto"/>
      </div>
    </div>
  );
}

export default App;
