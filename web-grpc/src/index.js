import { grpc } from "grpc-web-client";

import { DogTrack, DogTrackClient } from "../generated/dog_pb_service";
import { Dog, DogStatus, Location, TrackRequest } from "../generated/dog_pb";

window.DogStatus = DogStatus;

window.onSubmit = event => {
  event.preventDefault();
  const request = new TrackRequest();
  request.setLocationId(document.querySelector("#location-id").value);
  request.setFloorId("1");
  const grpcRequest = grpc.invoke(DogTrack.TrackDogs, {
    request: request,
    host: "http://localhost:9091",
    metadata: new grpc.Metadata({ HeaderTestKey1: "ClientValue1" }),
    onHeaders: headers => {
      console.log("onHeaders", headers);
    },
    onMessage: message => {
      const dog = message.toObject();
      console.log("message:", message);
      console.log("object:", dog);
      if (dog.status === DogStatus.LYME_DISEASE) {
        console.warn("oh my gosh!!!", dog.name, "has lyme disease!!!!!!!");
      }
    },
    onEnd: (status, statusMessage, trailers) => {
      console.log("onEnd", status, statusMessage, trailers);
    }
  });
};

// setTimeout(() => {
//   grpcRequest.close();
// }, 1000);
