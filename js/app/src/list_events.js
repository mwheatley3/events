import React from 'react';
import Table from 'react-bootstrap/Table';

export default function ListEvents(props) {
  return (
    <Table striped bordered hover size="sm">
      <thead>
        <tr>
          <th lg="1">#</th>
          <th>Type</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {props.events.map((ev, i) => (
          <tr key={i}>
            <td>{i + 1}</td>
            <td>{ev.type}</td>
            <td>{ev.value}</td>
          </tr>
        ))
      }
      </tbody>
    </Table>
    )
}
