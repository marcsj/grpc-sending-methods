window.onSubmit = event => {
  event.preventDefault();
  const locationID = document.querySelector("#location-id").value;
  const floorID = "1";
  const object = {
      location_id: locationID,
      floor_id: floorID,
  }
  // making a websocket request, but adding in a param for method overriding for our proxy server
  const ws = new WebSocket(
    `ws://${window.location.hostname}:8081/v1/dogs/track?method=POST`
  );
  ws.addEventListener("open", event => {
    ws.send(JSON.stringify(object));
    console.log("open", event);
  });
  ws.addEventListener("message", event => {
    console.log("message", event);
  });
  ws.addEventListener("error", event => {
    console.log("error", event);
  });
  ws.addEventListener("close", event => {
    console.log("close", event);
  });
};
