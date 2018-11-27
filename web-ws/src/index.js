window.onSubmit = event => {
  event.preventDefault();
  const locationID = document.querySelector("#location-id").value;
  const floorID = "1";
  const ws = new WebSocket(
    // `ws://${window.location.hostname}:8080/v1/dogs/track`
    `ws://${window.location.hostname}:8080/`
  );
  ws.addEventListener("open", event => {
    console.log("open", event);
    ws.send(
      JSON.stringify({
        command: "dogs.track",
        locationID,
        floorID
      })
    );
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
